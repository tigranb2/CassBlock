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
    new_line = f"/cassandra/cassandra{n}/"
    replace(file, "/cassandra/", new_line)
    new_line = f"10.0.0.{n}"
    replace(file, "localhost", new_line)

    file = f"/root/cassandra/conf{n}/cassandra-env.sh"
    new_line = f"7{n}99"
    replace(file, "7199", new_line)

    file = f"/root/cassandra/bin/cassandra{n}.in.sh"
    new_line = f"$CASSANDRA_HOME/conf{n}"
    replace(file, "$CASSANDRA_HOME/conf", new_line)

    file = f"/root/cassandra/bin/cassandra{n}"
    new_line = f"cassandra{n}.in.sh"
    replace(file, "cassandra.in.sh", new_line)


def replace(file, to_find, replace_to):
    f = open(file, 'r')
    file_data = f.read()
    f.close()

    new_data = file_data.replace(to_find, replace_to)

    f = open(file, 'w')
    f.write(new_data)
    f.close()


def main():
    shutil.copytree("/root/cassandra/", "/root/cassandra_backup/")  # backup of cassandra
    for n in range(num_of_nodes):
        create_dir(n + 1)
        edit_files(n + 1)

    #move the following to a different file!
    #shutil.rmtree("/root/cassandra/")  # remove modified cassandra folder
    #shutil.copytree("/root/cassandra_backup/", "/root/cassandra/")  # restore backup
    #shutil.rmtree("/root/cassandra_backup/")  # remove backup


if __name__ == '__main__':
    main()
