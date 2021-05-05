if [[ -z $1 ]]; then
  echo "Please specify the geth write frequency"
else
  chmod 700 simulateWrites
  python3 start.py $1
fi
