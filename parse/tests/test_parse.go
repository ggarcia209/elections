package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"projects/elections/donations"
	"strconv"
	"strings"
	"time"
)

type ICMap map[int]string
type CCMap map[int]interface{}

// Parse25Indv parses 25 rows of an individual contributions records file
// and creates a list of 25 items to be stored in the database
func Parse25Indv(file io.ReadSeeker, start int64) ([]*donations.IndvContribution, int64, error) {
	if _, err := file.Seek(start, 0); err != nil {
		return nil, start, err
	}
	scanner := bufio.NewScanner(file)
	fieldMap := make(ICMap)
	queue := []*donations.IndqvContribution{}

	// scanLines records the byte offset in order to recover from a failure
	offset := start
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	i := 0
	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)
		// convert non-string values from original strings
		txDateFmt := "01/02/2006"
		txDate, err := time.Parse(txDateFmt, fmtDateStr(fieldMap[13]))
		if err != nil {
			fmt.Println(err)
		}
		txAmt, _ := strconv.Atoi(fieldMap[14])
		fileNum, _ := strconv.Atoi(fieldMap[17])
		subId, _ := strconv.Atoi(fieldMap[20])

		donation := &donations.IndvContribution{
			CmteID:     fieldMap[0],
			AmndtInd:   fieldMap[1],
			ReportType: fieldMap[2],
			TxPGI:      fieldMap[3],
			TxType:     fieldMap[5],
			EntityType: fieldMap[6],
			Name:       fieldMap[7],
			City:       fieldMap[8],
			State:      fieldMap[9],
			Zip:        fieldMap[10],
			Employer:   fieldMap[11],
			Occupation: fieldMap[12],
			TxDate:     txDate,
			TxAmt:      txAmt,
			OtherID:    fieldMap[15],
			TxID:       fieldMap[16],
			FileNum:    fileNum,
			MemoCode:   fieldMap[18],
			MemoText:   fieldMap[19],
			SubID:      subId,
		}
		fmt.Println()
		fmt.Printf("%d:\t%v\n", i, donation)
		i++

		queue = append(queue, donation)
		if len(queue) == 10 {
			break
		}
	}
	return queue, offset, nil
}

/* Utility funcs */

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

const datapath = "~/go/src/projects/elections/parse/tests/test_indv.txt"

func main() {
	file, err := os.Open("test_indv.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/* scanner := bufio.NewScanner(file)
	m := make(map[int]interface{})

	// scanLines records the byte offset in order to recover from a failure
	offset := int64(0)
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		fmt.Println(row)
		scanRow(row, m)
	} */

	_, offset, _ := Parse25Indv(file, 0)
	fmt.Println(offset)
	_, offset, _ = Parse25Indv(file, offset)
	fmt.Println(offset)
	_, offset, _ = Parse25Indv(file, offset)
	fmt.Println(offset)
	_, offset, _ = Parse25Indv(file, offset)
	fmt.Println(offset)
}

func fmtDateStr(date string) string {
	newFmt := ""
	for i, c := range date {
		newFmt = newFmt + string(c)
		if i == 1 || i == 3 {
			newFmt = newFmt + "/"
		}
	}
	return newFmt
}
