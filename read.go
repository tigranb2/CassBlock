/*
func cassandraRead(sensor_id int, time_lower time.Time, time_upper time.Time) []test_sensor {
	var data []test_sensor
	m := map[string]interface{}{}
	query := fmt.Sprintf("SELECT * FROM test_sensor WHERE sensor_id=%v AND collect_time>='%s' AND collect_time<='%s'", sensor_id, time_lower.Format("2006-01-02 15:04:05.000"), time_upper.Format("2006-01-02 15:04:05.000"))
	
	//read from specifed range in car_stats
	iterable := Session.Query(query).Iter()
	for iterable.MapScan(m) {
		data = append(data, test_sensor{
			sensor_id: m["sensor_id"].(int),
			Collect_time: m["collect_time"].(time.Time),
			temperature: m["temperature"].(string),
			speed: m["speed"].(string),
		})
		m = map[string]interface{}{}
	}
	return data
}
*/