package main

import (
	"CassBlock/message"
	"fmt"
	"os"
	"os/exec"
	"time"
	"math/rand"
	"strconv"
	"github.com/gocql/gocql"
)

/* 
CREATE KEYSPACE test_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

CREATE TABLE test_sensor (
     sensor_id int,
	 write int,
     local_write_time bigint,
     speed text,
     PRIMARY KEY ((sensor_id), write)
   ) ;
*/

var Session *gocql.Session
var r = 0 //how many times data for all sensors has been written

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("Please provide frequency for writing to Geth...")
		return
	}
	gethWriteFrequency, _ := strconv.Atoi(arguments[1])

	cassandraInit("10.0.0.1")
	simulateWrites(gethWriteFrequency)
}

func cassandraInit(CONNECT string){
	var err error
	cluster := gocql.NewCluster(CONNECT) //connect to cassandra database
	cluster.Keyspace = "test_keyspace"
	Session, err = cluster.CreateSession() 
	if err != nil {
		fmt.Print(err)
	}
}

func simulateWrites(gethWriteFrequency int) {
	fmt.Println("Generating data...")
	randI := 0
	for {
		r++
		randI = rand.Intn(3000 - 1000) + 1000
		for id := 1; id <= 5; id++ {
			str := strconv.Itoa(randI/1000) //returns string of random int
			data := message.Test_sensor{Sensor_id: id, Write: r, Local_write_time: time.Now().UnixNano()/1000, Speed: (str+"km")}
			cassandraWrite(data)
			latencyTest(id, r)
		}

		if r % gethWriteFrequency == 0 {
			fmt.Println("geth write")
			//gethWrite("ws://10.0.0.1", "metadata")
		}

		time.Sleep(time.Duration(randI)*time.Millisecond) //sleeps for 1 - 3 seconds
	}
}

func cassandraWrite(data message.Test_sensor) {
	//create new row in test_table
	if err := Session.Query("INSERT INTO test_sensor(sensor_id,write,local_write_time,speed) VALUES(?, ?, ?, ?)", data.Sensor_id, data.Write, data.Local_write_time, data.Speed).Exec(); err != nil {
		fmt.Println(err)
	}
}

func latencyTest(sensorId, write int) {
	var local_write_time, WRITETIME int
	Session.Query(`SELECT local_write_time FROM test_sensor WHERE sensor_id = ? AND write = ?`, sensorId, write).Scan(&local_write_time)
	Session.Query(`SELECT WRITETIME FROM test_sensor WHERE sensor_id = ? AND write = ?`, sensorId, write).Scan(&WRITETIME)
	fmt.Println(local_write_time, WRITETIME, WRITETIME - local_write_time)

}

func gethWrite(connect, msg string){
	tx := fmt.Sprintf("eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1,data:web3.toHex('%v')})", msg)
	_, err := exec.Command("geth", "attach", connect, "--exec", tx).CombinedOutput() 
	if err != nil {
		fmt.Println(err)
	}
}