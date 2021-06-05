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
            results.append(float(line.split()[1]))

    return results


def average(data):
    return sum(data)/len(data)


def main():
    f = read_file("avg-latencies.txt")
    throughputs = extract_data(f, "Throughput:")
    throughput_avg = round(average(throughputs), 2)
    throughput_median = round(percentile(array(throughputs), 50), 2)
    throughput_99th = round(percentile(array(throughputs), 99), 2)

    cass_writes = extract_data(f, "CassW:")
    cass_write_avg = round(average(cass_writes), 2)
    cass_write_median = round(percentile(array(cass_writes), 50), 2)
    cass_write_99th = round(percentile(array(cass_writes), 99), 2)

    cass_reads = extract_data(f, "CassR:")
    cass_read_avg = round(average(cass_reads), 2)
    cass_read_median = round(percentile(array(cass_reads), 50), 2)
    cass_read_99th = round(percentile(array(cass_reads), 99), 2)

    geth_writes = extract_data(f, "GethW:")
    geth_write_avg = round(average(geth_writes), 2)
    geth_write_median = round(percentile(array(geth_writes), 50), 2)
    geth_write_99th = round(percentile(array(geth_writes), 99), 2)

    geth_reads = extract_data(f, "GethR:")
    geth_read_avg = round(average(geth_reads), 2)
    geth_read_median = round(percentile(array(geth_reads), 50), 2)
    geth_read_99th = round(percentile(array(geth_reads), 99), 2)

    data = f"Throughput: avg: {throughput_avg}     median: {throughput_median}     99th: {throughput_99th}\n" \
           f"Cassandra:\n        " \
           f"Writes:     avg: {cass_write_avg}     median: {cass_write_median}     99th: {cass_write_99th}\n        " \
           f"Reads:      avg: {cass_read_avg}     median: {cass_read_median}     99th: {cass_read_99th}\n" \
           f"Go-Ethereum:\n        " \
           f"Writes:     avg: {geth_write_avg}     median: {geth_write_median}     99th: {geth_write_99th}\n        " \
           f"Reads:      avg: {geth_read_avg}     median: {geth_read_median}     99th: {geth_read_99th}\n"

    file_name = f"{description}-data.txt"
    f = open(file_name, "a")
    f.write(data)

main()
