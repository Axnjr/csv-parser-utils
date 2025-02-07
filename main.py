import csv
import os
# for testing only, real functionlity is implemented from scratch !
import pandas as pd 
from typing import List, Callable

class CSV_Utils_Py:

    def __init__(self, file_path: str):
        self.file_path = file_path
        self.headers, self.rows = self._load_csv()


    def _load_csv(self) -> (List[str], List[List[str]]): # type: ignore

        if not os.path.exists(self.file_path):  
            print(f"File '{self.file_path}' not found. Creating an empty file.")
            with open(self.file_path, 'w', newline='') as csvfile:
                writer = csv.writer(csvfile)
                writer.writerow([])  # Empty header row
        
        with open(self.file_path, newline='') as csvfile:

            reader = csv.reader(csvfile)

            try:
                headers = next(reader)  # Read headers
                headers = [h.strip() for h in headers]  # Trim spaces
            except StopIteration:  # File is empty
                return [], []

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

        return headers, rows 

        # with open(self.file_path, newline='') as csvfile:
        #     reader = csv.reader(csvfile)
        #     headers = next(reader) 
        #     rows = list(reader)  
        # return headers, rows    


    def _validate(self, column: str):
        if column not in self.headers:
            raise ValueError(f"Column '{column}' not found in CSV.")


    def display_csv(self, num_rows: int = 3, add_index_col: bool = False):

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

        col_widths = [max(len(str(val)) for val in col) for col in zip(headers, *rows)]

        def print_row(row):
            formatted_values = [f"{value:<{col_widths[i]}}" for i, value in enumerate(row)]
            print("  ".join(formatted_values))

        print_row(headers)

        separators = "-" * (sum(col_widths) + (len(col_widths) - 1) * 2)
        print(separators) # for printing the "-----" line

        for row in rows:
            print_row(row)


    def filter_rows(
            self, 
            column: str, 
            condition: Callable[[str], bool],
            write_filtered_data_to_file: bool = False, 
            output_file_name: str = "filtered_csv_file.csv"
        ) -> List[List[str]]:

        self._validate(column)
        col_idx = self.headers.index(column)

        # ! less readable so i commented this one 
        # filtered_rows = [row for row in self.rows if condition(row[col_idx])]
        # return [self.headers] + filtered_rows

        filtered_rows = [self.headers]

        for row in self.rows:
            if condition(row[col_idx]):
                filtered_rows.append(row) 

        if write_filtered_data_to_file:
            CSV_Utils_Py.write_csv(output_file_name, filtered_rows)

        return filtered_rows


    def sort_csv(
            self, 
            column: str,  
            write_sorted_data_to_file: bool = False, 
            output_file_name: str = "sorted_csv_file.csv",
            ascending: bool = True
        ) -> List[List[str]]:
        
        self._validate(column)

        col_idx = self.headers.index(column)
        sorted_rows = sorted(self.rows, key=lambda row: row[col_idx], reverse=not ascending)

        if write_sorted_data_to_file: 
            print("\n SORTED DATA WRITTEN TO NEW FILE !! \n")
            CSV_Utils_Py.write_csv(output_file_name, sorted_rows)

        return [self.headers] + sorted_rows


    def aggregate_column(self, column: str, operation: str):

        self._validate(column)

        col_idx = self.headers.index(column)
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
        else:
            raise ValueError("Invalid operation. Choose from 'sum', 'avg', 'min', 'max'.")


    # def output_processed_csv(self, output_path: str = "processed_csv_file.csv"):
    #     if output_path is None:
    #         output_path = self.file_path  # Overwrite original file
    #     try:
    #         with open(output_path, 'w', newline='') as csvfile:
    #             writer = csv.writer(csvfile)
    #             writer.writerow(self.headers)
    #             writer.writerows(self.rows)
    #         print(f"Data saved to {output_path}")
    #     except Exception as e:
    #         print(f"Error writing file: {e}")

    @staticmethod
    def write_csv(file_path: str, data: List[List[str]]):
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


def run_test1():
    sample_data = [
        ['Name', 'Age', 'Score'], 
        ['Alice', '23', '85'], 
        ['Bob', '30', '90'], 
        ['Charlie', '25', '78']
    ]
    test_file = 'sample.csv'
    CSV_Utils_Py.write_csv(test_file, sample_data)

def run_test2():

    df = CSV_Utils_Py("python/ppl.csv")

    df.display_csv()

    # df.output_processed_csv("output_ppl.csv")

    filtered_data = df.filter_rows('Job Title', lambda x: x == "Games developer")
    print("Filtered Data:", filtered_data)
    
    sorted_data = df.sort_csv('Phone', ascending=False)
    print("Sorted Data:", sorted_data)
    
    total_score = df.aggregate_column('Index', 'sum')
    print("Total Score:", total_score)
    
    palindrome_count = df.count_valid_palindromes()
    print("Valid Palindromes Count:", palindrome_count)

    # os.remove(test_file)

    # data = pd.read_csv("python/ppl.csv")
    # print(data.to_string())

def run_test3():
    df = CSV_Utils_Py("python/fruits.csv")
    df.display_csv()

    df.output_processed_csv("output_fruits.csv")

    filtered_data = df.filter_rows('Price', lambda x: float(x) > 1.0, True)
    print("Filtered Data:", filtered_data)
    
    sorted_data = df.sort_csv('Date', ascending=False)
    print("Sorted Data:", sorted_data)
    
    total_score = df.aggregate_column('Price', 'sum')
    print("Total price:", total_score)
    
    palindrome_count = df.count_valid_palindromes()
    print("Valid Palindromes Count:", palindrome_count)

if __name__ == "__main__":
    run_test1()
    run_test2()
    run_test3()