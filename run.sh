if [[ -z $7 ]]; then
  echo "Please specify sensor node, Cassandra, Go-Ethereum, and row counts, rate parameter (λ), run mode (-c or -g), and test description.."
else
  rm -r ~/.ccm
  rm avg-latencies.txt
  ccm create test -v 3.0.24
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2 $3 $4 $5 $6
  python3 analysis.py $6 $7
  cat $7-data.txt
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi