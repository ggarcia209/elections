package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/elections/donations"
)

var cont = &donations.Contribution{
	TxType:     "40K",
	Occupation: "job",
	Employer:   "the trap",
	OtherID:    "P00033444",
	MemoCode:   "",
}

func main() {
	start := time.Now()
	bucket, incoming, memo := deriveTxTypes(cont)
	end := time.Since(start)
	fmt.Println(bucket, incoming, memo)
	fmt.Println("time elapsed: ", end)
	total := time.Since(start)
	fmt.Println("total time elapsed: ", total)

}

func deriveTxTypes(cont *donations.Contribution) (string, bool, bool) {
	// initialize return values
	incoming := false
	memo := false
	var bucket string

	// determine transaction type (incoming/outgoing)
	/* codes := map[string]int{
		"10": 10, "10J": 10, "11": 11, "11J": 11, "12": 12, "13": 13, "15": 15, "15C": 15, "15E": 15, "15F": 15,
		"15I": 15, "15J": 15, "15T": 15, "15Z": 15, "16": 16, "16C": 16, "16F": 16, "16G": 16, "16H": 16, "16J": 16, "16K": 16, "16L": 16, "16R": 16, "16U": 16,
		"17": 17, "18": 18, "19": 19, "20": 20, "21": 21, "22": 22,
		"23": 23, "24": 24, "28": 28, "29": 29, "30": 30,
		"31": 31, "32": 32, "40": 40, "41": 41, "42": 42,
	} */
	// re := regexp.MustCompile("[0-9]+")
	// strCode := strings.Join(re.FindAllString(cont.TxType, -1), "")
	numCode := cont.TxType
	if numCode < "20" || (numCode >= "30" && numCode < "33") {
		incoming = true
	}
	if cont.MemoCode == "X" {
		memo = true
	}

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(cont.OtherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		bucket = "committees"
	case idCode == "H" || idCode == "S" || idCode == "P":
		bucket = "candidates"
	default:
		if cont.Occupation == "" {
			bucket = "organizations"
		} else {
			bucket = "individuals"
		}
	}
	return bucket, incoming, memo
}
