package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elections/source/ui"
	"github.com/elections/source/util"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/elections/source/svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cert = "../cert/server.crt"

func main() {
	serverAddr := "127.0.0.1:9090"
	var opts []grpc.DialOption

	// Create the client TLS credentials
	fmt.Println("loading credentials...")
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		fmt.Printf("could not load tls cert: %s\n", err)
		os.Exit(1)
	}
	cr := grpc.WithTransportCredentials(creds)
	opts = append(opts, cr)

	fmt.Printf("dialing %s...\n", serverAddr)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure()) // change after testing
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewViewClient(conn)
	fmt.Println("connected to ", serverAddr)

	options := []string{"Search Data", "View Rankings", "View Yearly Totals", "Exit"}
	menu := ui.CreateMenu("client-view", options)
	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		switch menu.OptionsMap[ch] {
		case "Search Data":
			for {
				// get search query
				fmt.Println("enter search query: ")
				query := ui.GetQuery()
				resp, err := searchQuery(client, query)
				if err != nil {
					fmt.Println("err: ", err)
					fmt.Println("msg: ", resp.GetMsg())
					if resp.GetMsg() == "MAX_LENGTH" {
						fmt.Println("Too many results for query - Please refine the results by adding an additional search term.")
						continue
					} else {
						os.Exit(1)
					}
				}
				fmt.Println("printing results...")
				printResults(resp.Results)
				fmt.Println()
				fmt.Println("new search?")
				y := ui.Ask4confirm()
				if !y {
					fmt.Println("returning to menu...")
					break
				}
			}
		case "View Rankings":
			for {
				// get rankings
				req := createRankingsReq()
				resp, err := viewRankings(client, req)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("Got rankings...")
				res := resp.GetRankings()
				fmt.Println("ID: ", res.GetID())
				sorted := util.SortMapObjectTotals(res.GetRankingsList())
				for i, e := range sorted {
					fmt.Printf("%d)\t%s\t%.2f\n", i, e.ID, e.Total)
				}
				fmt.Println()
				fmt.Println("view new category?")
				y := ui.Ask4confirm()
				if !y {
					fmt.Println("returning to menu...")
					break
				}
			}
		case "View Yearly Totals":
			for {
				// get rankings
				req := createYrTotalReq()
				resp, err := viewYrTotals(client, req)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				res := resp.GetYearlyTotal()
				fmt.Println(res.GetID())
				fmt.Printf("Total: %.2f\n", res.GetTotal())
				fmt.Println()
				fmt.Println("view new category?")
				y := ui.Ask4confirm()
				if !y {
					fmt.Println("returning to menu...")
					break
				}
			}
		case "Exit":
			fmt.Println("Quitting...")
			os.Exit(1)
		}
	}
}

func searchQuery(client pb.ViewClient, query string, opts ...grpc.CallOption) (*pb.SearchResponse, error) {
	fmt.Println("Getting search results for query: ", query)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := createSearchRequest(query)
	resp, err := client.SearchQuery(ctx, &req)
	if err != nil {
		fmt.Println("searchQuery (client) failed: ", err)
		return resp, err
	}

	return resp, nil

}

func createSearchRequest(query string) pb.SearchRequest {
	req := pb.SearchRequest{
		UID:  "test007",
		Text: query,
		Msg:  "new-search-req",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createSearchRequest failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts

	return req
}

func printResults(res []*pb.SearchResult) {
	for i, r := range res {
		fmt.Printf("%d)\t%s - (%s, %s)\tYears: %v\n", i+1, r.Name, r.City, r.State, r.Years)
	}
}

func viewRankings(client pb.ViewClient, req pb.RankingsRequest, opts ...grpc.CallOption) (*pb.RankingsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Getting rankings...")
	resp, err := client.ViewRankings(ctx, &req)
	if err != nil {
		fmt.Println("viewRankings failed: ", err)
	}
	return resp, err
}

func createRankingsReq() pb.RankingsRequest {
	yr, bkt, cat, pty := getFields(true)
	req := pb.RankingsRequest{
		UID:      "test007",
		Year:     yr,
		Bucket:   bkt,
		Category: cat,
		Party:    pty,
		Msg:      "new-request",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createRankingsReq failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts
	return req
}

func viewYrTotals(client pb.ViewClient, req pb.YrTotalRequest, opts ...grpc.CallOption) (*pb.YrTotalResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Getting yearly total...")
	resp, err := client.ViewYrTotals(ctx, &req)
	if err != nil {
		fmt.Println("viewYrTotals failed: ", err)
	}
	return resp, err
}

func createYrTotalReq() pb.YrTotalRequest {
	yr, _, cat, pty := getFields(false)
	req := pb.YrTotalRequest{
		UID:      "test007",
		Year:     yr,
		Category: cat,
		Party:    pty,
		Msg:      "new request",
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println("createYrTotalReq failed: ", err)
		os.Exit(1)
	}
	req.Timestamp = ts
	return req
}

func getFields(rankings bool) (string, string, string, string) {
	year := ui.GetYear()
	buckets := []string{"individuals", "cmte_tx_data", "candidates"}
	cats := []string{"rec", "donor", "exp"}
	ptys := []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK"}
	var bucket string
	var category string
	var party string

	if rankings {
		menu := ui.CreateMenu("bkts", buckets)
		fmt.Println("choose bucket: ")
		ch, _ := ui.Ask4MenuChoice(menu)
		bucket = menu.OptionsMap[ch]
	}

	menu := ui.CreateMenu("cats", cats)
	fmt.Println("choose category: ")
	ch, _ := ui.Ask4MenuChoice(menu)
	category = menu.OptionsMap[ch]

	menu = ui.CreateMenu("ptys", ptys)
	fmt.Println("choose party: ")
	ch, _ = ui.Ask4MenuChoice(menu)
	party = menu.OptionsMap[ch]

	return year, bucket, category, party
}
