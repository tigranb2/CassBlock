package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"math/rand"
	"strconv"
	"github.com/gocql/gocql"
)

/*
CREATE TABLE test_sensor (
     sensor_id int,
	 write int,
     write_time int,
     speed text,
     PRIMARY KEY ((sensor_id), write)
   ) ;
*/

var Session *gocql.Session
var r = 0 //how many times data for all sensors has been written

type test_sensor struct {
	sensor_id int
	write int
	write_time int64 //unix time of write, in nanoseconds 
	speed string
}

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
			data := test_sensor{id, r, time.Now().UnixNano()/1000, (str+"km")}
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

func cassandraWrite(data test_sensor) {
	//create new row in test_table
	if err := Session.Query("INSERT INTO test_sensor(sensor_id,write,write_time,speed) VALUES(?, ?, ?, ?)", data.sensor_id, data.write, data.write_time, data.speed).Exec(); err != nil {
		fmt.Println(err)
	}
}

func latencyTest(sensorId, write int) {
	/*
		Read data based on sensorId & write
		Read WRITETIME of that row
		Compare
		Return write latency
	*/

	//var write_time int
	query := Session.Query(`SELECT write_time FROM tweet WHERE sensor_id = ?, write = ?`, sensorId, write)
	fmt.Println(query)

}

func gethWrite(connect, msg string){
	tx := fmt.Sprintf("eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1,data:web3.toHex('%v')})", msg)
	_, err := exec.Command("geth", "attach", connect, "--exec", tx).CombinedOutput() 
	if err != nil {
		fmt.Println(err)
	}
}