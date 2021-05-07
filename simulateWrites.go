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
     temperature text,
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

	cassandraInit("127.0.0.1")
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
	var cassLatencies, gethLatencies []message.Latencies
	fmt.Println("Generating data...")
	randI := 0
	for {
		r++
		randI = rand.Intn(3000 - 1000) + 1000
		for id := 1; id <= 5; id++ {
			str := strconv.Itoa(randI/1000) //returns string of random int
			cassWR := message.Test_sensor{Sensor_id: id, Write: r, Temperature: (str+"km"), Speed: (str+"km"), Latencies: message.Latencies{}} //stores info for cassandra write & read 
			cassandraTest(&cassWR)
			fmt.Printf("Sensor: %v, Write: %v, write latency: %vms, read latency: %vms\n", id, r, cassWR.Latencies.WriteLatency, cassWR.Latencies.ReadLatency)
			cassLatencies = append(cassLatencies, cassWR.Latencies)
		}

		if r % gethWriteFrequency == 0 {
			gethWR := message.Latencies{} //stores info for geth write & read
			gethTest("ws://10.0.0.1:8101", "metadata", &gethWR)
			gethLatencies = append(gethLatencies, gethWR)
			fmt.Printf("Go-Ethereum - write latency: %vms, read latency: %vms\n", gethWR.WriteLatency, gethWR.ReadLatency)
		}

		if r == 100 {
			break
		}
		time.Sleep(time.Duration(randI)*time.Millisecond) //sleeps for 1 - 3 seconds
	}

	cassWriteLatency, cassReadLatency := average(cassLatencies)
	gethWriteLatency, gethReadLatency := average(gethLatencies)
	fmt.Printf("\nCASSANDRA:\n	Average write latency: %vms\n	Average read latency: %vms\n", cassWriteLatency, cassReadLatency)
	fmt.Printf("\nGO-ETHEREUM:\n	Average write latency: %vms\n	Average read latency: %vms\n", gethWriteLatency, gethReadLatency)
}

func cassandraTest(data *message.Test_sensor) {
	s := time.Now().UnixNano()/1000000
	//create new row in test_table
	if err := Session.Query("INSERT INTO test_sensor(sensor_id,write,temperature,speed) VALUES(?, ?, ?, ?)", data.Sensor_id, data.Write, data.Temperature, data.Speed).Exec(); err != nil {
		fmt.Println(err)
	}

	data.Latencies.WriteLatency = int(time.Now().UnixNano()/1000000 - s)

	s = time.Now().UnixNano()/1000000
	//read new row in test_table
	if err := Session.Query(`SELECT speed FROM test_sensor WHERE sensor_id = ? AND write = ?`, data.Sensor_id, data.Write).Exec(); err != nil {
		fmt.Println(err)
	}	

	data.Latencies.ReadLatency = int(time.Now().UnixNano()/1000000 - s)
}

func gethTest(connect, msg string, gethWR *message.Latencies) {
	s := time.Now().UnixNano()/1000000
	//writes data into geth transaction
	tx := fmt.Sprintf("eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1,data:web3.toHex('%v')})", msg)
	output, err := exec.Command("geth", "attach", connect, "--exec", tx).CombinedOutput() 
	if err != nil {
		fmt.Println(err)
		return
	}

	transactionID := string(output)
	gethWR.WriteLatency =  int(time.Now().UnixNano()/1000000 - s)

	s = time.Now().UnixNano()/1000000
	tx = fmt.Sprintf("eth.getTransaction(%v)", transactionID)
	exec.Command("geth", "attach", connect, "--exec", tx).Run() 

	gethWR.ReadLatency = int(time.Now().UnixNano()/1000000 - s)
}



func average(arr []message.Latencies) (float32, float32) {
	var writeSum, readSum int
	for _, element := range arr {
		writeSum += element.WriteLatency
		readSum += element.ReadLatency
	}

	avgWriteLatency := float32(writeSum) / float32(len(arr))
	avgReadLatency := float32(readSum) / float32(len(arr))
	return avgWriteLatency, avgReadLatency
}