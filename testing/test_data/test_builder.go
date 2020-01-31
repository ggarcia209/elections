package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Entry struct {
	ID     string
	LinkID string
	Other  string
}

// input file paths
const candIn = "test_cands.txt"
const cmteIn = "../../raw_data/cmte/2018/cm.txt"
const icIn = "../../raw_data/indv_cont/2018/itcont.txt"
const ccIn = "../../raw_data/cmte_cont/2018/itoth.txt"
const disbIn = "../../raw_data/disb/2018/oppexp.txt"

// output file paths
const cmteOut = "test_cmtes.txt"
const icOut = "test_ics.txt"
const ccOut = "test_ccs.txt"
const disbOut = "test_disbs.txt"

var seen = make(map[string]bool)

func main() {
	// start with candIDs
	candIDs, err := scanCand(candIn)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanCand done")
	fmt.Println()

	// get linked cmte's
	cmteIDs, err := scanCmtes(cmteIn, cmteOut, candIDs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanCmtes done")
	fmt.Println()

	// get cmte conts
	otherIDs, err := scanCmteConts(ccIn, ccOut, cmteIDs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanCmteConts done")

	// get cmte's linked by cmte conts
	_, err = scanCmtes(cmteIn, cmteOut, otherIDs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanCmtes done")
	fmt.Println()

	// get indv conts
	err = scanIndvConts(icIn, icOut, cmteIDs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanIndvConts done")
	fmt.Println()

	// get disbursements
	err = scanDisbs(disbIn, disbOut, cmteIDs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("scanDisbs done")
	fmt.Println()

}

/* FILE SCAN & WRITE */

// get cmteIDs linked to candidates from test_cands.txt
func scanCand(inPath string) ([]Entry, error) {
	// open input file
	file, err := os.Open(inPath)
	if err != nil {
		fmt.Println("scanCand failed: ", err)
		return nil, fmt.Errorf("scanCand failed: %v", err)
	}

	candIDs := []Entry{} // default return value
	itemMap := make(map[int]string)

	// scan lines
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	fmt.Println("----- CAND IDs -----")
	i := 0
	for scanner.Scan() {
		row := scanner.Text()
		itemMap = scanRow(row, itemMap)
		e := Entry{LinkID: itemMap[9]} // CandID / PCC ID
		//fmt.Println(e)
		candIDs = append(candIDs, e)
		i++
	}
	fmt.Println("initial entries: ", i)
	fmt.Println("----- END CAND IDs -----")
	fmt.Println()

	return candIDs, nil
}

// get cmte records by CandID and PCC ID
// return cmteID and linked CandID
func scanCmtes(inPath, outPath string, linkIDs []Entry) ([]Entry, error) {
	// open input file
	input, err := os.Open(inPath)
	if err != nil {
		fmt.Println("scanCmtes failed: ", err)
		return nil, fmt.Errorf("scanCmtes failed: %v", err)
	}
	defer input.Close()

	// open output file
	output, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("scanCmtes failed: ", err)
		return nil, fmt.Errorf("scanCmtes failed: %v", err)
	}
	defer output.Close()

	// convert linkId's to set map
	set := sliceToMap(linkIDs)

	// initialize item map and cmteID slice
	itemMap := make(map[int]string)
	cmteIDs := []Entry{}

	// scan lines
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	fmt.Println("----- CMTE IDs -----")
	i := 0
	for scanner.Scan() {
		row := scanner.Text()

		itemMap = scanRow(row, itemMap)
		// search set for corresponding id or CandID)
		if set[itemMap[0]] == true { // cmte ID (rec) / cand ID
			if seen[itemMap[0]] == true {
				continue
			}
			seen[itemMap[0]] = true

			// if linked, append to cmteIDs and write out to test file
			e := Entry{ID: itemMap[0]} // cmteID / cand ID
			fmt.Println(e)
			cmteIDs = append(cmteIDs, e)
			if _, err := output.WriteString(fmt.Sprintf("%v\n", row)); err != nil {
				fmt.Println("scanCmtes failed: ", err)
				return nil, fmt.Errorf("scanCmtes failed: %v", err)
			}
			i++
		}
	}
	fmt.Println("lines written: ", i)
	fmt.Println("----- END CMTE IDs -----")
	fmt.Println()

	return cmteIDs, nil
}

// get cmte_cont records by cmteID and linked CandID
// return otherIDs - pass to 2nd round of scanCmtes
func scanCmteConts(inPath, outPath string, linkIDs []Entry) ([]Entry, error) {
	// open input file
	input, err := os.Open(inPath)
	if err != nil {
		fmt.Println("scanCmteConts failed: ", err)
		return nil, fmt.Errorf("scanCmteConts failed: %v", err)
	}
	defer input.Close()

	// open output file
	output, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("scanCmteConts failed: ", err)
		return nil, fmt.Errorf("scanCmteConts failed: %v", err)
	}
	defer output.Close()

	// convert linkId's to set map
	set := sliceToMap(linkIDs)

	// initialize item map and cmteID slice
	itemMap := make(map[int]string)
	otherIDs := []Entry{}

	// scan lines
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	fmt.Println("----- CMTE CONTS - OTHER IDs ------")
	i := 0
	for scanner.Scan() {
		row := scanner.Text()

		itemMap = scanRow(row, itemMap)
		// search set for corresponding id
		if set[itemMap[0]] == true { // cmteID in set
			// if linked, write out to test file & add otherID to cmteID slice
			e := Entry{ID: itemMap[15]} // cmte contributor IDs
			otherIDs = append(otherIDs, e)
			fmt.Println(itemMap[16], e)
			if _, err := output.WriteString(fmt.Sprintf("%v\n", row)); err != nil {
				fmt.Println("scanCmteConts failed: ", err)
				return nil, fmt.Errorf("scanCmteConts failed: %v", err)
			}
			i++
		}

	}
	fmt.Println("lines written: ", i)
	fmt.Println("----- END CMTE CONTs - OTHER IDs-----")
	fmt.Println()

	return otherIDs, nil
}

func scanIndvConts(inPath, outPath string, linkIDs []Entry) error {
	// open input file
	input, err := os.Open(inPath)
	if err != nil {
		fmt.Println("scanIndvConts failed: ", err)
		return fmt.Errorf("scanIndvConts failed: %v", err)
	}
	defer input.Close()

	// open output file
	output, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("scanIndvConts failed: ", err)
		return fmt.Errorf("scanIndvConts failed: %v", err)
	}
	defer output.Close()

	// convert linkId's to set map
	set := sliceToMap(linkIDs)

	// initialize item map and cmteID slice
	itemMap := make(map[int]string)

	// scan lines
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		row := scanner.Text()

		itemMap = scanRow(row, itemMap)
		// search set for corresponding id (ConnectedOrg or CandID)
		if set[itemMap[0]] == true {
			// if linked, write out to test file
			if _, err := output.WriteString(fmt.Sprintf("%v\n", row)); err != nil {
				fmt.Println("scanIndvConts failed: ", err)
				return fmt.Errorf("scanIndvConts failed: %v", err)
			}
		}
	}

	return nil
}

func scanDisbs(inPath, outPath string, linkIDs []Entry) error {
	// open input file
	input, err := os.Open(inPath)
	if err != nil {
		fmt.Println("scanDisbs failed: ", err)
		return fmt.Errorf("scanDisbs failed: %v", err)
	}
	defer input.Close()

	// open output file
	output, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("scanDisbs failed: ", err)
		return fmt.Errorf("scanDisbs failed: %v", err)
	}
	defer output.Close()

	// convert linkId's to set map
	set := sliceToMap(linkIDs)

	// initialize item map and cmteID slice
	itemMap := make(map[int]string)

	// scan lines
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		row := scanner.Text()

		itemMap = scanRow(row, itemMap)
		// search set for corresponding id (ConnectedOrg or CandID)
		if set[itemMap[0]] == true {
			// if linked, write out to test file
			if _, err := output.WriteString(fmt.Sprintf("%v\n", row)); err != nil {
				fmt.Println("scanDisbs failed: ", err)
				return fmt.Errorf("scanDisbs failed: %v", err)
			}
		}
	}

	return nil
}

/* UTILITY */

func sliceToMap(es []Entry) map[string]bool {
	m := make(map[string]bool)
	for _, e := range es {
		if e.ID == "" {
			m[e.ID] = false
		} else {
			m[e.ID] = true
		}
		if e.LinkID == "" {
			m[e.LinkID] = false
		} else {
			m[e.LinkID] = true
		}
		if e.Other == "" {
			m[e.Other] = false
		} else {
			m[e.Other] = true
		}
	}
	return m
}

func scanIDs(row string) string {
	scanner := bufio.NewScanner(strings.NewReader(row))
	scanner.Split(rowSplit)

	var id string
	for scanner.Scan() {
		id = scanner.Text()
		break
	}
	return id
}

func scanRow(row string, m map[int]string) map[int]string {
	scanner := bufio.NewScanner(strings.NewReader(row))
	scanner.Split(rowSplit)
	i := 0
	for scanner.Scan() {
		m[i] = scanner.Text()
		i++
	}
	return m
}

func rowSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a "|"
	if i := strings.Index(string(data), "|"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}
