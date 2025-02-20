/// PENDING !!
/// BINDINGS COMMING SOON ⚒️👷‍♂️🫷🫸✌️✨

use std::error::Error;
use std::fs::File; 
use std::io::{self, BufReader, BufWriter, Write};
use std::path::Path;
use csv::Reader;
use pyo3::prelude::*;

pub struct CsvPyBinder {
    pub file_path: &'static str,
    pub headers: Option<Vec<String>>,
    pub rows: Option<Vec<Vec<String>>>,
    pub coloumns: usize,
}

impl CsvPyBinder {

    pub fn new(&self, file_path: &'static str) -> Self {
        match self.load_csv() {
            Ok(data) => {
                return CsvPyBinder {
                    file_path,
                    headers: Some(data.0),
                    rows: Some(data.1),
                    coloumns: data.2
                }
            },
            Err(err) => {
                println!("Unable to load CSV File. {}", err);
                panic!("Unable to load CSV File.")
            }
        }
    }

    
    fn load_csv(&self) -> Result<(Vec<String>, Vec<Vec<String>>, usize), Box<dyn Error>> {

        let path = Path::new(&self.file_path);

        if !path.exists() {
            panic!("CSV File Not Found !!");
        }
        
        let file = File::open(&self.file_path)?;
        let mut reader = Reader::from_reader(BufReader::new(file));
        let mut rows: Vec<Vec<String>> = Vec::new();
        let mut headers = vec![];

        if let Some(result) = reader.headers().ok() {
            // trim spaces from the first row of the file i.e csv
            headers = result.iter().map(|h| h.trim().to_string()).collect();
        } 
        else {
            // Return empty if the file has no headers
            return Ok((headers, rows, 0)); 
        }

        for result in reader.records() {

            let record = result?;

            // trim spaces
            let mut row: Vec<String> = record
                .iter()
                .map(|cell| cell.trim().to_string())
                .collect()
            ;

            // Ignore completely empty rows
            if row.iter().all(|cell| cell.is_empty()) {
                continue;
            }

            // Handle missing or extra columns
            if row.len() < headers.len() {
                row.resize(headers.len(), "N/A".to_string());
            } 

            else if row.len() > headers.len() {
                row.truncate(headers.len());
            }

            rows.push(row);
        }

        let columns = &headers.len();
        Ok((headers, rows, *columns))
    }


    fn validate(&self, column: &str) -> Result<(), Box<dyn Error>> {
        if let Some(headers) = &self.headers {
            if headers.contains(&column.to_string()) {
                return Ok(());
            } else {
                return Err(Box::from(format!("Column: {} not found in CSV.", column)));
            }
        }
        Err(Box::from("Headers not loaded."))
    }


    pub fn display_csv(&self, num_rows: usize, add_index_col: bool) {

        if let Some(headers) = &self.headers {
            if headers.is_empty() {
                println!("CSV File is Empty, no headers found !");
                return;
            }
        }

        println!("Here are the first {} rows of the CSV file:\n", num_rows);

        let mut headers = self.headers.as_ref().unwrap().clone();
        let mut rows: Vec<Vec<String>> = self.rows.as_ref().unwrap().iter().take(num_rows).cloned().collect::<Vec<_>>();

        // Add index column if requested
        if add_index_col 
        {
            headers.insert(0, "INDEX".to_string());

            for (i, row) in rows.iter_mut().enumerate() {
                row.insert(0, i.to_string());
            }
        }

        // Helper function to print a row
        let print_row = |row: &[String]| {
            let formatted_values: Vec<String> = row
                .iter()
                .enumerate()
                .map(|(i, value)| format!("{:<width$}", value, width = self.coloumns))
                .collect();
            println!("{}", formatted_values.join("  "));
        };

        print_row(&headers);

        // Print rows
        for row in &rows {
            print_row(row);
        }
    }


