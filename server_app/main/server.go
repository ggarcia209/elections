package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/elections/source/server"
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
	crt = "../cert/server.crt"
	key = "../cert/server.key"
)

var rankingsCache server.RankingsMap
var yrTotalsCache server.YrTotalsMap
var searchDataCache server.SearchDataMap

func main() {
	fmt.Println("initializing disk cache...")
	server.InitServerDiskCache()
	fmt.Println("loading rankings and yearly total data from disk...")
	rankings, err := server.GetRankingsFromDisk()
	if err != nil {
		fmt.Println("failed to load rankings: ", err)
		os.Exit(1)
	}
	rankingsCache = rankings

	totals, err := server.GetYrTotalsFromDisk()
	if err != nil {
		fmt.Println("failed to load yearly totals: ", err)
		os.Exit(1)
	}
	yrTotalsCache = totals

	fmt.Println("initializing search cache...")
	sds, err := server.CreateSearchCache(rankingsCache)
	if err != nil {
		fmt.Println("failed to load yearly totals: ", err)
		os.Exit(1)
	}
	searchDataCache = sds

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

	var opts []grpc.ServerOption

	// Create the TLS credentials
	fmt.Println("loading credentials...")
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		fmt.Printf("could not load TLS keys: %s\n", err)
		os.Exit(1)
	}

	cr := grpc.Creds(creds)
	opts = append(opts, cr)

	fmt.Println("registering new server...")
	grpcServer := grpc.NewServer()
	pb.RegisterViewServer(grpcServer, newRPCServer())
	fmt.Println("now serving!")
	grpcServer.Serve(lis)

}

func newRPCServer() *viewServer {
	return &viewServer{}
}

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

	// find matching search results
	txt := in.GetText()
	common, err := server.SearchData(txt)
	if err != nil {
		fmt.Println(err)
		out.Msg = fmt.Sprintf("%s", err.Error())
		return out, err
	}
	sds, err := server.GetSearchResults(common, searchDataCache)
	if err != nil {
		errMsg := fmt.Errorf("%v\tSearchQuery failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(err)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}

	// convert to SearchResult message
	var results []*pb.SearchResult
	for _, sd := range sds {
		res := &pb.SearchResult{
			ID:       sd.ID,
			Bucket:   sd.Bucket,
			Name:     sd.Name,
			City:     sd.City,
			State:    sd.State,
			Employer: sd.Employer,
			Years:    sd.Years,
		}
		results = append(results, res)
	}
	out.Msg = "SUCCESS"
	out.Results = results
	if len(results) == 0 {
		out.Msg = "NO_RESULTS"
	}

	return out, nil
}

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
	srt := util.SortMapObjectTotals(rCache.Amts)
	rankings := []*pb.RankingEntry{}
	for _, e := range srt {
		sd := searchDataCache[e.ID]
		ranking := pb.RankingEntry{
			ID:     sd.ID,
			Name:   sd.Name,
			City:   sd.City,
			State:  sd.State,
			Years:  sd.Years,
			Amount: rCache.Amts[sd.ID],
		}
		rankings = append(rankings, &ranking)
	}
	res.RankingsList = rankings
	out.Rankings = &res
	out.Msg = "SUCCESS"

	return out, nil
}

// retrieve yearly total matching specified criteria
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

// retrieve object from cache/DynamoDB
func (s *viewServer) ViewIndividual(ctx context.Context, in *pb.GetIndvRequest) (*pb.GetIndvResponse, error) {
	fmt.Println("called ViewIndividualt...")
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

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	if len(years) == 0 {
		err := "NO_YEAR_SET"
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	year := years[0]

	obj, err := server.GetObjectFromDisk(year, in.GetObjectID(), in.GetBucket())
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	indv := obj.(server.Individual)
	recAmtsSrt := util.SortMapObjectTotals(indv.RecipientsAmt)
	recAmts := []*pb.TotalsMap{}
	for _, e := range recAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		recAmts = append(recAmts, entry)
	}
	senAmtsSrt := util.SortMapObjectTotals(indv.SendersAmt)
	senAmts := []*pb.TotalsMap{}
	for _, e := range senAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		senAmts = append(senAmts, entry)
	}

	indvPb := pb.Individual{
		ID:            indv.ID,
		Name:          indv.Name,
		City:          indv.City,
		State:         indv.State,
		Occupation:    indv.Occupation,
		Employer:      indv.Employer,
		TotalOutAmt:   indv.TotalOutAmt,
		TotalOutTxs:   indv.TotalOutTxs,
		AvgTxOut:      indv.AvgTxOut,
		TotalInAmt:    indv.TotalInAmt,
		TotalInTxs:    indv.TotalInTxs,
		AvgTxIn:       indv.AvgTxIn,
		NetBalance:    indv.NetBalance,
		RecipientsAmt: recAmts,
		RecipientsTxs: indv.RecipientsTxs,
		SendersAmt:    senAmts,
		SendersTxs:    indv.SendersTxs,
	}

	out.Individual = &indvPb
	out.Msg = "SUCCESS"

	return out, nil
}

