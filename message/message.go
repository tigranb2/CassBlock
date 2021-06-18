package message

type Test_sensor struct {
	Sensor_id, Row, Writeset int
	Speed float64
}

type BlockchainData struct {
	SensorId int
	Min, Max            float64
	Timestamp int64
}