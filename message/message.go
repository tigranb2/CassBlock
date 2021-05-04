package message

type Latencies struct {
	WriteLatency, ReadLatency int
}

type Test_sensor struct {
	Sensor_id, Write int
	Temperature, Speed string
	Latencies Latencies
}
