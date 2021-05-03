package message

type Test_sensor struct {
	Sensor_id, Write int
	Local_write_time int64 //unix time of write, in nanoseconds 
	Speed string
}

type WriteData struct {
	Sensor, Write, WriteLatency int
}