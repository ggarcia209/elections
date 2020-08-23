package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Menu represents a menu/choice of selections
// The Options field repesents the options in a pre-defined order
// The OptionsMap field is used to reference each option by user-input numeric value.
type Menu struct {
	Name       string
	Options    []string
	OptionsMap map[int]string
}

// Ask4confirm asks for Y/N confirmation in a CLI application
func Ask4confirm() bool {
	var s string

	fmt.Printf("(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

// CreateMenu creates a new Menu object from the given name and list of options.
func CreateMenu(name string, options []string) Menu {
	m := Menu{Name: name, Options: options, OptionsMap: make(map[int]string)}

	for i, opt := range options {
		m.OptionsMap[i] = opt
	}
	return m
}

// Ask4MenuChoice asks user to choose an option from
// the given menu and returns the int value of the selection.
// Returned int values correspond to the k/v pairs in m.OptionsMap
func Ask4MenuChoice(m Menu) (int, error) {
	fmt.Println("Please choose an option: ")
	for i, opt := range m.Options {
		fmt.Printf("\t%d)  %s\n", i, opt)
	}
	var s string
	for {
		fmt.Printf("choose option: ")
		_, err := fmt.Scan(&s)
		if err != nil {
			fmt.Println(err)
			return -1, fmt.Errorf("Ask$MenuChoice failed %v: ", err)
		}

		// format & check input validity
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		if menuInputOk(s) {
			iv, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				return -1, fmt.Errorf("Ask$MenuChoice failed %v: ", err)
			}

			//
			if m.OptionsMap[iv] != "" {
				fmt.Println("chose", m.OptionsMap[iv])
				return iv, nil
			}
			fmt.Println("Invalid option - please try again")
		}
	}
}

// verify input is numerical value in string format
func menuInputOk(s string) bool {
	nums := map[string]bool{
		"0": true, "1": true, "2": true, "3": true, "4": true,
		"5": true, "6": true, "7": true, "8": true, "9": true,
	}
	if nums[s] == true { // check single digit
		return true
	}
	split := strings.Split(s, "") // check multiple digit
	for _, sp := range split {
		if nums[sp] != true {
			fmt.Println("invalid option - please try again")
			return false
		}
	}
	return true
}

// GetPathFromUser asks a user for a filepath and returns the path.
func GetPathFromUser() string {
	var s string
	var y string
	for {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Println("Working Directory: ", wd)
		fmt.Printf("enter filepath: ")
		_, err = fmt.Scan(&s)
		if err != nil {
			panic(err)
		}

		fmt.Println("confirm filepath is correct: ", s)
		fmt.Printf("(y/N): ")
		_, err = fmt.Scan(&y)
		if err != nil {
			panic(err)
		}

		// check if filepath is valid
		if _, err := os.Stat(s); os.IsNotExist(err) {
			fmt.Printf("filepath %s does not exist. Please try again.\n", s)
			continue
		}

		y = strings.TrimSpace(y)
		y = strings.ToLower(y)
		if y == "y" || y == "yes" {
			fmt.Println("new path: ", s)
			return s
		}
	}
}

// GetYear gets year gets the chosen year from the user
func GetYear() string {
	years := map[string]bool{
		"2020": true, "2018": true, "2016": true, "2014": true, "2012": true,
		"2010": true, "2008": true, "2006": true, "2004": true, "2002": true,
		"2000": true, "1998": true, "1996": true, "1994": true, "1992": true,
		"1990": true, "1988": true, "1986": true, "1984": true, "1982": true,
		"1980": true, "cancel": true,
	}

	var s string
	var y string
	for {
		fmt.Printf("Enter year: ")
		_, err := fmt.Scan(&s)
		if err != nil {
			panic(err)
		}

		s = strings.TrimSpace(s)
		s = strings.ToLower(s)

		if !years[s] {
			fmt.Println("Invalid year - please try again")
			continue
		}

		fmt.Println("Confirm year: ", s)
		fmt.Printf("(y/N): ")
		_, err = fmt.Scan(&y)
		if err != nil {
			panic(err)
		}

		y = strings.TrimSpace(y)
		y = strings.ToLower(y)
		if y == "y" || y == "yes" {
			return s
		}
	}
}

// GetQuery gets unformatted query input from user
func GetQuery() string {
	var s string
	fmt.Printf("Enter search query: ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Query: ", s)
	return s
}

// GetDynamoQuery gets a user-input sort, partition key pair
// for retreiving an object from a DynamoDB Table
func GetDynamoQuery() map[string]string {
	var prt string
	var srt string

	fmt.Printf("Enter partition key: ")
	_, err := fmt.Scan(&prt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Enter sort key: ")
	_, err = fmt.Scan(&srt)
	if err != nil {
		panic(err)
	}
	q := map[string]string{prt: srt}
	return q
}
