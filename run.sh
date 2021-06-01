if [[ -z $2 ]]; then
  echo "Please specify node count, row count, and -n if running for the first time"
else
  if [[ $3 = "-n" ]]; then
    ccm create test -v 3.0.24
  fi
  go build ./simulateWrites.go
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh $1
  python3 start.py $1 $2
  killall java
  rm -r ~/go/src/CassBlock/blockchain/
  killall java
  killall geth
fi