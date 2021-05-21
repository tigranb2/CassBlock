package main

import (
	"CassBlock/message"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/gocql/gocql"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

/*
CREATE KEYSPACE test_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

CREATE TABLE test_sensor (
     sensor_id int,
	 row int,
     temperature text,
     speed text,
     PRIMARY KEY ((sensor_id), row)
   ) ;
*/

var Session *gocql.Session
var ip string

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("Please specify sensor id, row count...")
		return
	}
	id, _ := strconv.Atoi(arguments[1])       //numbers of rows per sensor
	rowCount, _ := strconv.Atoi(arguments[2]) //numbers of rows per sensor
	ip = fmt.Sprintf("10.0.0.%v", id)

	cassandraInit(ip)
	simulateWrites(id, rowCount)
}

func cassandraInit(CONNECT string) {
	var err error
	cluster := gocql.NewCluster(CONNECT) //connect to cassandra database
	cluster.Keyspace = "test_keyspace"
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Print(err)
	}
}

func simulateWrites(id, rowCount int) {
	var cassLatencies, gethLatencies []message.Latencies
	var s []message.Test_sensor
	write, row, randI := 0, 0, 0
	fmt.Println("Generating data...")
	for {
		row++
		write++
		randI = rand.Intn(3000-1000) + 1000
		str := strconv.Itoa(randI / 1000)                                                                                                  //returns string of random int
		cassWR := message.Test_sensor{Sensor_id: id, Row: row, Temperature: str + "km", Speed: str + "km", Latencies: message.Latencies{}} //stores info for cassandra write & read
		cassandraTest(&cassWR)
		fmt.Printf("Sensor: %v, Row: %v, write latency: %vms, read latency: %vms\n", 1, row, cassWR.Latencies.WriteLatency, cassWR.Latencies.ReadLatency)
		cassLatencies = append(cassLatencies, cassWR.Latencies)
		s = append(s, cassWR)

		if len(s)%rowCount == 0 {
			metadata := hash(s)
			gethWR := message.Latencies{} //stores info for geth write & read
			gethTest("ws://"+ip+"8101", metadata, &gethWR)
			gethLatencies = append(gethLatencies, gethWR)
			fmt.Printf("Go-Ethereum - write latency: %vms, read latency: %vms\n", gethWR.WriteLatency, gethWR.ReadLatency)
			s = []message.Test_sensor{}
			row = 0
		}

		if write == 100 {
			break
		}
		time.Sleep(time.Duration(randI) * time.Millisecond) //sleeps for 1 - 3 seconds
	}

	cassWriteLatency, cassReadLatency := average(cassLatencies)
	gethWriteLatency, gethReadLatency := average(gethLatencies)
	fmt.Printf("\nCASSANDRA:\n	Average write latency: %vms\n	Average read latency: %vms\n", cassWriteLatency, cassReadLatency)
	fmt.Printf("\nGO-ETHEREUM:\n	Average write latency: %vms\n	Average read latency: %vms\n", gethWriteLatency, gethReadLatency)
}

func cassandraTest(data *message.Test_sensor) {
	start := time.Now().UnixNano() / 1000000
	//create new row in test_table
	if err := Session.Query("INSERT INTO test_sensor(sensor_id,row,temperature,speed) VALUES(?, ?, ?, ?)", data.Sensor_id, data.Row, data.Temperature, data.Speed).Exec(); err != nil {
		fmt.Println(err)
	}

	data.Latencies.WriteLatency = int(time.Now().UnixNano()/1000000 - start)

	start = time.Now().UnixNano() / 1000000
	//read new row in test_table
	if err := Session.Query(`SELECT speed FROM test_sensor WHERE sensor_id = ? AND row = ?`, data.Sensor_id, data.Row).Exec(); err != nil {
		fmt.Println(err)
	}

	data.Latencies.ReadLatency = int(time.Now().UnixNano()/1000000 - start)
}

func gethTest(connect string, msg [32]byte, gethWR *message.Latencies) {
	start := time.Now().UnixNano() / 1000000
	//writes data into geth transaction
	tx := fmt.Sprintf("eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1,data:web3.toHex('%v')})", msg)
	output, err := exec.Command("geth", "attach", connect, "--exec", tx).CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}

	transactionID := string(output)
	fmt.Println(transactionID)
	gethWR.WriteLatency = int(time.Now().UnixNano()/1000000 - start)

	start = time.Now().UnixNano() / 1000000
	tx = fmt.Sprintf("eth.getTransaction(%v)", transactionID)
	exec.Command("geth", "attach", connect, "--exec", tx).Run()

	gethWR.ReadLatency = int(time.Now().UnixNano()/1000000 - start)
}

func hash(arr []message.Test_sensor) [32]byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(arr)
	return sha256.Sum256(buf.Bytes())
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
