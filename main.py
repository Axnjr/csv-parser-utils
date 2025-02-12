from collections import Counter
import csv
import json
import os
# for testing only, real functionlity is implemented from scratch !
# import pandas as pd 
import statistics
from typing import List, Callable
from copy import deepcopy

class CSV_Utils_Py:

    def __init__(self, file_path: str):
        
        self.file_path = file_path
        data = self._load_csv()

        self.headers: List[str]       = data[0] 
        self.rows   : List[List[str]] = data[1] 
        self.columns: int             = data[2]


    def _load_csv(self) -> (List[str], List[List[str]], int): # type: ignore

        if not os.path.exists(self.file_path):  

            print(f"File '{self.file_path}' not found. Creating an empty file.")
            raise Exception("CSV File not found !!")
            # with open(self.file_path, 'w', newline='') as csvfile:
            #     writer = csv.writer(csvfile)
            #     writer.writerow([])  # Empty header row
        
        with open(self.file_path, newline='') as csvfile:

            reader = csv.reader(csvfile)

            try:
                headers = next(reader) 
                headers = [h.strip() for h in headers]  # Trim spaces
            except StopIteration:  # File is empty
                return [], [], 0

            # storing reader iterator to list for re-using them later when the file is closed !
            rows = []
            for row in reader:
                # ignore completely empty rows
                if not any(cell.strip() for cell in row):
                    continue
                # missing values -> "N/A" if row is shorter than headers
                if len(row) < len(headers):
                    row.extend(["N/A"] * (len(headers) - len(row)))
                # extra columns: trim excess values
                elif len(row) > len(headers):
                    row = row[:len(headers)]

                rows.append(row)

        return headers, rows, len(headers) 

        # with open(self.file_path, newline='') as csvfile:
        #     reader = csv.reader(csvfile)
        #     headers = next(reader) 
        #     rows = list(reader)  
        # return headers, rows    


    def _update_csv(self, output_file_name: str, data_to_be_written: any, op: str):

        if output_file_name:
            CSV_Utils_Py.write_csv(output_file_name, data_to_be_written)
            print("New updated CSV file: ", output_file_name, " Created. Opeartion: ", op)

        else:
            # update the same CSV File if no output_file provided
            CSV_Utils_Py.write_csv(self.file_path, data_to_be_written)
            print("CSV File updated. Operation: ", op)


    def _validate(self, column: str):
        if column not in self.headers:
            raise ValueError(f"Column '{column}' not found in CSV.")
        

    def get_column_index(self, column: str) -> int:
        self._validate(column)
        return self.headers.index(column)


    def display_csv(self, num_rows: int = 3, add_index_col: bool = False):

        """
        prints the first `num_rows` of the CSV File
        :param add_index_col: A boolean to add an extra index column for better data visualization
        returns void
        """

        if not self.headers:
            print("CSV file is empty!")
            return

        print(f"Here are first {num_rows} rows of the csv file: \n")

        headers = self.headers
        rows = [row for _, row in zip(range(num_rows), self.rows)]

        # for easily identifying rows of the csv file, like pandas 
        if add_index_col:
            headers = ["INDEX"] + headers  
            rows = [[str(i)] + row for i, row in enumerate(rows)]  

        # headers: list[str], rows: list[str]
        # col_widths = [max(len(str(val)) for val in col) for col in zip(headers, *rows)]

        def print_row(row):
            formatted_values = [f"{value:<{self.columns}}" for i, value in enumerate(row)]
            # formatted_values = [f"{value:<{col_widths}}" for i, value in enumerate(row)]
            print("  ".join(formatted_values))

        print_row(headers)

        # ! SEPARATORS ARE NOT EFFICIENT FOR LARGE FILES, AS WE HAVE TO ZIP AND THE DO NESTED ITERATIONS !!
        # separators = "-" * (sum(col_widths) + (len(col_widths) - 1) * 2)
        # print(separators) # for printing the "-----" line

        for row in rows:
            print_row(row)


    def filter_rows(
            self, 
            column: str, 
            condition: Callable[[str], bool],
            output_file_name: str = None
        ) -> List[List[str]]:

        """
        updates / filters the specified rows with given lambda callback

        :param: `column`: the column to be filtered

        :param: `condition`: the lambda function to apply on column values

        :param: `output_file_name`: the current CSV file is updated if not given else a new file is created
        """

        col_idx = self.get_column_index(column)

        # ! less readable so i commented this one 
        # filtered_rows = [row for row in self.rows if condition(row[col_idx])]
        # return [self.headers] + filtered_rows

        filtered_rows = [self.headers]

        for row in self.rows:
            if condition(row[col_idx]):
                filtered_rows.append(row) 

        self._update_csv(output_file_name, filtered_rows)
        return filtered_rows


    def sort_csv(
            self, 
            column: str,  
            output_file_name: str,
            ascending: bool = True
        ) -> List[List[str]]:

        """
        sorts the specified `column` of the CSV File and updates or creates a new CSV File

        :param: `output_file_name`: If given some name, the sorted CSV data would be written in that file
        """
        
        col_idx = self.get_column_index(column)

        sorted_rows = sorted(self.rows, key=lambda row: row[col_idx], reverse=not ascending)

        self._update_csv(output_file_name, sorted_rows, "`sort_csv()`")

        return [self.headers] + sorted_rows


    def aggregate_column(self, column: str, operation: str):

        """
        Performs operations like: 'sum', 'min', 'max', 'std'

        :param: `column`: column name on which aggregation will be performed
        :param: `operation`: "sum" | "min" | "max' | "std"(standard deviation)
        :return: aggregated float value
        """

        col_idx = self.get_column_index(column)

        values = [float(row[col_idx]) for row in self.rows if row[col_idx].replace('.', '', 1).isdigit()]

        if not values:
            raise ValueError(f"Column '{column}' has no numeric values.")

        if operation == 'sum':
            return sum(values)
        elif operation == 'avg':
            return sum(values) / len(values)
        elif operation == 'min':
            return min(values)
        elif operation == 'max':
            return max(values)
        elif operation == 'std':
            return statistics.stdev(values)
        else:
            raise ValueError("Invalid operation. Choose from 'sum', 'avg', 'min', 'max'.")


    def output_processed_csv(self, output_path: str = "processed_csv_file.csv"):
        if output_path is None:
            output_path = self.file_path  # Overwrite original file
        try:
            with open(output_path, 'w', newline='') as csvfile:
                writer = csv.writer(csvfile)
                writer.writerow(self.headers)
                writer.writerows(self.rows)
            print(f"Data saved to {output_path}")
        except Exception as e:
            print(f"Error writing file: {e}")


    @staticmethod
    def write_csv(file_path: str, data: List[List[str]]):
        # create file if not exists
        with open(file_path, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerows(data)
        print(f"Data saved to {file_path}")


    @staticmethod
    def is_palindrome(word: str) -> bool:
        return word == word[::-1]


    def count_valid_palindromes(self) -> int:
        count = 0
        for row in self.rows:
            for word in row:
                word = word.strip().upper()
                if word and self.is_palindrome(word): 
                    count += 1
        return count
    

    def apply_func(self, column: str, func, output_file_name: str = None) -> List[List[str]]:

        """
        Apply a given function to all values in a specified column.

        :param column: Column name to apply function on.
        :param func: Function to apply.
        :param output_file_name: Optional filename to save modified data.
        :return: Modified CSV data as a list of lists.
        """

        col_idx = self.get_column_index(column)

        # Create a deep copy to avoid modifying original data
        modified_data = deepcopy(self.rows)

        for row in modified_data:
            try:
                row[col_idx] = func(row[col_idx])
            except Exception as e:
                print(f"Error processing row {row}: {e}")

        self._update_csv(output_file_name, [self.headers] + modified_data, "apply_func()")

        return [self.headers] + modified_data


    def export_json(self, json_file: str):

        """
        Convert CSV data to JSON format and optionally save it to a file.

        :param json_file: (Optional) Filename to save JSON data.
        :return: The JSON data as a string.
        """

        data = [dict(zip(self.headers, row)) for row in self.rows]
        json_data = json.dumps(data, indent=4)

        if json_file:
            with open(json_file, mode="w", encoding="utf-8") as file:
                file.write(json_data)

        return json_data


    def replace_all_vals(self, column: str, old_val: str, new_val: str, output_file_name: str = None):

        """
        Replaces all occurrences of `old_val` with `new_val` in the specified column 
        and updates the file or creates new file if given output file name.

        :param column: Column name where the replacement should occur.
        :param old_val: The value to be replaced.
        :param new_val: The new value to replace with.
        :output_file_name: The new file name in which updated data must be written
        """

        col_idx = self.get_column_index(column)

        for row in self.rows:
            if(row[col_idx].lower() == old_val.lower()):
                row[col_idx] = new_val

        self._update_csv(output_file_name, [self.headers] + self.rows, "`replace_all_vals()`")


    def replace_first_val(self, column: str, old_val: str, new_val: str, output_file_name: str = None):

        col_idx = self.get_column_index(column)

        for row in self.rows:
            if(row[col_idx].lower() == old_val.lower()):
                row[col_idx] = new_val
                break

        self._update_csv(output_file_name, [self.headers] + self.rows, "`replace_first_val()`")


    def remove_duplicates(self, column: str, output_file_name: str = None):

        """
        remove all duplictae values from the given `column`
        """

        col_idx = self.get_column_index(column)
        seen = set()
        unique_rows = []

        for row in self.rows:
            key = row[col_idx]
            if key not in seen:
                seen.add(key)
                unique_rows.append(row)

        self.rows = unique_rows
        self._update_csv(output_file_name, [self.headers] + self.rows, "remove_duplicates()")


    def summerize(self, preview_rows: int = 3):

        """
        Provides a summary of the CSV data, including:
        - Number of rows and columns
        - Column-wise data types
        - Column-wise missing values count
        - Unique value count per column
        - Most frequent value per column
        - Min, Max, Mean, Std Dev for numerical columns
        - First and Last few rows
        """

        num_rows = len(self.rows)
        num_cols = len(self.headers)
        total_vals = 0

        data_types     = {col: set() for col in self.headers}
        missing_values = {col: 0 for col in self.headers}
        unique_values  = {col: 0 for col in self.headers}
        column_data    = {col: [] for col in self.headers}
        numeric_stats = {
            col: {
                "min": None, 
                "max": None, 
                "mean": None, 
                "std_dev": None
            } for col in self.headers
        }

        mode_values = {}

        for row in self.rows:

            for col_idx, value in enumerate(row):

                col_name = self.headers[col_idx]

                if value == "" or value == None or value.isspace():
                    missing_values[col_name] += 1

                else:
                    total_vals += 1
                    data_types[col_name].add(type(value).__name__)
                    unique_values[col_name] += 1
                    column_data[col_name].append(value)

        # ! sets cannot be json serialized so convert sets to lists
        data_types = {col: list(types) for col, types in data_types.items()}

        for col, values in column_data.items():

            is_numeric = all(v.replace('.', '', 1).isdigit() for v in values if v)

            if is_numeric:

                numeric_values = [float(v) for v in values]

                numeric_stats[col] = {
                    "min": min(numeric_values),
                    "max": max(numeric_values),
                    "mean": round(sum(numeric_values) / len(numeric_values), 2),
                    "std_dev": round(statistics.stdev(numeric_values), 2) if len(numeric_values) > 1 else 0
                }

            if values:
                most_common_value = Counter(values).most_common(1)  # Returns the most frequent value
                mode_values[col] = most_common_value[0][0]  # Extract the actual value

        summary = {
            "Total Rows:": num_rows,
            "Total Columns:": num_cols,
            "Total Values:": total_vals,
            "Column Data Types:": data_types,
            "Missing Values Per Column:": missing_values,
            "Unique Values Per Column:": unique_values,
            "Most Frequent Value (Mode) Per Column": mode_values,
            "Numeric Stats": numeric_stats,
            "First Few Rows": self.rows[:preview_rows],
            "Last Few Rows": self.rows[-preview_rows:],
        }

        print(json.dumps(summary, indent=4))

