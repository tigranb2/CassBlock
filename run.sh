if [[ -z $2 ]]; then
  echo "Please specify node count, row count, and -n if running for the first time"
else
  go build ./simulateWrites.go
  chmod 700 cassandra.py
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2
  rm -r geth
  killall java
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi