package csv

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"
)

// CSVData represents structured CSV output.
type CSVData struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

//export LoadCSV
func (c *CSVData) LoadCSV(filePath *C.char) *C.char {

	goFilePath := C.GoString(filePath)
	file, err := os.Open(goFilePath)

	if err != nil {
		return C.CString(`{"error": "CSV File Not Found"}`)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return C.CString(`{"error": "Failed to read headers"}`)
	}
	headers = trimSpaces(headers)

	// ! STORE THE CSV FILE HEADERS IN MEMORY FOR DOING REDUNDANT FILE-IO
	c.Headers = headers

	// Read rows
	var rows [][]string
	records, err := reader.ReadAll()
	if err != nil {
		return C.CString(`{"error": "Failed to read rows"}`)
	}

	for _, record := range records {
		
		row := trimSpaces(record)

		if allEmpty(row) {
			continue
		}

		// Fill missing values
		for len(row) < len(headers) {
			row = append(row, "N/A")
		}

		if len(row) > len(headers) {
			row = row[:len(headers)]
		}

		rows = append(rows, row)
	}

	// ! STORE THE CSV FILE HEADERS IN MEMORY FOR DOING REDUNDANT FILE-IO
	c.Rows = rows

	// Convert to JSON
	csvData := CSVData{Headers: headers, Rows: rows}
	jsonData, _ := json.Marshal(csvData)

	// Return JSON as C string
	return C.CString(string(jsonData))
}

//export filterRows
func (c *CSVData) filterRows(column *C.char, condition *C.char, writeToFile *C.int, outputFile *C.char) **C.char {
	goColumn := C.GoString(column)
	goCondition := C.GoString(condition)
	goOutputFile := C.GoString(outputFile)
	writeToFileFlag := int(writeToFile) != 0

	// Find column index
	colIdx := -1
	for i, h := range c.Headers {
		if strings.TrimSpace(h) == goColumn {
			colIdx = i
			break
		}
	}

	if colIdx == -1 {
		return nil // Column not found
	}

	// Filter rows
	filteredRows := [][]string{c.Headers}
	for _, row := range c.Rows {
		if colIdx < len(row) && row[colIdx] == goCondition {
			filteredRows = append(filteredRows, row)
		}
	}

	// Write to file if required
	if writeToFileFlag {
		outFile, err := os.Create(goOutputFile)
		if err == nil {
			defer outFile.Close()
			writer := csv.NewWriter(outFile)
			writer.WriteAll(filteredRows)
			writer.Flush()
		}
	}

	// Convert to C-compatible char array
	result := make([]*C.char, len(filteredRows))
	for i, row := range filteredRows {
		result[i] = C.CString(strings.Join(row, ","))
	}

	return &result[0]
}

func trimSpaces(s []string) []string {
	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}
	return s
}

func allEmpty(row []string) bool {
	for _, cell := range row {
		if cell != "" {
			return false
		}
	}
	return true
}

func main() {}