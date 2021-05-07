package message

type Latencies struct {
	WriteLatency, ReadLatency float32
}

type Test_sensor struct {
	Sensor_id, Write int
	Temperature, Speed string
	Latencies Latencies
}
