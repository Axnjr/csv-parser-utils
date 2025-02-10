from main import CSV_Utils_Py
import ctypes

def run_test1():
    sample_data = [
        ['Name', 'Age', 'Score'], 
        ['Radha', '23', '85'], 
        ['Bob', '23', '90'], 
        ['Radha', '25', '78']
    ]
    test_file = 'sample.csv'
    CSV_Utils_Py.write_csv(test_file, sample_data)

def run_test2():

    df = CSV_Utils_Py("./ppl.csv")

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
    df = CSV_Utils_Py("./fruits.csv")
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

def run_go_binding_test1():
    # !TODO! not completed yet, was just playing around !!
    library = ctypes.cdll.LoadLibrary('./go_bindings/library.so')
    library.readCSV.restype = None
    # Set return type for the function (since it returns an int)
    library.countPalindromes.restype = ctypes.c_int

    # Call the function with the CSV file path
    csv_path = b"ppl.csv"  # Provide the actual CSV file path
    result = library.countPalindromes(csv_path)
    library.readCSV(csv_path)
    print("Valid Palindrome Count: ", result)

if __name__ == "__main__":

    run_test1()
    run_test2()
    run_test3()
    run_test4()
    
    # run_go_binding_test1()
   
    # fast_csv.CountValidPalindromes.argtypes = [ctypes.c_char_p]
    # fast_csv.CountValidPalindromes.restype = ctypes.c_int

    # count = fast_csv.CountValidPalindromes(b"ppl.csv")
    # print(count)

    # count = fast_csv.CountValidPalindromes(b"ppl.csv")
    # print("Palindrome Count:", count)