import shutil
import os
from sys import argv
from time import sleep

num_of_nodes = int(argv[1])
run_mode = str(argv[2])

init_keyspace = ('''~/cassandra/bin/cqlsh -e "CREATE KEYSPACE test_keyspace WITH replication '''
                 '''= {'class': 'SimpleStrategy', 'replication_factor': '1'};" 127.0.0.1''')
init_table = ('~/cassandra/bin/cqlsh -e "USE test_keyspace; CREATE TABLE test_sensor '
              'sensor_id int,row int,temperature text,speed text,PRIMARY KEY ((sensor_id), row));" 127.0.0.1')


def start_cluster():
    if run_mode == "-n":
        os.system("ccm create test -v 3.11.10")

    cmd = f"ccm populate -n {num_of_nodes}"
    os.system(cmd)
    os.system("ccm start --root")