from sys import argv
from numpy import percentile, array
description = str(argv[1])


def read_file(file):
    with open(file) as f:
        return f.readlines()


def extract_data(f, string_to_search):
    results = []
    for line in f:
        if string_to_search in line:
            results.append(line.split()[1])

    return results


def average(data):
    return sum(data)/len(data)


def main():
    f = read_file("avg-latencies.txt")
    cass_writes = extract_data(f, "CassW:")
    cass_write_avg = average(cass_writes)
    cass_write_median = percentile(array(cass_writes), 50)
    cass_write_99th = percentile(array(cass_writes), 99)

    cass_reads = extract_data(f, "CassR:")
    cass_read_avg = average(cass_reads)
    cass_read_median = percentile(array(cass_reads), 50)
    cass_read_99th = percentile(array(cass_reads), 99)

    geth_writes = extract_data(f, "GethW:")
    geth_write_avg = average(geth_writes)
    geth_write_median = percentile(array(geth_writes), 50)
    geth_write_99th = percentile(array(geth_writes), 99)

    geth_reads = extract_data(f, "GethR:")
    geth_read_avg = average(geth_reads)
    geth_read_median = percentile(array(geth_reads), 50)
    geth_read_99th = percentile(array(geth_reads), 99)

    data = f"Cassandra:\n        " \
           f"Writes:    avg: {cass_write_avg}    median: {cass_write_median}    99th: {cass_write_99th}\n        " \
           f"Reads:     avg: {cass_read_avg}    median: {cass_read_median}    99th: {cass_read_99th}\n" \
           f"Geth:\n        " \
           f"Writes:    avg: {geth_write_avg}    median: {geth_write_median}    99th: {geth_write_99th}\n        " \
           f"Reads:     avg: {geth_read_avg}    median: {geth_read_median}    99th: {geth_read_99th}"

    file_name = f"{description}-data.txt"
    f = open(file_name, "a")
    f.write(data)
