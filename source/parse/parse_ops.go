// Package parse contains operations for scanning the bulk data files and returning a list of objects derived from each file.
package parse

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/elections/source/donations"
)

type mapOfFields map[int]string

// ScanCandidates scans 100 lines of a candidates file
// and returns 100 Candidate objects per call
func ScanCandidates(file io.ReadSeeker, start int64) ([]interface{}, int64, error) {
	offset := int64(start)
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []interface{}{}

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
		if len(queue) == 10000 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return queue, offset, nil
}

// ScanCommittees scans 100 lines of a committees file
// and returns 100 Committee & corresponding CmteTxData objects per call
func ScanCommittees(file io.ReadSeeker, start int64) ([]interface{}, []interface{}, int64, error) {
	offset := int64(start)
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []interface{}{}
	dataQueue := []interface{}{}

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

		if cmte.Party == "" {
			cmte.Party = "UNK"
		}

		// initialize corresponding CmteTxData object
		txData := &donations.CmteTxData{
			CmteID: cmte.ID,
			CandID: cmte.CandID,
			Party:  cmte.Party,
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, cmte)
		dataQueue = append(dataQueue, txData)
		if len(queue) == 10000 {
			break
		}

		fieldMap = make(mapOfFields)
	}

	return queue, dataQueue, offset, nil
}

// ScanCmpnFin scans 100 lines of a committee financials file
// and returns 100 CmteFinancials objects per call
func ScanCmpnFin(file io.ReadSeeker, start int64) ([]interface{}, int64, error) {
	offset := start
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, offset, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []interface{}{}

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
		tr, err := strconv.ParseFloat(fieldMap[5], 32)
		if err != nil {
			tr = 0
		}
		tfa, err := strconv.ParseFloat(fieldMap[6], 32)
		if err != nil {
			tfa = 0
		}
		td, err := strconv.ParseFloat(fieldMap[7], 32)
		if err != nil {
			td = 0
		}
		tta, err := strconv.ParseFloat(fieldMap[8], 32)
		if err != nil {
			tta = 0
		}
		bop, err := strconv.ParseFloat(fieldMap[9], 32)
		if err != nil {
			bop = 0
		}
		cop, err := strconv.ParseFloat(fieldMap[10], 32)
		if err != nil {
			cop = 0
		}
		cc, err := strconv.ParseFloat(fieldMap[11], 32)
		if err != nil {
			cc = 0
		}
		cl, err := strconv.ParseFloat(fieldMap[12], 32)
		if err != nil {
			cl = 0
		}
		ol, err := strconv.ParseFloat(fieldMap[13], 32)
		if err != nil {
			ol = 0
		}
		clr, err := strconv.ParseFloat(fieldMap[14], 32)
		if err != nil {
			clr = 0
		}
		olr, err := strconv.ParseFloat(fieldMap[15], 32)
		if err != nil {
			olr = 0
		}
		dob, err := strconv.ParseFloat(fieldMap[16], 32)
		if err != nil {
			dob = 0
		}
		tic, err := strconv.ParseFloat(fieldMap[17], 32)
		if err != nil {
			tic = 0
		}
		gep, err := strconv.ParseFloat(fieldMap[24], 32)
		if err != nil {
			gep = 0
		}

		pcc, err := strconv.ParseFloat(fieldMap[25], 32)
		if err != nil {
			pcc = 0
		}

		ppc, err := strconv.ParseFloat(fieldMap[26], 32)
		if err != nil {
			ppc = 0
		}

		ir, err := strconv.ParseFloat(fieldMap[28], 32)
		if err != nil {
			ir = 0
		}

		cr, err := strconv.ParseFloat(fieldMap[29], 32)
		if err != nil {
			cr = 0
		}

		// txDateFmt := "01/02/2006"
		// ced, err := time.Parse(txDateFmt, fieldMap[24])
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// create object to be stored in database
		fin := &donations.CmpnFinancials{
			CandID:         fieldMap[0],
			Name:           fieldMap[1],
			PartyCd:        fieldMap[3],
			Party:          fieldMap[4],
			TotalReceipts:  float32(tr),
			TransFrAuth:    float32(tfa),
			TotalDisbsmts:  float32(td),
			TransToAuth:    float32(tta),
			COHBOP:         float32(bop),
			COHCOP:         float32(cop),
			CandConts:      float32(cc),
			CandLoans:      float32(cl),
			OtherLoans:     float32(ol),
			CandLoanRepay:  float32(clr),
			OtherLoanRepay: float32(olr),
			DebtsOwedBy:    float32(dob),
			TotalIndvConts: float32(tic),
			OfficeState:    fieldMap[18],
			OfficeDistrict: fieldMap[19],
			SpecElection:   fieldMap[20],
			PrimElection:   fieldMap[21],
			RunElection:    fieldMap[22],
			GenElection:    fieldMap[23],
			GenElectionPct: float32(gep),
			OtherCmteConts: float32(pcc),
			PtyConts:       float32(ppc),
			IndvRefunds:    float32(ir),
			CmteRefunds:    float32(cr),
			// CovgEndDate:     ced,
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, fin)
		if len(queue) == 10000 {
			break
		}
		fieldMap = make(mapOfFields)
	}
	return queue, offset, nil
}

