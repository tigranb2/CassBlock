mkdir -p /geth/keystore
geth --datadir geth init genesis.json
cp ./keys/UTC--2021-05-07T03-16-23.709229300Z--1feb3ff7be9be6a6182e6ece317a043a4f0337ab ./geth/keystore
