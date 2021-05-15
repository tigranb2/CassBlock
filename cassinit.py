import yaml
from sys import argv
import os

num_of_nodes = int(argv[1])


def create_dir(id):
    cmd = f"cp -r ~/cassandra/conf ~/cassandra/conf{id}"
    os.system(cmd)
    cmd = f"cp ~/cassandra/bin/cassandra.in.sh ~/cassandra/bin/cassandra{id}.in.sh"
    os.system(cmd)
    cmd = f"cp ~/cassandra/bin/cassandra ~/cassandra/bin/cassandra{id}"
    os.system(cmd)


def edit_files(id):
    file = f"~/casssandra/conf{id}/cassandra.yaml"
    with open(file) as f:
        y = yaml.safe_load(f)
        y['data_file_directories'] = f'- /var/lib/cassandra/cassandra{id}/data'
        y['commit_log_directory'] = f'/var/lib/cassandra/cassandra{id}/commitlog'
        y['saved_caches_directory'] = f'/var/lib/cassandra/cassandra{id}/saved_caches'
        y['listen_address'] = '10.0.0.{id}'
        y['rpc_address'] = '10.0.0.{id}'
        yaml.dump(y, default_flow_style=False, sort_keys=False)

def main():
    for id in num_of_nodes:
        create_dir(id+1)
        edit_files(id+1)

if __name__ == '__main__':
    main()