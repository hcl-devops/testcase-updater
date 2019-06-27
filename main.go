package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	testcasecsvPath string
)

func init() {
	flag.StringVar(&testcasecsvPath, "csv", "./testcase.csv", "Path to the csv file to process. This is mandatory")

	flag.Usage = func() {
		fmt.Printf("Usage %s [args]: \n\n", os.Args)
		fmt.Println("Arguments: ")
		flag.PrintDefaults()
	}
	flag.Parse()
}

/**
  This program will accept a csv file as input and process the file.
  The test case id will be used to update the Jira test cases.
*/
func main() {

	if testcasecsvPath == "" {
		flag.Usage()
		return
	}
	testcases, err := ProcessCSV(testcasecsvPath)
	if err != nil {
		fmt.Printf("Error process csv file %s, Error: %s\n", testcasecsvPath, err)
		return
	}
	//fmt.Printf("cases: %v\n", testcases)
	if len(testcases) == 0 {
		fmt.Printf("No test cases to update hence exiting\n")
		return
	}
	err = UpdateJiraIssues(testcases)
	if err != nil {
		fmt.Printf("Issues updating the Jira. Exiting\n")
	}
}
