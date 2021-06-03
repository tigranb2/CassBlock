if [[ -z $2 ]]; then
  echo "Please specify node count and row count..."
else
  rm -r ~/.ccm
  ccm create test -v 3.0.24
  ccm add node1 -i 10.0.0.1 -j 7100 -s
  ccm add node2 -i 10.0.0.2 -j 7200 -s
  ccm add node3 -i 10.0.0.3 -j 7300 -s
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi