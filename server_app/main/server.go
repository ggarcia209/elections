// The server service is primarily responsible for serving static files to the web client,
// and retrieving the information requested by the user from the Index service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/elections/source/dynamo"

	"github.com/golang/protobuf/ptypes"

	"github.com/elections/source/server"
	ind "github.com/elections/source/svc/index"
	pb "github.com/elections/source/svc/proto"
	"github.com/elections/source/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ViewServer implements the ViewServer gRPC interface
type viewServer struct {
	pb.UnimplementedViewServer
	mu sync.Mutex
}

var (
	crt        = "../cert/server.crt"
	key        = "../cert/server.key"
	clientCert = "../cert/server.crt"
)

var hostname string
var client ind.IndexClient

var rankingsCache server.RankingsMap
var yrTotalsCache server.YrTotalsMap
var searchDataCache server.SearchDataMap

var database *dynamo.DbInfo
var metadata *server.IndexData

func main() {
	var err error
	serverAddr := "127.0.0.1:9092" // index server address
	var opts []grpc.DialOption
	var sOpts []grpc.ServerOption

	fmt.Println("initializing disk cache...")
	server.InitServerDiskCache()

	fmt.Println("getting hostname...")
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Println("failed to get hostname")
		os.Exit(1)
	}
	fmt.Println("host: ", hostname)

	// Create the client TLS credentials
	fmt.Println("initializing index client...")
	fmt.Println("loading index client credentials...")
	cCreds, err := credentials.NewClientTLSFromFile(clientCert, "")
	if err != nil {
		fmt.Printf("could not load tls cert: %s\n", err)
		os.Exit(1)
	}
	ccr := grpc.WithTransportCredentials(cCreds)
	opts = append(opts, ccr)

	fmt.Printf("dialing %s...\n", serverAddr)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure()) // change after testing
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client = ind.NewIndexClient(conn)
	fmt.Println("connected to ", serverAddr)

	fmt.Println("getting rankings and yearly totals caches...")
	resp, err := getCaches(client, hostname)
	if err != nil {
		fmt.Println("failed to get caches: ", err)
		os.Exit(1)
	}
	rPb := resp.GetRankingsCache()
	if rankingsCache == nil {
		rankingsCache = make(server.RankingsMap)
	}
	for year, rmPb := range rPb.GetCache() {
		if rankingsCache[year] == nil {
			rankingsCache[year] = make(map[string]server.RankingsData)
		}
		for e, r := range rmPb.GetEntry() {
			entry := server.RankingsData{
				ID:       r.GetID(),
				Year:     r.GetYear(),
				Bucket:   r.GetBucket(),
				Category: r.GetCategory(),
				Party:    r.GetParty(),
				Rankings: r.GetRankings(),
			}
			rankingsCache[year][e] = entry
		}
	}
	ytPb := resp.GetTotalsCache()
	if yrTotalsCache == nil {
		yrTotalsCache = make(server.YrTotalsMap)
	}
	for year, ymPb := range ytPb.GetCache() {
		if yrTotalsCache[year] == nil {
			yrTotalsCache[year] = make(map[string]server.YrTotalData)
		}
		for e, y := range ymPb.GetTotals() {
			entry := server.YrTotalData{
				ID:       y.GetID(),
				Year:     y.GetYear(),
				Category: y.GetCategory(),
				Party:    y.GetParty(),
				Total:    y.GetTotal(),
			}
			yrTotalsCache[year][e] = entry
		}
	}

	fmt.Println("building search cache...")
	searchDataCache, err = server.CreateSearchCache(rankingsCache)
	if err != nil {
		fmt.Println("failed to build search cache: ", err)
		os.Exit(1)
	}

	// create http server and handler functions
	go func() {
		fmt.Println("initializing http server...")
		srv := server.InitHTTPServer("localhost:8081")
		fmt.Printf("server address: %v\nread timeout: %v\nwrite timeout: %v\n",
			srv.Addr, srv.ReadTimeout, srv.WriteTimeout)
		fmt.Printf("listening at: '%v'...\n", srv.Addr)
		server.RegisterHandlers()
		log.Fatal(srv.ListenAndServe())
	}()

	// create gRPC server
	port := 9090
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("listening at port %d...\n", port)

	// Create the server TLS credentials
	fmt.Println("loading credentials...")
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		fmt.Printf("could not load TLS keys: %s\n", err)
		os.Exit(1)
	}

	cr := grpc.Creds(creds)
	sOpts = append(sOpts, cr)

	fmt.Println("registering new server...")
	grpcServer := grpc.NewServer()
	pb.RegisterViewServer(grpcServer, newRPCServer())
	fmt.Println("now serving!")
	grpcServer.Serve(lis)
}

