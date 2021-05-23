from sys import argv

num_of_nodes = int(argv[1])
nodes = [f"10.0.0.{i}" for i in range(1, num_of_nodes + 1)]
miner_thread = 3
network_id = 714715
port = 30303

gen_enode = (f"ENODE_ADDRESS=\"enode://$(bootnode -nodekey $DDR1/geth/nodekey -writeaddress)"
             f"@{nodes[0]}:{port}\"")

miner_start = (
    f"geth --nodiscover --ipcdisable --networkid {network_id} "
    f"--syncmode 'full' --port {port} --datadir=$DDR1 "
    "--ws --wsaddr 10.0.0.1 --wsport 8101 --rpcapi eth,web3,personal,net,admin,miner "
    f"--gasprice '1' --mine --minerthreads {miner_thread} "
    "--etherbase='0x1feb3ff7be9be6a6182e6ece317a043a4f0337ab' "
    "--unlock '0x1feb3ff7be9be6a6182e6ece317a043a4f0337ab' -allow-insecure-unlock --password ./password.sec &"
)

node_start = (
    f"nohup geth --nodiscover --ipcdisable --networkid {network_id} --port {port} "
    "--syncmode 'full' --datadir=$DDR{n} --bootnodes $ENODE_ADDRESS "
    "--ws --wsaddr 10.0.0.{n} --wsport 8101 --rpcapi eth,web3,personal,net,admin,miner "
    f"--gasprice '1' --mine --minerthreads {miner_thread} "
    "--etherbase='0x1feb3ff7be9be6a6182e6ece317a043a4f0337ab' "
    "--unlock '0x1feb3ff7be9be6a6182e6ece317a043a4f0337ab' -allow-insecure-unlock --password ./password.sec &"
)
