~/cassandra/bin/cassandra1 -R
sleep 60 # waits 1 minute for cassandra to come online
  
if [[ $1 == "rerun" ]]; then
    ~/cassandra/bin/cqlsh -e "USE test_keyspace; DROP TABLE test_sensor;" 
else 
    ~/cassandra/bin/cqlsh -e "CREATE KEYSPACE test_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};"
fi
~/cassandra/bin/cqlsh -e "USE test_keyspace; CREATE TABLE test_sensor (sensor_id int,row int,temperature text,speed text,PRIMARY KEY ((sensor_id), row));"
service cassandra stop