    pub fn filter_rows<F: Fn(&str) -> bool>(
        &self, 
        column: &str, 
        condition: F, 
        write_filtered_data_to_file: bool, 
        output_file_name: &str
    ) -> Vec<Vec<String>> {

        let _ = self.validate(column);
        let col_idx = self.headers.as_ref().unwrap().iter().position(|col_name| col_name == column);

        let mut filtered_rows: Vec<Vec<String>> = Vec::new();

        if let Some(idx) = col_idx {
            for row in self.rows.as_ref().unwrap_or(&vec![]) {
                if condition(&row[idx]) {
                    filtered_rows.push(row.clone());
                }
            }
        }

        if write_filtered_data_to_file {
            let _ = Self::write_csv(output_file_name, &filtered_rows);
        }

        filtered_rows
    }


    pub fn sort_csv(
        &self, 
        column: &str, 
        write_filtered_data_to_file: bool, 
        output_file_name: &str,
        ascending: bool
    ) -> Vec<Vec<String>> {

        let _ = self.validate(column);
        let col_idx = self.headers
            .as_ref()
            .unwrap()
            .iter()
            .position(|col_name| col_name == column)
            .expect("Specified Column not found the CSV file !")
        ;

        let mut sorted_rows = self.rows.as_ref().unwrap().to_owned();

        sorted_rows.sort_by(|a, b| {
            if ascending {
                a[col_idx].cmp(&b[col_idx])
            }
            else {
                b[col_idx].cmp(&a[col_idx])
            }
        });

        if write_filtered_data_to_file {
            let _ = Self::write_csv(output_file_name, &sorted_rows);
        }

        sorted_rows
    }


    /// operations include: sum, min, max, avg. <br>
    /// column: the column whose values you want to aggregate
    pub fn aggregate_column(self, column: &str, operation: &str) -> Result<f64, String> {

        let _ = self.validate(column);
        
        let col_idx = self.headers
            .as_ref()
            .unwrap()
            .iter()
            .position(|col_name| col_name == column)
            .expect("Specified Column not found the CSV file !")
        ;

        let values: Vec<f64> = self.rows
            .as_ref()
            .unwrap()
            .iter()
            .filter_map(|row| row[col_idx].parse::<f64>().ok())
            .collect()
        ;

        if values.is_empty() {
            return Err(format!("Column '{}' has no numeric values.", column));
        }

        match operation {
            "sum" => Ok(values.iter().sum()),
            "avg" => Ok(values.iter().sum::<f64>() / values.len() as f64),
            "min" => Ok(*values.iter().min_by(|a, b| a.partial_cmp(b).unwrap()).unwrap()),
            "max" => Ok(*values.iter().max_by(|a, b| a.partial_cmp(b).unwrap()).unwrap()),
            _ => Err("Invalid operation. Choose from 'sum', 'avg', 'min', 'max'.".to_string()),
        }

    }

    fn write_csv(file_path: &str, data: &[Vec<String>]) -> io::Result<()> {
        let file = File::create(file_path)?;
        let mut writer = BufWriter::new(file);
    
        for row in data {
            writeln!(writer, "{}", row.join(","))?;
        }
    
        println!("Data saved to {}", file_path);
        Ok(())
    }

    fn is_palindrome(word: &str) -> bool {
        let clean_word = word.trim().to_uppercase();
        clean_word.chars().eq(clean_word.chars().rev())
    }

    pub fn count_valid_palindromes(&self) -> usize {
        self.rows
            .as_ref()
            .unwrap()
            .iter()
            .flat_map(|row| row.iter())
            .filter(|word| !word.trim().is_empty() && Self::is_palindrome(word))
            .count()
    }


}

/// Formats the sum of two numbers as string.
#[pyfunction]
fn sum_as_string(a: usize, b: usize) -> PyResult<String> {
    Ok((a + b).to_string())
}

/// A Python module implemented in Rust.
#[pymodule]
fn redhat_test(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(sum_as_string, m)?)?;
    Ok(())
}