// retrieve object from cache/DynamoDB
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

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	year := years[0]

	obj, err := server.GetObjectFromDisk(year, in.GetObjectID(), "committees")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	cmte := obj.(server.Committee)
	cmtePb := pb.Committee{
		ID:           cmte.ID,
		Name:         cmte.Name,
		TresName:     cmte.TresName,
		City:         cmte.City,
		State:        cmte.State,
		Zip:          cmte.Zip,
		Designation:  cmte.Designation,
		Type:         cmte.Type,
		Party:        cmte.Party,
		FilingFreq:   cmte.FilingFreq,
		OrgType:      cmte.OrgType,
		ConnectedOrg: cmte.ConnectedOrg,
		CandID:       cmte.CandID,
	}
	out.Committee = &cmtePb

	obj, err = server.GetObjectFromDisk(year, in.GetObjectID(), "cmte_tx_data")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	cmteTx := obj.(server.CmteTxData)
	indvAmtsSrt := util.SortMapObjectTotals(cmteTx.TopIndvContributorsAmt)
	indvAmts := []*pb.TotalsMap{}
	for _, e := range indvAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		indvAmts = append(indvAmts, entry)
	}
	cmteAmtsSrt := util.SortMapObjectTotals(cmteTx.TopCmteOrgContributorsAmt)
	cmteAmts := []*pb.TotalsMap{}
	for _, e := range cmteAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		cmteAmts = append(cmteAmts, entry)
	}
	trAmtsSrt := util.SortMapObjectTotals(cmteTx.TransferRecsAmt)
	trAmts := []*pb.TotalsMap{}
	for _, e := range trAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		trAmts = append(trAmts, entry)
	}
	expAmtsSrt := util.SortMapObjectTotals(cmteTx.TopExpRecipientsAmt)
	expAmts := []*pb.TotalsMap{}
	for _, e := range expAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		expAmts = append(expAmts, entry)
	}

	cmteTxPb := pb.CmteTxData{
		CmteID:                    cmteTx.CmteID,
		CandID:                    cmteTx.CandID,
		ContributionsInAmt:        cmteTx.ContributionsInAmt,
		ContributionsInTxs:        cmteTx.ContributionsInTxs,
		AvgContributionIn:         cmteTx.AvgContributionIn,
		OtherReceiptsInAmt:        cmteTx.OtherReceiptsInAmt,
		OtherReceiptsInTxs:        cmteTx.OtherReceiptsInTxs,
		AvgOtherIn:                cmteTx.AvgOtherIn,
		TotalIncomingAmt:          cmteTx.TotalIncomingAmt,
		TotalIncomingTxs:          cmteTx.TotalIncomingTxs,
		AvgIncoming:               cmteTx.AvgIncoming,
		TransfersAmt:              cmteTx.TransfersAmt,
		TransfersTxs:              cmteTx.TransfersTxs,
		AvgTransfer:               cmteTx.AvgTransfer,
		ExpendituresAmt:           cmteTx.ExpendituresAmt,
		ExpendituresTxs:           cmteTx.ExpendituresTxs,
		AvgExpenditure:            cmteTx.AvgExpenditure,
		TotalOutgoingAmt:          cmteTx.TotalOutgoingAmt,
		TotalOutgoingTxs:          cmteTx.TotalOutgoingTxs,
		AvgOutgoing:               cmteTx.AvgOutgoing,
		NetBalance:                cmteTx.NetBalance,
		TopIndvContributorsAmt:    indvAmts,
		TopIndvContributorsTxs:    cmteTx.TopIndvContributorsTxs,
		TopCmteOrgContributorsAmt: cmteAmts,
		TopCmteOrgContributorsTxs: cmteTx.TopCmteOrgContributorsTxs,
		TransferRecsAmt:           trAmts,
		TransferRecsTxs:           cmteTx.TransferRecsTxs,
		TopExpRecipientsAmt:       expAmts,
		TopExpRecipientsTxs:       cmteTx.TopExpRecipientsTxs,
	}
	out.TxData = &cmteTxPb

	obj, err = server.GetObjectFromDisk(year, in.GetObjectID(), "cmte_fin")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	cmteFin := obj.(server.CmteFinancials)
	if cmteFin.CmteID != "" { // object exists
		cmteFinPb := pb.CmteFinancials{
			CmteID:          cmteFin.CmteID,
			TotalReceipts:   cmteFin.TotalReceipts,
			TxsFromAff:      cmteFin.TxsFromAff,
			IndvConts:       cmteFin.IndvConts,
			OtherConts:      cmteFin.OtherConts,
			CandCont:        cmteFin.CandCont,
			TotalLoans:      cmteFin.TotalLoans,
			TotalDisb:       cmteFin.TotalDisb,
			TxToAff:         cmteFin.TxToAff,
			IndvRefunds:     cmteFin.IndvRefunds,
			OtherRefunds:    cmteFin.OtherRefunds,
			LoanRepay:       cmteFin.LoanRepay,
			CashBOP:         cmteFin.CashBOP,
			CashCOP:         cmteFin.CashCOP,
			DebtsOwed:       cmteFin.DebtsOwed,
			NonFedTxsRecvd:  cmteFin.NonFedTxsRecvd,
			ContToOtherCmte: cmteFin.ContToOtherCmte,
			IndExp:          cmteFin.IndExp,
			PartyExp:        cmteFin.PartyExp,
			NonFedSharedExp: cmteFin.NonFedSharedExp,
		}
		out.Financials = &cmteFinPb
		out.Msg = "SUCCESS"
	} else { // record does not exist for specified committee
		out.Msg = "SUCCESS_NO_FIN"
	}

	return out, nil
}

// retrieve object from cache/DynamoDB
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

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	year := years[0]

	obj, err := server.GetObjectFromDisk(year, in.GetObjectID(), "candidates")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	cand := obj.(server.Candidate)
	recAmtsSrt := util.SortMapObjectTotals(cand.DirectRecipientsAmts)
	recAmts := []*pb.TotalsMap{}
	for _, e := range recAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		recAmts = append(recAmts, entry)
	}
	senAmtsSrt := util.SortMapObjectTotals(cand.DirectSendersAmts)
	senAmts := []*pb.TotalsMap{}
	for _, e := range senAmtsSrt {
		entry := &pb.TotalsMap{ID: e.ID, Total: e.Total}
		senAmts = append(senAmts, entry)
	}

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
		DirectRecipientsAmts: recAmts,
		DirectRecipientsTxs:  cand.DirectRecipientsTxs,
		DirectSendersAmts:    senAmts,
		DirectSendersTxs:     cand.DirectSendersTxs,
	}
	out.Candidate = &candPb

	/* obj, err = server.GetObjectFromDisk(year, in.GetObjectID(), "cmpn_fin")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	cf := obj.(server.CmpnFinancials)
	if cf.CandID != "" { // object exists
		cfPb := pb.CmpnFinancials{
			CandID:         cf.CandID,
			Name:           cf.Name,
			PartyCd:        cf.PartyCd,
			Party:          cf.Party,
			TotalReceipts:  cf.TotalReceipts,
			TransFrAuth:    cf.TransFrAuth,
			TotalDisbsmts:  cf.TotalDisbsmts,
			TransToAuth:    cf.TransToAuth,
			COHBOP:         cf.COHBOP,
			COHCOP:         cf.COHCOP,
			CandConts:      cf.CandConts,
			CandLoans:      cf.CandLoans,
			OtherLoans:     cf.OtherLoans,
			CandLoanRepay:  cf.CandLoanRepay,
			OtherLoanRepay: cf.OtherLoanRepay,
			DebtsOwedBy:    cf.DebtsOwedBy,
			TotalIndvConts: cf.TotalIndvConts,
			SpecElection:   cf.SpecElection,
			PrimElection:   cf.PrimElection,
			RunElection:    cf.RunElection,
			GenElection:    cf.GenElection,
			GenElectionPct: cf.GenElectionPct,
			OtherCmteConts: cf.OtherCmteConts,
			PtyConts:       cf.PtyConts,
			IndvRefunds:    cf.IndvRefunds,
			CmteRefunds:    cf.CmteRefunds,
		}
		out.Financials = &cfPb
		out.Msg = "SUCCESS"
	} else { // record does not exist for specified committee
		out.Msg = "SUCCESS_NO_FIN"
	} */

	out.Msg = "SUCCESS"

	return out, nil
}

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
	sds, err := server.LookupByID(IDs)
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupObjByID failed: %v", time.Now(), err.Error())
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	results := []*pb.SearchResult{}

	// convert to SearchResult message
	for _, sd := range sds {
		res := &pb.SearchResult{
			ID:       sd.ID,
			Bucket:   sd.Bucket,
			Name:     sd.Name,
			City:     sd.City,
			State:    sd.State,
			Employer: sd.Employer,
			Years:    sd.Years,
		}
		results = append(results, res)
	}

	out.Results = results
	out.Msg = "SUCCESS"

	fmt.Println("returning results...")
	return out, nil
}

// One empty request, ZERO processing, followed by one empty response
// (minimum effort to do message serialization).
func (s viewServer) NoOp(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
