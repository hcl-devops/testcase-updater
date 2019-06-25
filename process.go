package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// TestCase structure
type TestCase struct {
	TestID     string
	TestStepID string
	Status     string
}

/*
 process the csv and return csv structure
*/
func ProcessCSV(csvpath string) ([]TestCase, error) {
	var testcases []TestCase
	if fileMode, err := os.Stat(csvpath); err == nil && !fileMode.IsDir() {
		fmt.Printf("File is not a directory.")
		csvFile, err := os.Open(csvpath)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		if err != nil {
			return nil, err
		}
		i := 0
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return testcases, err
			}
			if i == 0 {
				i++
				continue
			}
			testcases = append(testcases, TestCase{
				TestID:     line[0],
				TestStepID: line[1],
				Status:     line[2],
			})
			i++
		}
	} else {
		fmt.Printf("Err reading file %s\n", err)
	}
	return testcases, nil
}
