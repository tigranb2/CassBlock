#!/usr/bin/env bash
node_count=$1

# delete all geth host folders
rm -rf ~/go/src/CassBlock/blockchain/*

# make node_count # of folders
cur_node=1
while [ $cur_node -le $node_count ]; do
  mkdir -p ~/go/src/CassBlock/blockchain/ethData$cur_node/keystore
  source ~/.bashrc

  # test and set env var
  dirName=DDR${cur_node}
  if [[ -z ${!dirName} || ${!dirName} != ~/go/src/CassBlock/blockchain/ethData$cur_node ]]; then
    echo export DDR$cur_node=~/go/src/CassBlock/blockchain/ethData$cur_node >>~/.bashrc
  fi
  source ~/.bashrc

  # geth init
  geth init --datadir ${!dirName} ~/go/src/CassBlock/genesis.json

  # copy keys
  cp ~/go/src/CassBlock/keys/UTC--2021-05-07T03-16-23.709229300Z--1feb3ff7be9be6a6182e6ece317a043a4f0337ab ~/go/src/CassBlock/blockchain/ethData$cur_node/keystore
  ((cur_node++))
done