if [[ -z $1 ]]; then
  echo "Please specify the geth write frequency"
else
  go build ./simulateWrites.go
  chmod 700 cassrun.sh
  chmod 700 gethrun.sh
  chmod 700 simulateWrites
  . gethrun.sh
  python3 start.py $1 $2
  rm -r geth
  killall java
fi