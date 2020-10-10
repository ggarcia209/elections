package qwikscript

import (
	"bufio"
	"fmt"
	"os"

	"github.com/elections/source/persist"

	"github.com/elections/source/ui"
	"github.com/elections/source/util"
)

func PrintRecordAtOffset() {
	fmt.Println("enter filepath")
	path := ui.GetPathFromUser()
	start := 0

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("PrintRecordAtOffset failed: ", err)
		os.Exit(1)
	}
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("PrintRecordAtOffset failed: ", err)
		os.Exit(1)
	}

	offset := int64(start)
	if _, err := file.Seek(offset, 0); err != nil {
		fmt.Println("PrintRecordAtOffset failed: ", err)
		os.Exit(1)
	}

	fmt.Println(fi.Name())
	fmt.Println("offset: ", offset)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		break
	}
}

func ConvertInt64() {
	// bill := int64(6607314240)
	bill := int64(4288771025)
	fmt.Println("bill start: ", bill)
	bytes := util.Itob(bill)
	fmt.Println("bytes: ", bytes)
	bill = util.Btoi(bytes)
	fmt.Println("bill end: ", bill)

	fmt.Println("Getting")
	err := persist.LogOffset("2018", "test", bill)
	if err != nil {
		fmt.Println("convertInt64 failed: ", err)
		os.Exit(1)
	}
	fmt.Println("offset logged: ", bill)
	offset, err := persist.GetOffset("2018", "test")
	if err != nil {
		fmt.Println("convertInt64 failed: ", err)
		os.Exit(1)
	}
	fmt.Println("GotOffset from disk: ", offset)
}
