from sys import argv

num_of_miners = int(argv[1])
miner_thread = 3
network_id = 714715
port = 30303

miner_start = (
    f"nohup geth --nodiscover --ipcdisable --networkid {network_id} "
    f"--syncmode 'full' --port {port} --datadir geth "
    "--ws --wsaddr 10.0.0.1 --wsport 8101 --rpcapi eth,web3,personal,net,admin,miner "
    f"--gasprice '1' --mine --minerthreads {miner_thread} "
    "--unlock 0 -allow-insecure-unlock --password ./password.sec"
)