/* Web Server functions */

// register new RPC server
func newRPCServer() *viewServer {
	return &viewServer{}
}

// SearchQuery finds a list of objects matching a given search term
func (s *viewServer) SearchQuery(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	// intitialize response object
	out := &pb.SearchResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tSearchQuery failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(err)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// make RPC call to Index service to find matching search results
	txt := in.GetText()
	resp, err := searchIndex(client, txt, hostname)
	if err != nil {
		fmt.Println("search index err: ", err)
		fmt.Println("msg: ", resp.GetMsg())
		out.Msg = resp.GetMsg()
		return out, err
	}
	results := []*pb.SearchResult{}
	for _, r := range resp.GetResults() {
		wrap := &pb.SearchResult{
			ID:       r.GetID(),
			Name:     r.GetName(),
			City:     r.GetCity(),
			State:    r.GetState(),
			Employer: r.GetEmployer(),
			Bucket:   r.GetBucket(),
			Years:    r.GetYears(),
		}
		results = append(results, wrap)
	}

	out.Results = results
	out.Msg = "SUCCESS"

	if len(out.Results) == 0 {
		out.Msg = "NO_RESULTS"
	}

	return out, nil
}

// LookupObjByID finds object summary data for a list of IDs
func (s *viewServer) LookupObjByID(ctx context.Context, in *pb.LookupRequest) (*pb.LookupResponse, error) {
	fmt.Println("called LookupObjByID...")
	// intitialize response object
	out := &pb.LookupResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupObjByID failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// find matching search results
	IDs := in.GetObjectIds()
	resp, err := lookupObjects(client, IDs, hostname) // rpc call to index service
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupObjByID failed: %v", time.Now(), err.Error())
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	results := []*pb.SearchResult{}

	// convert to proto.SearchResult message
	for _, r := range resp.GetResults() {
		wrap := &pb.SearchResult{
			ID:       r.GetID(),
			Name:     r.GetName(),
			City:     r.GetCity(),
			State:    r.GetState(),
			Employer: r.GetEmployer(),
			Bucket:   r.GetBucket(),
			Years:    r.GetYears(),
		}
		results = append(results, wrap)
	}

	out.Results = results
	out.Msg = "SUCCESS"

	fmt.Println("returning results...")
	return out, nil
}

// ViewRankings returns Rankings datasets from the in memory cache
func (s *viewServer) ViewRankings(ctx context.Context, in *pb.RankingsRequest) (*pb.RankingsResponse, error) {
	// intitialize response object
	out := &pb.RankingsResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewRankings failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// get object from cache
	year, bucket, cat, pty := in.GetYear(), in.GetBucket(), in.GetCategory(), in.GetParty()
	ID := fmt.Sprintf("%s-%s-%s-%s", year, bucket, cat, pty)
	rCache := rankingsCache[year][ID]

	// encode result
	res := pb.RankingsResult{
		ID:       rCache.ID,
		Year:     rCache.Year,
		Bucket:   rCache.Bucket,
		Category: rCache.Category,
		Party:    rCache.Party,
	}

	// sort IDs
	srt := util.SortMapObjectTotals(rCache.Rankings)
	rankings := []*pb.RankingEntry{}
	for _, e := range srt {
		sd := searchDataCache[e.ID]
		ranking := pb.RankingEntry{
			ID:     sd.ID,
			Name:   sd.Name,
			City:   sd.City,
			State:  sd.State,
			Years:  sd.Years,
			Amount: rCache.Rankings[sd.ID],
		}
		rankings = append(rankings, &ranking)
	}
	res.RankingsList = rankings
	out.Rankings = &res
	out.Msg = "SUCCESS"

	return out, nil
}

