// WORK PENDING
// NOT COMPLETE
// ALMOST COMPLETE

// PENDING:
// - export_json()
// - apply_func()
// - aggregate_column()

// COMMING SOON !!
// COMMING SOON !!
// COMMING SOON !!

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type CSV_Utils_Go struct {
	file_path string
	headers   []string
	rows      [][]string
	columns   int
}

// # constructor
func new_csv_utils_go(file_path string) (*CSV_Utils_Go, error) {
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

// _write_csv is a static method to write CSV data to a file.
func _write_csv(file_path string, data [][]string) {
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
		_write_csv(output_file_name, data_to_be_written)
		fmt.Println("New updated CSV file: ", output_file_name, " Created. Opeartion: ", op)
	} else {
		// update the same CSV File if no output_file provided
		_write_csv(c.file_path, data_to_be_written)
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
func (c *CSV_Utils_Go) get_column_index(column string, from_func_name string) int {
	err := c._validate(column)
	if err != nil {
		// requested column not found in the csv file, panic !!
		panic(
			fmt.Sprintf(
			"Error: '%s', Check if you have specified the right column name in function: `'%s'`", 
			err, 
			from_func_name))
	}
	for i, h := range c.headers {
		if h == column {
			return i
		}
	}
	return -1 // , errors.New("unexpected error in get_column_index")
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

// Helper function to determine if a string is numeric
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Helper function to round float to given precision
func roundFloat(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}

// Helper function to compute standard deviation (sample std dev)
func standardDeviation(vals []float64) float64 {
	if len(vals) <= 1 {
		return 0
	}
	mean := 0.0
	for _, v := range vals {
		mean += v
	}
	mean = mean / float64(len(vals))
	var sumSquares float64
	for _, v := range vals {
		sumSquares += (v - mean) * (v - mean)
	}
	return math.Sqrt(sumSquares / float64(len(vals)-1))
}

// Helper function to deep copy a 2D slice of strings
func deepCopy2D(src [][]string) [][]string {
	dest := make([][]string, len(src))
	for i, row := range src {
		dest[i] = make([]string, len(row))
		copy(dest[i], row)
	}
	return dest
}

// Helper function to get first N rows
func firstNRows(rows [][]string, n int) [][]string {
	if n > len(rows) {
		n = len(rows)
	}
	return rows[:n]
}

// Helper function to get last N rows
func lastNRows(rows [][]string, n int) [][]string {
	if n > len(rows) {
		n = len(rows)
	}
	return rows[len(rows)-n:]
}

// summerize provides a summary of the CSV data.
func (c *CSV_Utils_Go) summerize(preview_rows int) {
	/*
	   Provides a summary of the CSV data, including:
	   - Number of rows and columns
	   - Column-wise data types
	   - Column-wise missing values count
	   - Unique value count per column
	   - Most frequent value per column
	   - Min, Max, Mean, Std Dev for numerical columns
	   - First and Last few rows
	*/
	num_rows := len(c.rows)
	num_cols := len(c.headers)
	total_vals := 0

	data_types := make(map[string]map[string]bool)
	missing_values := make(map[string]int)
	unique_values := make(map[string]int)
	column_data := make(map[string][]string)
	numeric_stats := make(map[string]map[string]interface{})
	mode_values := make(map[string]string)

	// Initialize maps for headers
	for _, col := range c.headers {
		data_types[col] = make(map[string]bool)
		missing_values[col] = 0
		unique_values[col] = 0
		column_data[col] = []string{}
		numeric_stats[col] = map[string]interface{}{
			"min":     nil,
			"max":     nil,
			"mean":    nil,
			"std_dev": nil,
		}
	}

	// Process each row
	for _, row := range c.rows {
		for col_idx, value := range row {
			col_name := c.headers[col_idx]
			if value == "" || strings.TrimSpace(value) == "" {
				missing_values[col_name] += 1
			} else {
				total_vals++
				// All values are strings in Go, so type is always "string"
				data_types[col_name]["string"] = true
				unique_values[col_name]++
				column_data[col_name] = append(column_data[col_name], value)
			}
		}
	}

	// Convert data_types sets to lists
	data_types_list := make(map[string][]string)
	for col, typesMap := range data_types {
		typesList := []string{}
		for t := range typesMap {
			typesList = append(typesList, t)
		}
		data_types_list[col] = typesList
	}

	// For each column, compute numeric stats and mode
	for col, values := range column_data {
		is_numeric := true
		numeric_values := []float64{}
		for _, v := range values {
			if !isNumeric(v) {
				is_numeric = false
				break
			} else {
				f, _ := strconv.ParseFloat(v, 64)
				numeric_values = append(numeric_values, f)
			}
		}
		if is_numeric && len(numeric_values) > 0 {
			minVal := numeric_values[0]
			maxVal := numeric_values[0]
			sum := 0.0
			for _, num := range numeric_values {
				if num < minVal {
					minVal = num
				}
				if num > maxVal {
					maxVal = num
				}
				sum += num
			}
			meanVal := roundFloat(sum/float64(len(numeric_values)), 2)
			stdDev := 0.0
			if len(numeric_values) > 1 {
				stdDev = roundFloat(standardDeviation(numeric_values), 2)
			}
			numeric_stats[col] = map[string]interface{}{
				"min":     minVal,
				"max":     maxVal,
				"mean":    meanVal,
				"std_dev": stdDev,
			}
		}
		if len(values) > 0 {
			// Calculate most frequent value using frequency counter
			freq := make(map[string]int)
			for _, v := range values {
				freq[v]++
			}
			type kv struct {
				Key   string
				Value int
			}
			var freqList []kv
			for k, v := range freq {
				freqList = append(freqList, kv{k, v})
			}
			sort.Slice(freqList, func(i, j int) bool {
				return freqList[i].Value > freqList[j].Value
			})
			if len(freqList) > 0 {
				mode_values[col] = freqList[0].Key
			}
		}
	}

	summary := map[string]interface{}{
		"Total Rows:":                           num_rows,
		"Total Columns:":                        num_cols,
		"Total Values:":                         total_vals,
		"Column Data Types:":                    data_types_list,
		"Missing Values Per Column:":            missing_values,
		"Unique Values Per Column:":             unique_values,
		"Most Frequent Value (Mode) Per Column": mode_values,
		"Numeric Stats":                         numeric_stats,
		"First Few Rows":                        firstNRows(c.rows, preview_rows),
		"Last Few Rows":                         lastNRows(c.rows, preview_rows),
	}

	summaryBytes, err := json.MarshalIndent(summary, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling summary: %v\n", err)
		return
	}
	fmt.Println(string(summaryBytes))
}

// remove_duplicates removes duplicate rows based on the specified column.
func (c *CSV_Utils_Go) remove_duplicates(column string, output_file_name string) [][]string {
	/*
	   remove all duplictae values from the given `column`
	*/

	// fmt.Println("REMOVING DUPS !! DEBUG !")

	col_idx := c.get_column_index(column, "remove_duplicates")

	seen := make(map[string]bool)
	unique_rows := make([][]string, 0)

	for _, row := range c.rows {
		key := row[col_idx]
		if !seen[key] {
			seen[key] = true
			unique_rows = append(unique_rows, row)
		}
	}

	c.rows = unique_rows
	combined := append([][]string{c.headers}, c.rows...)

	c._update_csv(output_file_name, combined, "remove_duplicates()")

	return combined
}

// # replace_first_val replaces the first occurrence of `old_val` with `new_val` in the specified column.
func (c *CSV_Utils_Go) replace_first_val(
	column string,
	old_val string,
	new_val string,
	output_file_name string) {

	/*
	   Replaces all occurrences of `old_val` with `new_val` in the specified column
	   and updates the file or creates new file if given output file name.

	   :param column: Column name where the replacement should occur.
	   :param old_val: The value to be replaced.
	   :param new_val: The new value to replace with.
	   :output_file_name: The new file name in which updated data must be written
	*/

	col_idx := c.get_column_index(column, "replace_first_val")

	for i := range c.rows {
		if strings.EqualFold(c.rows[i][col_idx], old_val) {
			c.rows[i][col_idx] = new_val
			break
		}
	}

	combined := append([][]string{c.headers}, c.rows...)
	c._update_csv(output_file_name, combined, "`replace_first_val()`")
	// return nil
}

// replace_all_vals replaces all occurrences of old_val with new_val in the specified column.
func (c *CSV_Utils_Go) replace_all_vals(
	column string,
	old_val string,
	new_val string,
	output_file_name string) {

	/*SAME AS ABOVE FUNCTION*/

	col_idx := c.get_column_index(column, "replace_all_vals")

	for i := range c.rows {
		if strings.EqualFold(c.rows[i][col_idx], old_val) {
			c.rows[i][col_idx] = new_val
			break
		}
	}

	combined := append([][]string{c.headers}, c.rows...)
	c._update_csv(output_file_name, combined, "`replace_all_vals()`")
	// return nil
}

// is_palindrome is a static method to check if a given word is a palindrome.
func CSV_Utils_Py_is_palindrome(word string) bool {
	runes := []rune(word)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		if runes[i] != runes[n-1-i] {
			return false
		}
	}
	return true
}

// count_valid_palindromes counts palindrome words in the CSV rows.
func (c *CSV_Utils_Go) count_valid_palindromes() int {
	count := 0
	for _, row := range c.rows {
		for _, word := range row {
			word = strings.ToUpper(strings.TrimSpace(word))
			if word != "" && CSV_Utils_Py_is_palindrome(word) {
				count++
			}
		}
	}
	return count
}

// # filter_rows filters rows based on a condition on a specific column.
// It updates/filters the rows using the provided lambda callback.
func (c *CSV_Utils_Go) filter_rows(
	column string,
	condition func(string) bool,
	output_file_name string) ([][]string, error) {

	/*
	   updates / filters the specified rows with given lambda callback

	   :param: `column`: the column to be filtered

	   :param: `condition`: the lambda function to apply on column values

	   :param: `output_file_name`: the current CSV file is updated if not given else a new file is created
	*/

	col_idx := c.get_column_index(column, "filter_rows")

	filtered_rows := [][]string{c.headers}

	for _, row := range c.rows {
		if condition(row[col_idx]) {
			filtered_rows = append(filtered_rows, row)
		}
	}

	c._update_csv(output_file_name, filtered_rows, "")
	return filtered_rows, nil
}

// # sort_csv sorts the CSV data based on a specific column and writes the sorted data.
func (c *CSV_Utils_Go) sort_csv(
	column string,
	output_file_name string,
	ascending bool) ([][]string, error) {

	/*
	   sorts the specified `column` of the CSV File and updates or creates a new CSV File
	   :param: `output_file_name`: If given some name, the sorted CSV data would be written in that file
	*/

	col_idx := c.get_column_index(column, "sort_csv")

	// make a copy of rows for sorting to avoid modifying original order if needed
	sorted_rows := make([][]string, len(c.rows))
	copy(sorted_rows, c.rows)

	sort.Slice(sorted_rows, func(i, j int) bool {
		if ascending {
			return sorted_rows[i][col_idx] < sorted_rows[j][col_idx]
		}
		return sorted_rows[i][col_idx] > sorted_rows[j][col_idx]
	})

	c._update_csv(output_file_name, sorted_rows, "`sort_csv()`")

	combined := [][]string{c.headers}
	combined = append(combined, sorted_rows...)
	return combined, nil
}


// # aggregate_column performs aggregation on a numeric column using an operation.
func (c *CSV_Utils_Go) aggregate_column(column string, operation string) (float64, error) {

	/*
	Performs operations like: 'sum', 'min', 'max', 'std'

	:param: `column`: column name on which aggregation will be performed
	:param: `operation`: "sum" | "min" | "max' | "std"(standard deviation)
	:return: aggregated float value
	*/

	col_idx := c.get_column_index(column, "aggregate_column")

	values := []float64{}

	for _, row := range c.rows {

		// Check if the value is numeric by attempting conversion
		valStr := row[col_idx]

		// Remove one dot if present to mimic Python's replace('.', '', 1) for isdigit check
		valStrForCheck := strings.Replace(valStr, ".", "", 1)

		if _, err := strconv.ParseFloat(valStrForCheck, 64); err == nil {
			val, err := strconv.ParseFloat(valStr, 64)
			if err == nil {
				values = append(values, val)
			}
		}
	}

	if len(values) == 0 {
		return 0, fmt.Errorf("Column name: `'%s'` has no numeric values", column)
	}

	switch operation {
	case "sum":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum, nil
	case "avg":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values)), nil
	case "min":
		min := values[0]
		for _, v := range values {
			if v < min {
				min = v
			}
		}
		return min, nil
	case "max":
		max := values[0]
		for _, v := range values {
			if v > max {
				max = v
			}
		}
		return max, nil
	case "std":
		std_dev := standardDeviation(values)
		return std_dev, nil
	default:
		return 0, errors.New("invalid operation. Choose from 'sum', 'avg', 'min', 'max'")
	}
}


// output_processed_csv writes the current CSV data to the specified output path.
func (c *CSV_Utils_Go) output_processed_csv(output_path string) {

	if output_path == "" {
		output_path = c.file_path // Overwrite original file
	}

	file, err := os.Create(output_path)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write(c.headers)
	err = writer.WriteAll(c.rows)

	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Data saved to %s\n", output_path)
}


// apply_func applies a provided function to all values in the specified column.
func (c *CSV_Utils_Go) apply_func(
	column string, 
	funcToApply interface{}, 
	output_file_name string) ([][]string, error) {

	/*
	Apply a given function to all values in a specified column.

	:param column: Column name to apply function on.
	:param func: Function to apply.
	:param output_file_name: Optional filename to save modified data.
	:return: Modified CSV data as a list of lists.
	*/

	col_idx := c.get_column_index(column, "apply_func")

	modified_data := deepCopy2D(c.rows)

	for _, row := range modified_data {

		// Using recover to catch any panic during func execution
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Error processing row %v: %v\n", row, r)
				}
			}()

			// row[col_idx] = funcToApply(row[col_idx])

			switch f := funcToApply.(type) {

			case func(string) string:
				row[col_idx] = f(row[col_idx])

			case func(int) int:
				num, err := strconv.Atoi(row[col_idx])
				if err == nil {
					row[col_idx] = strconv.Itoa(f(num))
				}

			case func(float64) float64:
				num, err := strconv.ParseFloat(row[col_idx], 64)
				if err == nil {
					row[col_idx] = fmt.Sprintf("%f", f(num))
				} else {
					fmt.Printf("Skipping non-float value: %v\n", row[col_idx])
				}

			case func(float32) float32:
				num, err := strconv.ParseFloat(row[col_idx], 64)
				if err == nil {
					row[col_idx] = fmt.Sprintf("%f", f(float32(num)))
				} else {
					fmt.Printf("Skipping non-float value: %v\n", row[col_idx])
				}
				
			default:
				fmt.Printf("Unsupported function type for row %v\n", row)
			}

		}()
	}

	combined := append([][]string{c.headers}, modified_data...)
	c._update_csv(output_file_name, combined, "")

	return combined, nil
}


// export_json converts the CSV data to JSON format and optionally writes it to a file.
func (c *CSV_Utils_Go) export_json(json_file string) (string, error) {

	/*
    Convert CSV data to JSON format and optionally save it to a file.

	:param json_file: (Optional) Filename to save JSON data.
	:return: The JSON data as a string.
	*/

	data := make([]map[string]string, 0)

	for _, row := range c.rows {

		entry := make(map[string]string)

		for i, header := range c.headers {
			if i < len(row) {
				entry[header] = row[i]
			} else {
				entry[header] = ""
			}
		}

		data = append(data, entry)
	}

	jsonBytes, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		return "", err
	}

	json_data := string(jsonBytes)

	if json_file != "" {
		err := os.WriteFile(json_file, []byte(json_data), 0644)
		if err != nil {
			return "", err
		}
	}

	return json_data, nil
}