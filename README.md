# CassBlock
## Installation
Execute line by line:
```shell
sudo su 
mkdir ~/go/src -p
cd ~/go/src # install directory 
git clone https://github.com/tigranb2/CassBlock.git
cd CassBlock && . setup.sh # installs dependencies (Go, Python, Geth, Mininet)
```

## Usage
### Configuration: 
Open config.yaml to edit the number of Mininet hosts:
```shell
nano config.yaml
```
You will see:
```yaml
network:
  topo:
    class: SingleSwitchTopo
    args:
      - 1 # set number of Mininet hosts. Should be the same as number of nodes
```    

### Running:
Execute:
```shell
. run.sh {num_of_nodes} {num_of_rows} {think_time} {think_time_variation} {description} # ex: . run.sh 50 10 80 20 50_node_10_row
# num_of_nodes is the number of sensors, Cassandra nodes, and geth nodes (same as number in config.yaml)
# think_time is the time between writes to Cassandra, in milliseconds
# think_time_varaition is the amount of time think_time may vary by at most
# description is used to name the file where performance data is stored ({description}-data.txt)
```

After running, the average, median, and 99th percentile of the throughputs and write/read latencies for Cassandra and Go-Ethereum will be printed
