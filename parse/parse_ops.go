// update: moved GetOffset and LogOffset outside of Parse logic, added offset as input to functions
//	offsets not logged until returned data is successfully persisted
package parse

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/elections/donations"
)

type mapOfFields map[int]string

// Parse25Contributions parses 25 rows of an individual contributions records file
// or 25 rows of a committee contribuutions file and creates a list of 25 Contribution
// objects to be processed, and a list of
// and a list of Individual (Donor) items to be stored (25 at most)

// check returned donors against donors saved to disk, update id_lookup
func Parse25Contributions(year string, file io.ReadSeeker, start int64) ([]*donations.Contribution, []interface{}, int64, error) {
	// seek to starting byte offset
	offset := int64(start)
	if _, err := file.Seek((offset), 0); err != nil {
		return nil, nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	IcQueue := []*donations.Contribution{}
	ObjQueue := []interface{}{}
	donors := make(map[string]map[string]interface{})

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// skip if MemoCode == "X" (not included in itemization totals)
		if fieldMap[18] == "X" {
			continue
		}

		// convert non-string values from original strings
		txDateFmt := "01/02/2006"
		txDate, err := time.Parse(txDateFmt, fmtDateStr(fieldMap[13]))
		if err != nil {
			fmt.Println(err)
		}
		txAmt, _ := strconv.ParseFloat(fieldMap[14], 32)
		fileNum, _ := strconv.Atoi(fieldMap[17])
		subID, _ := strconv.Atoi(fieldMap[20])

		// create object to be stored in database
		donation := &donations.Contribution{
			CmteID:     fieldMap[0],
			AmndtInd:   fieldMap[1],
			ReportType: fieldMap[2],
			TxPGI:      fieldMap[3],
			ImgNum:     fieldMap[4],
			TxType:     fieldMap[5],
			EntityType: fieldMap[6],
			Name:       fieldMap[7],
			City:       fieldMap[8],
			State:      fieldMap[9],
			Zip:        fieldMap[10],
			Employer:   fieldMap[11],
			Occupation: fieldMap[12],
			TxDate:     txDate,
			TxAmt:      float32(txAmt),
			OtherID:    fieldMap[15],
			TxID:       fieldMap[16],
			FileNum:    fileNum,
			MemoCode:   fieldMap[18],
			MemoText:   fieldMap[19],
			SubID:      subID,
		}

		// if donor is an Organization
		if donation.Occupation == "" {
			org, err := findOrg(year, donation, donors)
			if err != nil {
				fmt.Println("Parse25IndvCont failed: findPerson failed: ", err)
				return nil, nil, start, fmt.Errorf("Parse25IndvCont failed: findPerson failed: %v", err)
			}
			donation.OtherID = org.ID // set DonorID field after retreiving/creating donor

			// * temporary cache *
			if donors[donation.Name][donation.Zip] == nil { // new/not seen in current call
				// add new entry to map
				// DQueue = append(DQueue, donor)
				donors[donation.Name][donation.Zip] = org
			}
		} else { // donor is an Individual
			donor, err := findPerson(year, donation, donors)
			if err != nil {
				fmt.Println("Parse25IndvCont failed: findPerson failed: ", err)
				return nil, nil, start, fmt.Errorf("Parse25IndvCont failed: findPerson failed: %v", err)
			}
			donation.OtherID = donor.ID // set DonorID field after retreiving/creating donor

			// * temporary cache *
			job := donation.Employer + " - " + donation.Occupation
			if donors[donation.Name][job] == nil { // new/not seen in current call
				// add new entry to map
				// DQueue = append(DQueue, donor)
				donors[donation.Name][job] = donor
			}
		}

		// add donation to queue of items, stop at 25 items
		IcQueue = append(IcQueue, donation)
		if len(IcQueue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	for _, jobMap := range donors {
		for _, donor := range jobMap {
			if len(donor.(*donations.Organization).Transactions) == 0 || len(donor.(*donations.Individual).Transactions) == 0 { // new donor
				ObjQueue = append(ObjQueue, donor)
			}
		}
	}

	return IcQueue, ObjQueue, offset, nil
}

// Parse25Candidate parses 25 rows of a candidate records file
// and creates a list of 25 items to be stored in the database
func Parse25Candidate(file io.ReadSeeker, start int64) ([]*donations.Candidate, int64, error) {
	offset := int64(start)
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []*donations.Candidate{}

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// create object to be stored in database
		cand := &donations.Candidate{
			ID:          fieldMap[0],
			Name:        fieldMap[1],
			Party:       fieldMap[2],
			ElectnYr:    fieldMap[3],
			OfficeState: fieldMap[4],
			Office:      fieldMap[5],
			PCC:         fieldMap[9],
			City:        fieldMap[12],
			State:       fieldMap[13],
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, cand)
		if len(queue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return queue, offset, nil
}

// Parse25CmteLink parses 25 rows of a candidate-committee links records file
// and creates a list of 25 items to be stored in the database
func Parse25CmteLink(file io.ReadSeeker, start int64) ([]*donations.CmteLink, int64, error) {
	offset := int64(start)
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []*donations.CmteLink{}

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// convert values from string
		ey, _ := strconv.Atoi(fieldMap[1])

		// create object to be stored in database
		link := &donations.CmteLink{
			CandID:   fieldMap[0],
			ElectnYr: ey,
			CmteID:   fieldMap[3],
			CmteType: fieldMap[4],
			CmteDsgn: fieldMap[5],
			LinkID:   fieldMap[6],
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, link)
		if len(queue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return queue, offset, nil
}

// Parse25Committee parses 25 rows of a committee records file
// and creates a list of 25 items to be stored in the database
func Parse25Committee(file io.ReadSeeker, start int64) ([]*donations.Committee, int64, error) {
	offset := int64(start)
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []*donations.Committee{}

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// create object to be stored in database
		cmte := &donations.Committee{
			ID:           fieldMap[0],
			Name:         fieldMap[1],
			TresName:     fieldMap[2],
			City:         fieldMap[5],
			State:        fieldMap[6],
			Zip:          fieldMap[7],
			Designation:  fieldMap[8],
			Type:         fieldMap[9],
			Party:        fieldMap[10],
			FilingFreq:   fieldMap[11],
			OrgType:      fieldMap[12],
			ConnectedOrg: fieldMap[13],
			CandID:       fieldMap[14],
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, cmte)
		if len(queue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return queue, offset, nil
}

// Parse25Disbursements parses 25 rows of a disbursements records file,
// creates a list of 25 Disbursement items to be stored in the database,
// and a list of DisbRecipient items to be stored (25 at most)

// increment recID by len(DQueue) after each call
// check returned recipients against recipients saved to disk, update id_lookup
func Parse25Disbursements(year string, file io.ReadSeeker, start int64) ([]*donations.Disbursement, []*donations.Organization, int64, error) {
	// seek to starting byte offset
	offset := int64(start)
	if _, err := file.Seek((offset), 0); err != nil {
		return nil, nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	DQueue := []*donations.Disbursement{}
	RQueue := []*donations.Organization{}
	recs := make(map[string]map[string]interface{})

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// convert non-string values from original strings
		txDateFmt := "01/02/2006"
		txDate, err := time.Parse(txDateFmt, fieldMap[12])
		if err != nil {
			fmt.Println(err)
		}
		txAmt, _ := strconv.ParseFloat(fieldMap[13], 32)
		fileNum, _ := strconv.Atoi(fieldMap[22])
		subID, _ := strconv.Atoi(fieldMap[21])

		// create object to be stored in database
		disb := &donations.Disbursement{
			CmteID:       fieldMap[0],
			Name:         fieldMap[8],
			City:         fieldMap[9],
			State:        fieldMap[10],
			Zip:          fieldMap[11],
			TxDate:       txDate,
			TxAmt:        float32(txAmt),
			TxPGI:        fieldMap[14],
			Purpose:      fieldMap[15],
			Category:     fieldMap[16],
			CategoryDesc: fieldMap[17],
			MemoTxt:      fieldMap[19],
			EntityType:   fieldMap[20],
			SubID:        subID,
			FileNum:      fileNum,
			TxID:         fieldMap[23],
			BackRefTxID:  fieldMap[24],
		}

		rec, err := findOrg(year, disb, recs)
		// fmt.Println(donors)
		if err != nil {
			fmt.Println("Parse25Disbursements failed: findPerson failed: ", err)
			return nil, nil, start, fmt.Errorf("Parse25Disbursements failed: findPerson failed: %v", err)
		}

		disb.RecID = rec.ID // set RecID field after retreiving recipientt

		// * temporary cache *
		name := disb.Name
		zip := fmt.Sprintf("%s", disb.Zip)
		if recs[name][zip] == nil { // new
			// add new entry to map
			// RQueue = append(RQueue, rec)
			recs[name][zip] = rec
		}

		// add donation to queue of items, stop at 25 items
		DQueue = append(DQueue, disb)
		if len(DQueue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	for _, zipMap := range recs {
		for _, rec := range zipMap {
			if len(rec.(*donations.Organization).Transactions) == 0 { // new Organization
				RQueue = append(RQueue, rec.(*donations.Organization))
			}
		}
	}

	return DQueue, RQueue, offset, nil
}

// Parse25CmteFin parses 25 rows of a committee financial records file
// and creates a list of 25 items to be stored in the database
func Parse25CmteFin(file io.ReadSeeker, start int64) ([]*donations.CmteFinancials, int64, error) {
	// seek to starting byte offset
	if _, err := file.Seek(start, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []*donations.CmteFinancials{}

	// scanLines records the byte offset in order to recover from a failure
	offset := start
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		row := scanner.Text()

		// scan row and map field values
		fieldMap = scanRow(row, fieldMap)

		// convert non-string values from original strings
		tr, err := strconv.Atoi(fieldMap[5])
		if err != nil {
			tr = 0
		}
		tfa, err := strconv.Atoi(fieldMap[6])
		if err != nil {
			tfa = 0
		}
		ic, err := strconv.Atoi(fieldMap[7])
		if err != nil {
			ic = 0
		}
		oc, err := strconv.Atoi(fieldMap[8])
		if err != nil {
			oc = 0
		}
		cc, err := strconv.Atoi(fieldMap[9])
		if err != nil {
			cc = 0
		}
		tl, err := strconv.Atoi(fieldMap[10])
		if err != nil {
			tl = 0
		}
		td, err := strconv.Atoi(fieldMap[11])
		if err != nil {
			td = 0
		}
		tta, err := strconv.Atoi(fieldMap[12])
		if err != nil {
			tta = 0
		}
		ir, err := strconv.Atoi(fieldMap[13])
		if err != nil {
			ir = 0
		}
		or, err := strconv.Atoi(fieldMap[14])
		if err != nil {
			or = 0
		}
		lr, err := strconv.Atoi(fieldMap[15])
		if err != nil {
			lr = 0
		}
		cbop, err := strconv.Atoi(fieldMap[16])
		if err != nil {
			cbop = 0
		}
		ccop, err := strconv.Atoi(fieldMap[17])
		if err != nil {
			ccop = 0
		}
		do, err := strconv.Atoi(fieldMap[18])
		if err != nil {
			cc = 0
		}
		nft, err := strconv.Atoi(fieldMap[19])
		if err != nil {
			cc = 0
		}
		ctoc, err := strconv.Atoi(fieldMap[20])
		if err != nil {
			ctoc = 0
		}
		ie, err := strconv.Atoi(fieldMap[21])
		if err != nil {
			ie = 0
		}
		pe, err := strconv.Atoi(fieldMap[22])
		if err != nil {
			pe = 0
		}
		nfe, err := strconv.Atoi(fieldMap[23])
		if err != nil {
			nfe = 0
		}

		txDateFmt := "01/02/2006"
		ced, err := time.Parse(txDateFmt, fmtDateStr(fieldMap[13]))
		if err != nil {
			fmt.Println(err)
		}

		// create object to be stored in database
		fin := &donations.CmteFinancials{
			CmteID:          fieldMap[0],
			Type:            fieldMap[2],
			TotalReceipts:   tr,
			TxsFromAff:      tfa,
			IndvConts:       ic,
			OtherConts:      oc,
			CandCont:        cc,
			TotalLoans:      tl,
			TotalDisb:       td,
			TxToAff:         tta,
			IndvRefunds:     ir,
			OtherRefunds:    or,
			LoanRepay:       lr,
			CashBOP:         cbop,
			CashCOP:         ccop,
			DebtsOwed:       do,
			NonFedTxsRecvd:  nft,
			ContToOtherCmte: ctoc,
			IndExp:          ie,
			PartyExp:        pe,
			NonFedSharedExp: nfe,
			CovgEndDate:     ced,
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, fin)
		if len(queue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
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

// DEPRECATED
/*

// Parse25CmteCont parses 25 rows of a committee contributions records file
// and creates a list of 25 items to be stored in the database
func Parse25CmteCont(file io.ReadSeeker, start int64) ([]*donations.CmteContribution, int64, error) {
	// seek to starting byte offset
	offset := int64(start)
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []*donations.CmteContribution{}

	// scanLines records the byte offset in order to recover from a failure
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		offset += int64(advance)
		return
	}
	scanner.Split(scanLines)

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
		txAmt, _ := strconv.ParseFloat(fieldMap[14], 32)
		fileNum, _ := strconv.Atoi(fieldMap[17])
		subID, _ := strconv.Atoi(fieldMap[20])

		// create object to be stored in database
		donation := &donations.CmteContribution{
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
			TxAmt:      float32(txAmt),
			OtherID:    fieldMap[15],
			TxID:       fieldMap[16],
			FileNum:    fileNum,
			MemoCode:   fieldMap[18],
			MemoText:   fieldMap[19],
			SubID:      subID,
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, donation)
		if len(queue) == 25 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return queue, offset, nil
}
*/
