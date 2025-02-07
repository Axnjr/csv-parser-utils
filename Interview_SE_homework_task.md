Interview app - CSV File Processing Utility v1:
========================================================================
Build a CSV file utility tool that will help you perform various operations on CSV files. 
The code should be easy to deploy and run and must be hosted in a publicly available git repository. 
Your solution can be as robust or as minimal as you want but please keep in mind that this is your opportunity to show us what you can do. :)
You may use Python or Golang for the implementation of this utility. And, share your sample CSV file. 

Problem:
Write a program that implements a CSV utility with the following functionalities:
1. Read CSV and Display Data: Given a CSV file, read its contents and display the first few rows (e.g., 3 rows) on the console. The CSV file contains columns with headers, and the data is separated by commas.
2. Filter Rows Based on Criteria: Implement a function that filters rows based on a user-provided condition. For example, filter rows where the value in a specific column is greater than a given threshold or where the value of a string column contains a certain substring.
3. Sort Rows by Column: Allow the user to sort the rows in the CSV file based on a specific column. The sort can be done in ascending or descending order.
4. Aggregate Data (Sum, Average, etc.): Implement an aggregation function that performs simple statistics on a specific numeric column (e.g., sum, average, min, max). The user should be able to specify the column and the type of aggregation to perform.
5. Write to a New CSV File: After performing operations such as filtering or sorting, allow the user to save the resulting data into a new CSV file.
6. Check and display the number of palindromes in the CSV that consist solely of the letters A, D, V, B, and N.
7. Include tests to ensure that all the file functionalities mentioned above are working with your solution

Input:
- A CSV file with data where the first row contains column headers.
Example: Given the following CSV file fruits_sales_data.csv:
Date,Product,Quantity,Price
2025-01-01,Apple,10,1.2
2025-01-02,Banana,5,0.8
2025-01-03,Apple,15,1.3
2025-01-04,Orange,8,1.5
2025-01-05,Banana,12,0.7
- The user will interact with the utility, providing the necessary parameters (e.g., column name, condition, etc.) for the operations.
- The user can optionally specify the name of the new file to save changes.


Output:
- For the read and display data functionality, output the first few rows (with headers) to the console.
- For the filter rows functionality, output the filtered rows.
- For the sort rows functionality, output the rows sorted by the selected column.
- For the aggregate data functionality, output the aggregated result (sum, average, etc.).
- For the write to new CSV file functionality, create a new CSV file with the modified data.
- For display number of palindromes in CSV file, output of palindromes that satisfy above criteria

Constraints:
- The CSV file may have a large number of rows, so your solution should efficiently handle file reading and writing.
- Assume that the CSV file will have headers, and the columns will contain a mix of strings and numeric values.
- The user can specify a condition for filtering, sorting, or aggregation based on the column names provided in the header.
- You should implement the program to handle edge cases such as empty rows, missing values, and improper formats.
========================================================================

