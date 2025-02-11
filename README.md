# `CSV_Utils_Py` 
# Homework task for SWE position at Redhat.

## Public Methods:
- `display_csv()`
- `filter_rows()`
- `sort_csv()`
- `aggregate_column()`
- `count_valid_palindromes()`
- `write_csv()`
- `remove_duplicates()`
- `replace_all_vals()`
- `reaplace_first_val()`
- `read_csv()`
- `summerize()`
- `export_json()`
- `apply_func()`
- `output_processed_csv()`

## Private Methods:
- `_update_csv()`
- `_validate()`
- `_load_csv()`

## Usage / test cases
```py
def run_test1():

    sample_data = [
        ['Name', 'Age', 'Score'], 
        ['Alice', '23', '85'], 
        ['Bob', '30', '90'], 
        ['Charlie', '25', '78']
    ]
    test_file = 'sample.csv'
    CSV_Utils_Py.write_csv(test_file, sample_data)

# ------------------------------------------------------------------------------------------------------------------ #

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

# ------------------------------------------------------------------------------------------------------------------ #

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

# ------------------------------------------------------------------------------------------------------------------ #

def run_test4():
    
    df = CSV_Utils_Py("./ppl.csv")

    # replace all "First Name": Shelby to Radha
    df.replace_all_vals("First Name", "Shelby", "Radha")

    # removes duplicates from "Job" Column
    df.remove_duplicates("Job")

    # Converts all first names to uppercase
    df.apply_func("First Name", lambda s: s.upper())  

    # Export data to JSON
    json_data = df.export_json("ppl.json")
    print(json_data)

    # counts palindromes in the csv file
    df.count_valid_palindromes()

    # show the summary of the csv file after analysis
    df.summerize()

if __name__ == "__main__":
    run_test1()
    run_test2()
    run_test3()
    run_test4()

```

## 🦀 `Rust` and `Go` implementation comming soon, just for fun and better performance 💪