// ScanCmteFin scans 100 lines of a committee financials file
// and returns 100 CmteFinancials objects per call
func ScanCmteFin(file io.ReadSeeker, start int64) ([]interface{}, int64, error) {
	offset := start
	// seek to starting byte offset
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, offset, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	queue := []interface{}{}

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

		// txDateFmt := "01/02/2006"
		// ced, err := time.Parse(txDateFmt, fieldMap[24])
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// create object to be stored in database
		fin := &donations.CmteFinancials{
			CmteID:          fieldMap[0],
			Type:            fieldMap[2],
			TotalReceipts:   float32(tr),
			TxsFromAff:      float32(tfa),
			IndvConts:       float32(ic),
			OtherConts:      float32(oc),
			CandCont:        float32(cc),
			TotalLoans:      float32(tl),
			TotalDisb:       float32(td),
			TxToAff:         float32(tta),
			IndvRefunds:     float32(ir),
			OtherRefunds:    float32(or),
			LoanRepay:       float32(lr),
			CashBOP:         float32(cbop),
			CashCOP:         float32(ccop),
			DebtsOwed:       float32(do),
			NonFedTxsRecvd:  float32(nft),
			ContToOtherCmte: float32(ctoc),
			IndExp:          float32(ie),
			PartyExp:        float32(pe),
			NonFedSharedExp: float32(nfe),
			// CovgEndDate:     ced,
		}

		// add donation to queue of items, stop at 25 items
		queue = append(queue, fin)
		if len(queue) == 10000 {
			break
		}
		fieldMap = make(mapOfFields)
	}
	return queue, offset, nil
}

// ScanContributions scans 500 lines of a contributions file
// and returns 500 Contribution objects per call
func ScanContributions(year string, file io.ReadSeeker, start int64) ([]*donations.Contribution, int64, error) {
	// seek to starting byte offset
	offset := int64(start)
	if _, err := file.Seek((offset), 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	icQueue := []*donations.Contribution{}

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
		/* txDateFmt := "01/02/2006"
		txDate, err := time.Parse(txDateFmt, fmtDateStr(fieldMap[13]))
		if err != nil {
			fmt.Println(err)
			fmt.Println("txID: ", fieldMap[16])
			return nil, start, fmt.Errorf("ParseContributions failed: %v", err)
		} */
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
			// TxDate:     txDate,
			TxAmt:    float32(txAmt),
			OtherID:  fieldMap[15],
			TxID:     fieldMap[16],
			FileNum:  fileNum,
			MemoCode: fieldMap[18],
			MemoText: fieldMap[19],
			SubID:    subID,
		}

		// add donation to queue of items, stop at 25 items
		icQueue = append(icQueue, donation)
		if len(icQueue) == 100000 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return icQueue, offset, nil
}

// ScanDisbursements scans 500 lines of a disbursements file
// and returns 500 Disbursement objects per call
func ScanDisbursements(year string, file io.ReadSeeker, start int64) ([]*donations.Disbursement, int64, error) {
	// seek to starting byte offset
	offset := int64(start)
	if _, err := file.Seek((offset), 0); err != nil {
		return nil, start, err
	}

	scanner := bufio.NewScanner(file)
	fieldMap := make(mapOfFields)
	dQueue := []*donations.Disbursement{}

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
		/* txDateFmt := "01/02/2006"
		txDate, err := time.Parse(txDateFmt, fieldMap[12])
		if err != nil {
			fmt.Println(err)
		} */
		txAmt, _ := strconv.ParseFloat(fieldMap[13], 32)
		fileNum, _ := strconv.Atoi(fieldMap[22])
		subID, _ := strconv.Atoi(fieldMap[21])

		// create object to be stored in database
		disb := &donations.Disbursement{
			CmteID: fieldMap[0],
			Name:   fieldMap[8],
			City:   fieldMap[9],
			State:  fieldMap[10],
			Zip:    fieldMap[11],
			// TxDate:       txDate,
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

		// add donation to queue of items, stop at 25 items
		dQueue = append(dQueue, disb)
		if len(dQueue) == 100000 {
			break
		}
		fieldMap = make(mapOfFields)
	}

	return dQueue, offset, nil
}

// Scan25CmteLink parses 25 rows of a candidate-committee links records file
// and creates a list of 25 items to be stored in the database
func Scan25CmteLink(file io.ReadSeeker, start int64) ([]*donations.CmteLink, int64, error) {
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
