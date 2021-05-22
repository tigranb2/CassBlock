if [[ -z $2 ]]; then
  echo "Please specify node count and row count"
else
  go build ./simulateWrites.go
  chmod 700 cassandra.py
  chmod 700 gethrun.sh
  chmod 700 simulateWrites

  python3 cassandra.py $1 # initializes cassandra directories
  . gethrun.sh $1 # initializes geth directories
  python3 start.py $1 $2 $3 # optional 3rd argument for mode

  rm -r ~/go/src/CassBlock/blockchain/
  rm -r /root/cassandra/
  cp -r /root/cassandra_backup/ /root/cassandra/
  rm -r root/cassandra_backup/
  killall java
  killall geth
fi