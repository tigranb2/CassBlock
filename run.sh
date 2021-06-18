if [[ -z $5 ]]; then
  echo "Please specify node count, row count, rate parameter (λ), run mode (-c or -g), and test description.."
else
  rm -r ~/.ccm
  rm avg-latencies.txt
  ccm create test -v 3.0.24
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2 $3 $4
  python3 analysis.py $4 $5
  cat $5-data.txt
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi