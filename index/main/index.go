// The index service is primarily responsible for information retrieval.
// Search query results (list of common IDs) are retrieved from BoltDB local storage.
// All other data, inlcuding both summary and complete object data, is retrieved from DynamoDB.
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
	pb "github.com/elections/source/svc/index"
	"github.com/elections/source/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// indexServer implements the IndexServer gRPC interface
type indexServer struct {
	pb.UnimplementedIndexServer
	mu sync.Mutex
}

var (
	crt = "../cert/server.crt"
	key = "../cert/server.key"
)

var rankingsCache server.RankingsMap
var yrTotalsCache server.YrTotalsMap
var searchDataCache server.SearchDataMap

var database *dynamo.DbInfo
var metadata *server.IndexData

func main() {
	var err error
	fmt.Println("initializing disk cache...")
	server.InitServerDiskCache()
	database, err = server.InitDynamo()
	if err != nil {
		fmt.Println("failed to load rankings: ", err)
		os.Exit(1)
	}
	fmt.Println("Retrieving index metadata...")
	metadata, err = server.GetIndexData()
	if err != nil {
		fmt.Println("failed to load rankings: ", err)
		os.Exit(1)
	}
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

	// create gRPC server
	port := 9092
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
	pb.RegisterIndexServer(grpcServer, newRPCServer())
	fmt.Println("now serving!")
	grpcServer.Serve(lis)

}

func newRPCServer() *indexServer {
	return &indexServer{}
}

func (s *indexServer) GetCaches(ctx context.Context, in *pb.GetCachesRequest) (*pb.GetCachesResponse, error) {
	// intitialize response object
	out := &pb.GetCachesResponse{
		ServerID: in.GetServerID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tSearchQuery failed: %v\n", time.Now(), err)
		fmt.Println(err)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	rCache := &pb.RankingsCache{}
	for year, rankings := range rankingsCache {
		entry := &pb.RankingsMap{}
		rMap := make(map[string]*pb.Rankings)
		for ID, data := range rankings {
			rMap[ID] = &pb.Rankings{
				ID:       data.ID,
				Year:     data.Year,
				Bucket:   data.Bucket,
				Category: data.Category,
				Party:    data.Party,
				Rankings: data.Rankings,
			}
		}
		entry.Entry = rMap
		if rCache.Cache == nil {
			rCache.Cache = make(map[string]*pb.RankingsMap)
		}
		rCache.Cache[year] = entry
	}
	ytCache := &pb.TotalsCache{}
	for year, totals := range yrTotalsCache {
		entry := &pb.YrTotalsMap{}
		ytMap := make(map[string]*pb.Totals)
		for ID, data := range totals {
			ytMap[ID] = &pb.Totals{
				ID:       data.ID,
				Year:     data.Year,
				Category: data.Category,
				Party:    data.Party,
				Total:    data.Total,
			}
		}
		entry.Totals = ytMap
		if ytCache.Cache == nil {
			ytCache.Cache = make(map[string]*pb.YrTotalsMap)
		}
		ytCache.Cache[year] = entry
	}

	out.RankingsCache = rCache
	out.TotalsCache = ytCache
	out.Msg = "SUCCESS"

	return out, nil
}

func (s *indexServer) SearchIndex(ctx context.Context, in *pb.SearchIndexRequest) (*pb.SearchIndexResponse, error) {
	// intitialize response object
	out := &pb.SearchIndexResponse{
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
	fmt.Println("SEARCH QUERY: ", txt)
	common, err := server.SearchData(metadata, txt)
	if err != nil {
		fmt.Println(err)
		out.Msg = fmt.Sprintf("%s", err.Error())
		return out, err
	}
	sds, err := server.GetSearchResults(database, common, searchDataCache)
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

func (s *indexServer) LookupObjects(ctx context.Context, in *pb.LookupObjRequest) (*pb.LookupObjResponse, error) {
	fmt.Println("called LookupObjByID...")
	// intitialize response object
	out := &pb.LookupObjResponse{
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
	sds, err := server.LookupByID(database, IDs)
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

// retrieve object from cache/DynamoDB
func (s *indexServer) GetIndividual(ctx context.Context, in *pb.LookupIndvRequest) (*pb.LookupIndvResponse, error) {
	fmt.Println("called LookupIndividual...")
	out := &pb.LookupIndvResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
		Years:    in.GetYears(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	sd, err := server.LookupByID(database, []string{out.ObjectID})
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Years = sd[0].Years

	out.Timestamp = ts

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	if len(years) == 0 {
		err := "NO_YEAR_SET"
		errMsg := fmt.Errorf("%v\tLookupIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}

	q := server.CreateQueryFromSearchData(sd[0])
	objs, err := server.GetObjectFromDynamo(database, q, sd[0].Bucket, years)
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupIndividual failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	indv := objs[0].(server.Individual)

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
func (s *indexServer) GetCommittee(ctx context.Context, in *pb.LookupCmteRequest) (*pb.LookupCmteResponse, error) {
	fmt.Println("called LookupCommittee...")
	out := &pb.LookupCmteResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	sd, err := server.LookupByID(database, []string{out.ObjectID})
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	if len(sd) == 0 {
		msg := "ITEM_NOT_FOUND"
		out.Msg = msg
		fmt.Println(msg)
		return out, fmt.Errorf(msg)
	}
	out.Years = sd[0].Years

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	query := server.CreateQueryFromSearchData(sd[0])
	st := time.Now()
	objs, err := server.GetObjectFromDynamo(database, query, sd[0].Bucket, years)
	if err != nil {
		if err.Error() == "TABLE_NOT_FOUND" {
			out.Msg = err.Error()
			return out, err
		}
		errMsg := fmt.Errorf("%v\tLookupCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	fmt.Println("get obj from dynamo time: ", time.Since(st))
	cmte := objs[0].(server.Committee)

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

	sd[0].Bucket = "cmte_tx_data"
	query = server.CreateQueryFromSearchData(sd[0])
	st = time.Now()
	objs, err = server.GetObjectFromDynamo(database, query, sd[0].Bucket, years)
	// obj, err = server.GetObjectFromDisk(year, in.GetObjectID(), "cmte_tx_data")
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	fmt.Println("get obj from dynamo time: ", time.Since(st))
	cmteTx := objs[0].(server.CmteTxData)
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

	/* obj, err = server.GetObjectFromDisk(year, in.GetObjectID(), "cmte_fin")
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupCommittee failed: %v\tUID: %s", time.Now(), err, out.UID)
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
	} */
	out.Msg = "SUCCESS"

	return out, nil
}

// retrieve object from cache/DynamoDB
func (s *indexServer) GetCandidate(ctx context.Context, in *pb.LookupCandRequest) (*pb.LookupCandResponse, error) {
	fmt.Println("called LookupCandidate...")
	out := &pb.LookupCandResponse{
		UID:      in.GetUID(),
		ObjectID: in.GetObjectID(),
		Bucket:   in.GetBucket(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("%v\tLookupCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	sd, err := server.LookupByID(database, []string{out.ObjectID})
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Years = sd[0].Years

	// Get object binary and return in response
	/* support for multiple years and aggregated datasets will be available in future version */
	years := in.GetYears()
	query := server.CreateQueryFromSearchData(sd[0])
	st := time.Now()
	obj, err := server.GetObjectFromDynamo(database, query, sd[0].Bucket, years)
	// obj, err := server.GetObjectFromDisk(year, in.GetObjectID(), "candidates")
	if err != nil {
		errMsg := fmt.Errorf("%v\tViewCandidate failed: %v\tUID: %s", time.Now(), err, out.UID)
		fmt.Println(errMsg)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	fmt.Println("get obj from dynamo time: ", time.Since(st))
	cand := obj[0].(server.Candidate)
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
