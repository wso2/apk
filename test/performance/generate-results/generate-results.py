import json
import csv
import getpass
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("-n", "--name", help="Test name")
args = parser.parse_args()

username = getpass.getuser()

# Specify which columns to include in the results_summary CSV file
columns = ['user_count', 'payload_size', 'samples', 'mean', 'throughput', 'errors', 'errorPercentage', 'p90', 'p95', 'p99']
columns_display = ['Concurrent Users', 'Message Size (Bytes)', 'Total requests', 'Average Response Time (ms)', 'Throughput (Requests/sec)', 'Error %', 'Error Count', '90th Percentile of Response Time (ms)', '95th Percentile of Response Time (ms)', '99th Percentile of Response Time (ms)', "Adapter CPU", "Adapter Memory", "Enforcer Memory", "Enforcer Memory", "Router Memory", "Router Memory", "Nginx Memory", "Nginx Memory"]

# Specify the values for payload_size and user_count
payload_sizes = ['50B', '1024B', '10240B', '102400B']
payload_sizes_display = {'50B': '50B', '1024B': '1KiB', '10240B': '10KiB', '102400B': '100KiB'}
user_counts = [10, 50, 100, 200, 500, 1000]

def read_resource_usage(filename):
    with open(filename, 'r') as file:
        data = file.read()
        
        rows = data.split("\n")[1:-1]
        rows_dict = {}
        for row in rows:
            _, _, name, cpu, memory = row.split()
            rows_dict[name] = {'cpu': cpu, 'memory': memory}

        names = ["adapter", "enforcer", "router"]
        row = []
        for name in names:
            row.append(rows_dict[name]['cpu'])
            row.append(rows_dict[name]['memory'])
        return row

# Open the CSV file for writing or appending
with open(f'results_summary-{args.name}.csv', 'w', newline='') as csv_file:

    # Create a CSV writer object
    writer = csv.writer(csv_file)

    # Check if the file is empty
    is_empty = csv_file.tell() == 0

    # Loop through each combination of payload_size and user_count
    for payload_size in payload_sizes:
        for user_count in user_counts:
            # Generate the filename for the current combination of payload_size and user_count
            results_dir = f"/home/{username}/results/{args.name}/passthrough/1g_heap/{user_count}_users/{payload_size}/0ms_sleep"

            with open(f"{results_dir}/results-measurement-summary.json", 'r') as json_file:
                # Load the JSON data
                data = json.load(json_file)

                # Extract the HTTP Request data
                http_request = data['HTTP Request']

                # Extract the values for the selected columns
                row = [user_count, payload_sizes_display[payload_size]] + [http_request[col] for col in columns[2:]] + read_resource_usage(f"{results_dir}/resources-10min.txt")

                # Write the header row if the file is empty
                if is_empty:
                    writer.writerow(columns_display)
                    is_empty = False

                # Write the data row
                writer.writerow(row)
