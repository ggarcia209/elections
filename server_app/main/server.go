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
	fmt.Println("called method SearchQuery...")
	// intitialize response object
	out := &pb.SearchResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("SearchQuery failed: %v\tUID: %s\tTimeStamp: %v", err, out.UID, out.Timestamp)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// find matching search results
	txt := in.GetText()
	fmt.Printf("Getting search results for '%s'...\n", txt)
	sds, err := server.SearchData(txt)
	if err != nil {
		out.Msg = err.Error()
		fmt.Println("SearchQuery (server) error: ", err.Error())
		return out, err
	}

	// convert to SearchResult message
	var results []*pb.SearchResult
	for _, sd := range sds {
		res := &pb.SearchResult{
			ID:    sd.ID,
			Name:  sd.Name,
			City:  sd.City,
			State: sd.State,
			Years: sd.Years,
		}
		results = append(results, res)
	}
	out.Msg = "SUCCESS"
	out.Results = results
	if len(results) == 0 {
		out.Msg = "NO_RESULTS"
	}
	fmt.Println("Returning results...")

	return out, nil
}

var rankingsCache server.RankingsMap
var yrTotalsCache server.YrTotalsMap

func (s *viewServer) ViewRankings(ctx context.Context, in *pb.RankingsRequest) (*pb.RankingsResponse, error) {
	fmt.Println("called method ViewRankings")
	// intitialize response object
	fmt.Println("creating response...")
	out := &pb.RankingsResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("SearchQuery failed: %v\tUID: %s\tTimeStamp: %v", err, out.UID, out.Timestamp)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// get object from cache
	fmt.Println("getting object from cache...")
	year, bucket, cat, pty := in.GetYear(), in.GetBucket(), in.GetCategory(), in.GetParty()
	ID := fmt.Sprintf("%s-%s-%s-%s", year, bucket, cat, pty)
	fmt.Println("ID in: ", ID)
	cache := rankingsCache[year][ID]

	// encode result
	fmt.Println("encoding result...")
	res := pb.RankingsResult{
		ID:           cache.ID,
		Year:         cache.Year,
		Bucket:       cache.Bucket,
		Category:     cache.Category,
		Party:        cache.Party,
		RankingsList: cache.Amts,
	}
	fmt.Println("ID: ", res.ID)

	out.Rankings = &res

	return out, nil
}

// retrieve yearly total matching specified criteria
func (s *viewServer) ViewYrTotals(ctx context.Context, in *pb.YrTotalRequest) (*pb.YrTotalResponse, error) {
	fmt.Println("called method ViewYrTotals")
	// intitialize response object
	fmt.Println("creating response...")
	out := &pb.YrTotalResponse{
		UID: in.GetUID(),
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		errMsg := fmt.Errorf("ViewYrTotals failed: %v\tUID: %s\tTimeStamp: %v", err, out.UID, out.Timestamp)
		out.Msg = fmt.Sprintf("%s", errMsg)
		return out, errMsg
	}
	out.Timestamp = ts

	// get object from cache
	fmt.Println("getting object from cache...")
	year, cat, pty := in.GetYear(), in.GetCategory(), in.GetParty()
	ID := fmt.Sprintf("%s-%s-%s", year, cat, pty)
	cache := yrTotalsCache[year][ID]

	// encode result
	fmt.Println("encoding result...")
	res := pb.YrTotalResult{
		ID:       cache.ID,
		Year:     cache.Year,
		Category: cache.Category,
		Party:    cache.Party,
		Total:    cache.Total,
	}

	out.YearlyTotal = &res

	return out, nil
}

// retrieve object from cache/DynamoDB
func (s *viewServer) ViewObject(ctx context.Context, in *pb.GetObjRequest) (*pb.GetObjResponse, error) {
	return nil, nil
}

// One empty request, ZERO processing, followed by one empty response
// (minimum effort to do message serialization).
func (s viewServer) NoOp(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
