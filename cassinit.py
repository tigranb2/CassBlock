import yaml
import shutil
from sys import argv

num_of_nodes = int(argv[1])


def create_dir(n):
    dst = f"/root/cassandra/conf{n}"
    shutil.copytree("/root/cassandra/conf", dst)
    dst = f"/root/cassandra/bin/cassandra{n}.in.sh"
    shutil.copy("/root/cassandra/bin/cassandra.in.sh", dst)
    dst = f"/root/cassandra/bin/cassandra{n}"
    shutil.copy("/root/cassandra/bin/cassandra", dst)


def edit_files(n):
    file = f"/root/cassandra/conf{n}/cassandra.yaml"
    yaml_replace(file, n)

    file = f"/root/cassandra/conf{n}/cassandra-env.sh"
    new_line = f"7{n}99"
    replace(file, "7199", new_line)

    file = f"/root/cassandra/conf{n}/log4j-server.properties"
    new_line = f"/var/log/cassandra/cassandra{n}/system.log"
    replace(file, "/var/log/cassandra/system.log", new_line)
    replace(file, "INFO.stdout.R", "DEBUG.stdout.R")

    file = f"/root/cassandra/bin/cassandra{n}.in.sh"
    new_line = f"$CASSANDRA_HOME/conf{n}"
    replace(file, "$CASSANDRA_HOME/conf", new_line)

    file = f"/root/cassandra/bin/cassandra{n}"
    new_line = f"/cassandra{n}.in.sh"
    replace(file, "/cassandra.in.sh", new_line)


def yaml_replace(file, n):
    # read yaml file
    with open(file) as f:
        y = yaml.load(f)

    dst = f'/var/lib/cassandra/cassandra{n}/data'
    y['data_file_directories'] = dst
    dst = f'/var/lib/cassandra/cassandra{n}/commitlog'
    y['commitlog_directory'] = dst
    dst = f'/var/lib/cassandra/cassandra{n}/saved_caches'
    y['saved_caches_directory'] = dst
    dst = f'10.0.0.{n}'
    y['listen_address'] = dst
    dst = f'10.0.0.{n}'
    y['rpc_address'] = dst

    # write updated parameters to file
    with open(file, "w") as f:
        yaml.dump(y, f)


def replace(file, to_find, replace_to):
    f = open(file, 'r')
    file_data = f.read()
    f.close()

    new_data = file_data.replace(to_find, replace_to)

    f = open(file, 'w')
    f.write(new_data)
    f.close()


def main():
    #shutil.copytree("/root/cassandra/", "/root/cassandra_backup/")  # backup of cassandra
    for n in range(num_of_nodes):
        create_dir(n + 1)
        edit_files(n + 1)

    #move the following to a different file!
    #shutil.rmtree("/root/cassandra/")  # remove modified cassandra folder
    #shutil.copytree("/root/cassandra_backup/", "/root/cassandra/")  # restore backup
    #shutil.rmtree("/root/cassandra_backup/")  # remove backup


if __name__ == '__main__':
    main()
