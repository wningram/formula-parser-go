package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type CsvFile struct {
	FileName string
	Headers  []string
	DataRows [][]string
}

func (f CsvFile) GetColumn(colName string) []string {
	var colNdx int = -1
	var result []string

	for ndx, col := range f.Headers {
		if colName == strings.Trim(col, " ") {
			colNdx = ndx
		}
	}

	// Check if result was foudn
	if colNdx == -1 {
		return []string{}
	}

	// Get column data
	for _, row := range f.DataRows {
		result = append(result, row[colNdx])
	}

	return result
}

func (f CsvFile) String() string {
	result := ""
	// Print headers
	result += "HEADERS:\n------------\n"
	for _, v := range f.Headers {
		result += v + "\n"
	}
	// Print data
	result += "\nDATA:\n-----------\n"
	for _, v := range f.DataRows {
		for _, v2 := range v {
			result += v2 + "\n"
		}
		result += "\n"
	}

	return result
}

func LoadCsv(fileName string) (CsvFile, error) {
	var (
		headers []string
		data    [][]string
	)
	csvFile, err := os.Open(fileName)
	if err != nil {
		return CsvFile{}, err
	}

	cr := csv.NewReader(csvFile)
	rowNdx := 0
	// Read the CSV file
	for {
		row, err := cr.Read()
		rowNdx++
		if err == io.EOF {
			break
		} else if err != nil {
			return CsvFile{}, fmt.Errorf("Could not read CSV file: %s", err)
		}

		if rowNdx == 1 {
			headers = row
		} else {
			data = append(data, row)
		}
	}

	return CsvFile{
		FileName: fileName,
		Headers:  headers,
		DataRows: data,
	}, nil
}
