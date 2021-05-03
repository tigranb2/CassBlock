if [[ -z $1 ]]; then
  echo "Please specify the geth write frequency"
else
  chmod 700 simulateWrites
    python3 start.py $1
    python3 analysis.py
  done
  cat anaylsis.txt
fi