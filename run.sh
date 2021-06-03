if [[ -z $2 ]]; then
  echo "Please specify node count and row count..."
else
  rm -r ~/.ccm
  rm avg-latencies.txt
  ccm create test -v 3.0.24
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi
