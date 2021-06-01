init_keyspace = ('''~/cassandra/bin/cqlsh -e "CREATE KEYSPACE test_keyspace WITH replication '''
                 '''= {'class': 'SimpleStrategy', 'replication_factor': '1'};" 127.0.0.1''')
init_table = ('~/cassandra/bin/cqlsh -e "USE test_keyspace; CREATE TABLE test_sensor '
              'sensor_id int,row int,temperature text,speed text,PRIMARY KEY ((sensor_id), row));" 127.0.0.1')