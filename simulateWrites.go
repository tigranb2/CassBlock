package main

import (
	"CassBlock/message"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/gocql/gocql"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

/*
CREATE KEYSPACE test_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

CREATE TABLE test_sensor (
     sensor_id int,
	 row int,
     round int,
     speed double,
     PRIMARY KEY ((sensor_id), row)
WITH CLUSTERING ORDER BY (row DESC)
   ) ;
*/

var (
	Session                                                        *gocql.Session
	ip                                                             string
	operations, round                                              int
	lambda                                                         float64
	cassRLatencies, gethRLatencies, cassWLatencies, gethWLatencies []int
	gethTxs                                                        []string
	gethMode                                                       bool
)

func main() {
	arguments := os.Args
	if len(arguments) < 5 {
		fmt.Println("Please specify node count, row count, rate parameter (λ), and run mode (-c or -g)...")
		return
	}
	id, _ := strconv.Atoi(arguments[1])              //numbers of rows per sensor
	rowCount, _ := strconv.Atoi(arguments[2])        //numbers of rows per sensor
	lambda, _ = strconv.ParseFloat(arguments[3], 64) //time between messages
	gethMode = arguments[4] == "-g"                  //determines whether geth will be used or not

	ip = fmt.Sprintf("10.0.0.%v", id)

	cassandraInit(ip) //connect to cassandra database
	simulateWrites(id, rowCount)
}

func cassandraInit(CONNECT string) {
	var err error
	cluster := gocql.NewCluster(CONNECT)
	cluster.Keyspace = "test_keyspace"
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Print(err)
	}
}