// ViewYrTotals returns Yearly Totals datasets from the in memory cache
func (s *viewServer) ViewYrTotals(ctx context.Context, in *pb.YrTotalRequest) (*pb.YrTotalResponse, error) {
	// intitialize response object
	out := &pb.YrTotalResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewYrTotals failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// get object from cache
	year, cat := in.GetYear(), in.GetCategory()
	totals := []*pb.YrTotalResult{}
	ptys := []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK"}

	for _, pty := range ptys {
		ID := fmt.Sprintf("%s-%s-%s", year, cat, pty)
		cache := yrTotalsCache[year][ID]
		// encode result
		res := pb.YrTotalResult{
			ID:       cache.ID,
			Year:     cache.Year,
			Category: cache.Category,
			Party:    cache.Party,
			Total:    cache.Total,
		}
		totals = append(totals, &res)
	}

	out.YearlyTotal = totals
	out.Msg = "SUCCESS"

	return out, nil
}

// ViewIndividual retrieves an Individual dataset from the Index service
func (s *viewServer) ViewIndividual(ctx context.Context, in *pb.GetIndvRequest) (*pb.GetIndvResponse, error) {
	fmt.Println("called ViewIndividual...")
	out := &pb.GetIndvResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts
	years := in.GetYears() // years requested
	if len(years) == 0 {
		err := "NO_YEAR_SET"
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	// rpc call to index service
	resp, err := lookupIndividual(client, out.ObjectID, hostname, years)
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Years = resp.GetYears() // years available

	indv := resp.GetIndividual()
	indvPb := pb.Individual{
		ID:            indv.GetID(),
		Name:          indv.GetName(),
		City:          indv.GetCity(),
		State:         indv.GetState(),
		Occupation:    indv.GetOccupation(),
		Employer:      indv.GetEmployer(),
		TotalOutAmt:   indv.GetTotalOutAmt(),
		TotalOutTxs:   indv.GetTotalOutTxs(),
		AvgTxOut:      indv.GetAvgTxOut(),
		TotalInAmt:    indv.GetTotalInAmt(),
		TotalInTxs:    indv.GetTotalInTxs(),
		AvgTxIn:       indv.GetAvgTxIn(),
		NetBalance:    indv.GetNetBalance(),
		RecipientsAmt: wrapTotals(indv.GetRecipientsAmt()),
		RecipientsTxs: indv.GetRecipientsTxs(),
		SendersAmt:    wrapTotals(indv.GetSendersAmt()),
		SendersTxs:    indv.GetSendersTxs(),
	}
	out.Individual = &indvPb
	out.Msg = "SUCCESS"

	return out, nil
}

// ViewCommittee retrieves a Committee dataset from the Index service
func (s *viewServer) ViewCommittee(ctx context.Context, in *pb.GetCmteRequest) (*pb.GetCmteResponse, error) {
	fmt.Println("called ViewCommittee...")
	out := &pb.GetCmteResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts
	years := in.GetYears() // years requested
	if len(years) == 0 {
		err := "NO_YEAR_SET"
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	// rpc call to index service
	resp, err := lookupCommittee(client, out.ObjectID, hostname, years)
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Years = resp.GetYears() // years available

	// Get object binary and return in respons
	cmte := resp.GetCommittee()
	cmtePb := pb.Committee{
		ID:           cmte.GetID(),
		Name:         cmte.GetName(),
		TresName:     cmte.GetTresName(),
		City:         cmte.GetCity(),
		State:        cmte.GetState(),
		Zip:          cmte.GetZip(),
		Designation:  cmte.GetDesignation(),
		Type:         cmte.GetType(),
		Party:        cmte.GetParty(),
		FilingFreq:   cmte.GetFilingFreq(),
		OrgType:      cmte.GetOrgType(),
		ConnectedOrg: cmte.GetConnectedOrg(),
		CandID:       cmte.GetCandID(),
	}
	out.Committee = &cmtePb

	cmteTx := resp.GetTxData()
	cmteTxPb := pb.CmteTxData{
		CmteID:                    cmteTx.GetCmteID(),
		CandID:                    cmteTx.GetCandID(),
		ContributionsInAmt:        cmteTx.GetContributionsInAmt(),
		ContributionsInTxs:        cmteTx.GetContributionsInTxs(),
		AvgContributionIn:         cmteTx.GetAvgContributionIn(),
		OtherReceiptsInAmt:        cmteTx.GetOtherReceiptsInAmt(),
		OtherReceiptsInTxs:        cmteTx.GetOtherReceiptsInTxs(),
		AvgOtherIn:                cmteTx.GetAvgOtherIn(),
		TotalIncomingAmt:          cmteTx.GetTotalIncomingAmt(),
		TotalIncomingTxs:          cmteTx.GetTotalIncomingTxs(),
		AvgIncoming:               cmteTx.GetAvgIncoming(),
		TransfersAmt:              cmteTx.GetTransfersAmt(),
		TransfersTxs:              cmteTx.GetTransfersTxs(),
		AvgTransfer:               cmteTx.GetAvgTransfer(),
		ExpendituresAmt:           cmteTx.GetExpendituresAmt(),
		ExpendituresTxs:           cmteTx.GetExpendituresTxs(),
		AvgExpenditure:            cmteTx.GetAvgExpenditure(),
		TotalOutgoingAmt:          cmteTx.GetTotalOutgoingAmt(),
		TotalOutgoingTxs:          cmteTx.GetTotalOutgoingTxs(),
		AvgOutgoing:               cmteTx.GetAvgOutgoing(),
		NetBalance:                cmteTx.GetNetBalance(),
		TopIndvContributorsAmt:    wrapTotals(cmteTx.GetTopIndvContributorsAmt()),
		TopIndvContributorsTxs:    cmteTx.GetTopIndvContributorsTxs(),
		TopCmteOrgContributorsAmt: wrapTotals(cmteTx.GetTopCmteOrgContributorsAmt()),
		TopCmteOrgContributorsTxs: cmteTx.GetTopCmteOrgContributorsTxs(),
		TransferRecsAmt:           wrapTotals(cmteTx.GetTransferRecsAmt()),
		TransferRecsTxs:           cmteTx.GetTransferRecsTxs(),
		TopExpRecipientsAmt:       wrapTotals(cmteTx.GetTopExpRecipientsAmt()),
		TopExpRecipientsTxs:       cmteTx.GetTopExpRecipientsTxs(),
	}
	out.TxData = &cmteTxPb
	out.Msg = "SUCCESS"

	return out, nil
}

// ViewCandidate retrieves a Candidate dataset from the Index service
func (s *viewServer) ViewCandidate(ctx context.Context, in *pb.GetCandRequest) (*pb.GetCandResponse, error) {
	fmt.Println("called ViewCandidate...")
	out := &pb.GetCandResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	years := in.GetYears() // years requested
	if len(years) == 0 {
		err := "NO_YEAR_SET"
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	// rpc call to index service
	resp, err := lookupCandidate(client, out.ObjectID, hostname, years)
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Years = resp.GetYears() // years available

	cand := resp.GetCandidate()
	candPb := pb.Candidate{
		ID:                   cand.ID,
		Name:                 cand.Name,
		Party:                cand.Party,
		OfficeState:          cand.OfficeState,
		Office:               cand.Office,
		PCC:                  cand.PCC,
		City:                 cand.City,
		State:                cand.State,
		Zip:                  cand.Zip,
		OtherAffiliates:      cand.OtherAffiliates,
		TransactionsList:     cand.TransactionsList,
		TotalDirectInAmt:     cand.TotalDirectInAmt,
		TotalDirectInTxs:     cand.TotalDirectInTxs,
		AvgDirectIn:          cand.AvgDirectIn,
		TotalDirectOutAmt:    cand.TotalDirectOutAmt,
		TotalDirectOutTxs:    cand.TotalDirectOutTxs,
		AvgDirectOut:         cand.AvgDirectOut,
		NetBalanceDirectTx:   cand.NetBalanceDirectTx,
		DirectRecipientsAmts: wrapTotals(cand.GetDirectRecipientsAmts()),
		DirectRecipientsTxs:  cand.DirectRecipientsTxs,
		DirectSendersAmts:    wrapTotals(cand.GetDirectSendersAmts()),
		DirectSendersTxs:     cand.DirectSendersTxs,
	}
	out.Candidate = &candPb

	out.Msg = "SUCCESS"

	return out, nil
}

// NoOp - One empty request, ZERO processing, followed by one empty response
func (s viewServer) NoOp(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// wrap ind.TotalsMap in pb.TotalsMap
func wrapTotals(m []*ind.TotalsMap) []*pb.TotalsMap {
	wrap := []*pb.TotalsMap{}
	for _, e := range m {
		wrap = append(wrap, &pb.TotalsMap{ID: e.GetID(), Total: e.GetTotal()})
	}
	return wrap
}

/* Index Client functions */
func getCaches(client ind.IndexClient, hostname string, opts ...grpc.CallOption) (*ind.GetCachesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createGetCachesRequest(hostname)
	resp, err := client.GetCaches(ctx, &req)
	if err != nil {
		fmt.Println("searchIndex (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createGetCachesRequest(hostname string) ind.GetCachesRequest {
	req := ind.GetCachesRequest{
		ServerID: hostname,
		Msg:      "new-search-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}
func searchIndex(client ind.IndexClient, query, hostname string, opts ...grpc.CallOption) (*ind.SearchIndexResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createSearchIndexRequest(query, hostname)
	resp, err := client.SearchIndex(ctx, &req)
	if err != nil {
		fmt.Println("search index err: ", err.Error())
		if err.Error() == "DeadlineExceeded" {
			msg := "DEADLINE_EXCEEDED"
			fmt.Println(msg)
			return resp, fmt.Errorf(msg)
		}
		fmt.Println("searchIndex (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createSearchIndexRequest(query, hostname string) ind.SearchIndexRequest {
	req := ind.SearchIndexRequest{
		UID:      "test007",
		ServerID: hostname,
		Text:     query,
		Msg:      "new-search-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}

func lookupObjects(client ind.IndexClient, ids []string, hostname string, opts ...grpc.CallOption) (*ind.LookupObjResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createLookupRequest(ids, hostname)
	resp, err := client.LookupObjects(ctx, &req)
	if err != nil {
		fmt.Println("lookupObjects (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createLookupRequest(ids []string, hostname string) ind.LookupObjRequest {
	req := ind.LookupObjRequest{
		UID:       "test007",
		ServerID:  hostname,
		ObjectIds: ids,
		Msg:       "new-search-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}

func lookupIndividual(client ind.IndexClient, ID, hostname string, years []string, opts ...grpc.CallOption) (*ind.LookupIndvResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createLookupIndvRequest(ID, hostname, years)
	resp, err := client.GetIndividual(ctx, &req)
	if err != nil {
		fmt.Println("lookupIndividual (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createLookupIndvRequest(ID, hostname string, years []string) ind.LookupIndvRequest {
	req := ind.LookupIndvRequest{
		UID:      "test007",
		ServerID: hostname,
		ObjectID: ID,
		Bucket:   "individuals",
		Years:    years,
		Msg:      "new-lookup-indv-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}

func lookupCommittee(client ind.IndexClient, ID, hostname string, years []string, opts ...grpc.CallOption) (*ind.LookupCmteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createLookupCmteRequest(ID, hostname, years)
	resp, err := client.GetCommittee(ctx, &req)
	if err != nil {
		fmt.Println("lookupCommittee (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createLookupCmteRequest(ID, hostname string, years []string) ind.LookupCmteRequest {
	req := ind.LookupCmteRequest{
		UID:      "test007",
		ServerID: hostname,
		ObjectID: ID,
		Bucket:   "committees",
		Years:    years,
		Msg:      "new-lookup-cmte-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}

func lookupCandidate(client ind.IndexClient, ID, hostname string, years []string, opts ...grpc.CallOption) (*ind.LookupCandResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createLookupCandRequest(ID, hostname, years)
	resp, err := client.GetCandidate(ctx, &req)
	if err != nil {
		fmt.Println("lookupCandidate (client) failed: ", err)
		return resp, err
	}

	return resp, nil
}

func createLookupCandRequest(ID, hostname string, years []string) ind.LookupCandRequest {
	req := ind.LookupCandRequest{
		UID:      "test007",
		ServerID: hostname,
		ObjectID: ID,
		Bucket:   "candidates",
		Years:    years,
		Msg:      "new-lookup-cand-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}
