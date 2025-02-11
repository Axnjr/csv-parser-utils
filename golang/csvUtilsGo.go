// WORK PENDING
// NOT COMPLETE
// IMPLEMENTED A FEW PRIVATE METHODS, THAT'S IT:

// COMMING SOON !!
// COMMING SOON !!
// COMMING SOON !!

package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"strconv"
	// "math"
	// "sort"
	// "encoding/json"
)

type CSV_Utils_Go struct {
	file_path string
	headers   []string
	rows      [][]string
	columns   int
}


// constructor
func NewCSV_Utils_Py(file_path string) (*CSV_Utils_Go, error) {
	csvUtil := &CSV_Utils_Go{file_path: file_path}
	headers, rows, columns, err := csvUtil._load_csv()
	if err != nil {
		return nil, err
	}
	csvUtil.headers = headers
	csvUtil.rows = rows
	csvUtil.columns = columns
	return csvUtil, nil
}


func (c *CSV_Utils_Go) _load_csv() ([]string, [][]string, int, error) {
	// Check if file exists
	if _, err := os.Stat(c.file_path); os.IsNotExist(err) {
		fmt.Printf("File '%s' not found. Creating an empty file.\n", c.file_path)
		return nil, nil, 0, errors.New("CSV File not found")
	}

	// Open CSV file
	file, err := os.Open(c.file_path)
	if err != nil {
		return nil, nil, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// Try reading headers
	var headers []string
	firstLine, err := reader.Read()
	if err == io.EOF {
		return []string{}, [][]string{}, 0, nil
	} else if err != nil {
		return nil, nil, 0, err
	}

	// Trim spaces in headers
	for _, h := range firstLine {
		headers = append(headers, strings.TrimSpace(h))
	}

	rows := make([][]string, 0)

	// Read all rows from csv reader
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, 0, err
		}
		// ignore completely empty rows
		empty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}

		// missing values -> "N/A" if row is shorter than headers
		if len(row) < len(headers) {
			for i := 0; i < (len(headers) - len(row)); i++ {
				row = append(row, "N/A")
			}
		} else if len(row) > len(headers) {
			// extra columns: trim excess values
			row = row[:len(headers)]
		}

		rows = append(rows, row)
	}

	return headers, rows, len(headers), nil
}


// CSV_Utils_Py_write_csv is a static method to write CSV data to a file.
func CSV_Utils_Py_write_csv(file_path string, data [][]string) {
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	fmt.Printf("Data saved to %s\n", file_path)
}


func (c *CSV_Utils_Go) _update_csv(output_file_name string, data_to_be_written [][]string, op string) {
	if output_file_name != "" {
		CSV_Utils_Py_write_csv(output_file_name, data_to_be_written)
		fmt.Println("New updated CSV file: ", output_file_name, " Created. Opeartion: ", op)
	} else {
		// update the same CSV File if no output_file provided
		CSV_Utils_Py_write_csv(c.file_path, data_to_be_written)
		fmt.Println("CSV File updated. Operation: ", op)
	}
}


// _validate checks if the given column exists in the headers, else returns an error
func (c *CSV_Utils_Go) _validate(column string) error {
	for _, h := range c.headers {
		if h == column {
			return nil
		}
	}
	return fmt.Errorf("%s", fmt.Sprintf("Column '%s' not found in CSV", column))
}


// get_column_index returns the index of the specified column
func (c *CSV_Utils_Go) get_column_index(column string) (int, error) {
	err := c._validate(column)
	if err != nil {
		return -1, err
	}
	for i, h := range c.headers {
		if h == column {
			return i, nil
		}
	}
	return -1, errors.New("unexpected error in get_column_index")
}


// display_csv prints the first num_rows of the CSV data.
// If add_index_col is true, an extra index column is added.
func (c *CSV_Utils_Go) display_csv(num_rows int, add_index_col bool) {
	/*
	   prints the first `num_rows` of the CSV File
	   :param add_index_col: A boolean to add an extra index column for better data visualization
	   returns void
	*/
	if len(c.headers) == 0 {
		fmt.Println("CSV file is empty!")
		return
	}

	fmt.Printf("Here are first %d rows of the csv file: \n\n", num_rows)

	headers := c.headers
	rows := make([][]string, 0)

	for i := 0; i < num_rows && i < len(c.rows); i++ {
		rows = append(rows, c.rows[i])
	}

	// for easily identifying rows of the csv file, like pandas in python
	if add_index_col {
		headers = append([]string{"INDEX"}, headers...)
		newRows := make([][]string, 0)
		for i, row := range rows {
			indexStr := strconv.Itoa(i)
			newRow := append([]string{indexStr}, row...)
			newRows = append(newRows, newRow)
		}
		rows = newRows
	}

	// print_row helper function
	print_row := func(row []string) {
	
		formatted_values := make([]string, 0)

		for range row {
			// Use fixed width equal to number of columns (as per Python code)
			// This is a simplification of the Python formatting.
			// In Python, each column is formatted with width self.columns.
			// Here we use that same self.columns.
			// Note: This might not exactly mimic Python behavior if columns vary.
			// Using fmt.Sprintf with left alignment.
			formatted_values = append(formatted_values, "")
		}

		for i, value := range row {
			formatted_values[i] = fmt.Sprintf("%-*s", c.columns, value)
		}

		fmt.Println(strings.Join(formatted_values, "  "))
	}

	print_row(headers)
	// Print each row
	for _, row := range rows {
		print_row(row)
	}
}


func main() {
	// This main function is provided for testing purposes.
	// It can be modified as needed to test the functionality of CSV_Utils_Py.
	// Example usage:
	
	csvUtil, err := NewCSV_Utils_Py("C:/Users/yaksh/redhat-test/ppl.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	csvUtil.display_csv(3, true)
	// !todo => csvUtil.summerize(3)

}