func simulateWrites(id, rowCount int) {
	var (
		currentWrites []message.Test_sensor
		throughput    float32
	)
	row := 0

	fmt.Println("Generating data...")
	start := time.Now().Unix()
	go read(id, rowCount) //goroutine that reads randomly
	for {
		row++
		round++
		value := rand.NormFloat64() //returns random value from N(0, 1)

		data := message.Test_sensor{Sensor_id: id, Row: row, Round: round, Speed: value} //stores info for cassandra write & read
		cassWLatency := cassandraWrite(data)                                                //writes to Cassandra
		fmt.Printf("Cassandra write - Sensor: %v, Row: %v, latency: %vms\n", id, row, cassWLatency)
		cassWLatencies = append(cassWLatencies, cassWLatency)
		currentWrites = append(currentWrites, data)

		if gethMode && len(currentWrites)%rowCount == 0 {
			min, max := extrema(currentWrites)
			metadata := encode(id, min, max)                        //byte array of data
			gethWLatency := gethWrite("ws://"+ip+":8101", metadata) //writes to Geth
			fmt.Printf("Go-Ethereum write - latency: %vms\n\n", gethWLatency)
			gethWLatencies = append(gethWLatencies, gethWLatency)

			currentWrites = []message.Test_sensor{}
			row = 0 //to overwrite rows 1..rowCount
		} else if len(currentWrites)%rowCount == 0 {
			currentWrites = []message.Test_sensor{}
			row = 0 //to overwrite rows 1..rowCount
		}

		if round == 100 { //exit after writing 100 times to Cassandra
			timeTaken := time.Now().Unix() - start
			throughput = float32(int64(operations) / timeTaken)
			break
		}

		sleepTime := rand.ExpFloat64() / lambda //return value from exponential dist.
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	cassWriteLatency := average(cassWLatencies)
	cassReadLatency := average(cassRLatencies)
	writeString := ""
	if gethMode {
		gethWriteLatency := average(gethWLatencies)
		gethReadLatency := average(gethRLatencies)
		writeString = fmt.Sprintf("Throughput: %v\nCassW: %v\nCassR: %v\nGethW: %v\nGethR: %v\n", throughput, cassWriteLatency, cassReadLatency, gethWriteLatency, gethReadLatency)
	} else {
		writeString = fmt.Sprintf("Throughput: %v\nCassW: %v\nCassR: %v\n", throughput, cassWriteLatency, cassReadLatency)
	}

	writeToFile("avg-latencies.txt", writeString) //write data to avg-latencies.txt
}

func read(id, rowCount int) {
	var currentReads []message.Test_sensor
	for {
		if round > 1 { //read from cassandra
			start := time.Now().UnixNano() / 1000000

			m := map[string]interface{}{}
			var rows, previousRows []message.Test_sensor
			iterable := Session.Query("SELECT * FROM test_sensor WHERE sensor_id = ?;", id).Iter() //read all cassandra rows for this sensor
			for iterable.MapScan(m) {
				row := message.Test_sensor{Sensor_id: m["sensor_id"].(int), Row: m["row"].(int), Round: m["round"].(int), Speed: m["speed"].(float64)}
				rows = append(rows, row)
				m = map[string]interface{}{}
			}

			latest := message.Test_sensor{}
			for _, row := range rows {
				if row.Round < round {
					previousRows = append(previousRows, row) //collection of rows from previous rounds
				} else if row.Round == round {
					latest = row //latest write
				}
			}

			if latest.Sensor_id != 0 { //latest write found
				previousSetMedian := median(previousRows)
				if latest.Speed >= previousSetMedian-2 && latest.Speed <= previousSetMedian+2 { //check if latest write is within 2 z-score of previous set median
					cassRLatency := int(time.Now().UnixNano()/1000000 - start) //cassandra read latency
					cassRLatencies = append(cassRLatencies, cassRLatency)
					currentReads = append(currentReads, latest)
					fmt.Printf("Cassandra read: %v with latency: %vms\n", latest.Speed, cassRLatency)
				}
			} else {
				continue
			}
			operations++
		}

		if gethMode && len(currentReads)%rowCount == 0 && len(currentReads) > 0 && len(gethTxs) > 1{ //read from geth
			start := time.Now().UnixNano() / 1000000
			latest := currentReads[len(currentReads)-1] //latest read from cassandra

			latestTx := gethTxs[len(gethTxs)-1]
			cmd := fmt.Sprintf("eth.getTransaction(%v).input", latestTx)
			output, err := exec.Command("geth", "attach", "ws://"+ip+":8101", "--exec", cmd).CombinedOutput() //reads transaction
			if err != nil {
				fmt.Println("Transaction not found...")
			}
			data := decode(string(output)) //struct w/ min and max values written
			min := data.Min
			max := data.Max

			if latest.Speed >= min && latest.Speed <= max { //basic outlier test
				gethRLatency := int(time.Now().UnixNano()/1000000 - start) //geth read latency
				gethRLatencies = append(gethRLatencies, gethRLatency)
				fmt.Printf("Go-Ethereum read: %v with latency: %vms\n", latest, gethRLatency)
				operations++
			}

			currentReads = []message.Test_sensor{}
		}

		iaTime := rand.ExpFloat64() / lambda //return value from exponential dist.
		time.Sleep(time.Duration(iaTime) * time.Second)
	}
}

func cassandraWrite(data message.Test_sensor) int {
	start := time.Now().UnixNano() / 1000000
	//create new row in test_table
	if err := Session.Query("INSERT INTO test_sensor(sensor_id,row,round,speed) VALUES(?, ?, ?, ?)", data.Sensor_id, data.Row, data.Round, data.Speed).Exec(); err != nil {
		fmt.Println(err)
	}
	operations++
	return int(time.Now().UnixNano()/1000000 - start) //write latency
}

func gethWrite(connect string, data string) int {
	start := time.Now().UnixNano() / 1000000
	tx := fmt.Sprintf("eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:1,data:'0x%v'})", data) //writes data into geth transaction
	output, err := exec.Command("geth", "attach", connect, "--exec", tx).CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	transactionID := string(output)
	gethTxs = append(gethTxs, transactionID)
	operations++
	return int(time.Now().UnixNano()/1000000 - start) //write latency
}

func encode(sensorId int, min, max float64) string {
	data := message.BlockchainData{SensorId: sensorId, Min: min, Max: max, Timestamp: time.Now().Unix()}
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(data); err != nil { //encode struct to byte array
		return ""
	}
	return hex.EncodeToString(buf.Bytes()) //return hex of byte array
}

func decode(src string) message.BlockchainData {
	strEnc :=  src[3 : len(src)-1] //remove quotation marks and 0x
	dataBytes, _:= hex.DecodeString(strEnc)
	buf := bytes.NewReader(dataBytes)
	data := message.BlockchainData{}
	gob.NewDecoder(buf).Decode(&data) //decode byte array to struct
	return data
}

func extrema(arr []message.Test_sensor) (float64, float64) { //calculate rows with min and max speed in array
	min, max := arr[0].Speed, arr[0].Speed
	for i := 1; i < len(arr); i++ {
		if arr[i].Speed < min {
			min = arr[i].Speed
		} else if arr[i].Speed > max {
			max = arr[i].Speed
		}
	}

	return min, max
}

func median(arr []message.Test_sensor) float64 { //find median speed value of array
	values := []float64{}
	for _, element := range arr {
		values = append(values, element.Speed)
	}

	sort.Float64s(values)
	if len(values)%2 == 0 {
		return (values[len(values)/2] + values[(len(values)/2)-1]) / 2
	} else {
		return values[(len(values)-1)/2]
	}
}

func average(arr []int) float32 { //return average of array
	sum := 0
	for _, element := range arr {
		sum += element
	}
	return float32(sum) / float32(len(arr))
}

func writeToFile(file, writeString string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //creates file if it doesn't exist
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.WriteString(writeString) //writes latency data to file
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = f.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
