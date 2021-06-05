if [[ -z $3 ]]; then
  echo "Please specify node count, row count, think time (ms), think time varaition (ms), and a test description..."
else
  rm -r ~/.ccm
  rm avg-latencies.txt
  ccm create test -v 3.0.24
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2 $3 $4
  python3 analysis.py $5
  cat $5-data.txt
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi