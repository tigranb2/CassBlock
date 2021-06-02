init_keyspace = ('''ccm node1 cqlsh -x "CREATE KEYSPACE test_keyspace WITH replication '''
                 '''= {'class': 'SimpleStrategy', 'replication_factor': '1'};"''')
init_table = ('ccm node1 cqlsh -x "USE test_keyspace; CREATE TABLE test_sensor ( '
              'sensor_id int,row int,temperature text,speed text,PRIMARY KEY ((sensor_id), row));"')
