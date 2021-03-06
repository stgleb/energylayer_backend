package main

import (
	. "../../../storage"
	influx "../../influx_storage"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	//"math/rand"
	"strconv"
	"time"
)

const (
	DB_NAME   = "energylayer"
	PASSWORD = "1234"
	USER_NAME = "test"
)

func parseMeasurements(result []client.Result) ([]Measurement, error) {
	measurements := make([]Measurement, 0, 100)

	for _, r := range result {
		for _, row := range r.Series {
			m := make(map[string]int)

			for _, val := range row.Values {
				var timestamp int64
				// Convert tuple of measurements to map <column: value>
				for index, column := range row.Columns {
					if index == 0 {
						t, err := time.Parse(time.RFC3339, val[index].(string))
						if err != nil {
							log.Fatal(err)
						}
						timestamp = t.Unix()
					} else {
						if i, err := strconv.ParseInt(fmt.Sprintf("%s", val[index]), 10, 64); err == nil {
							m[column] = int(i)
						}
					}
				}

				// Convert map to measurement
				measurement := Measurement{
					Current:     m[CURRENT],
					Voltage:     m[VOLTAGE],
					Power:       m[POWER],
					Temperature: m[TEMPERATURE],
					Gpio:        m[GPIO],
					DeviceId:    m[influx.DEVICE_ID],
					Timestamp:   timestamp,
				}

				measurements = append(measurements, measurement)
			}
		}
	}

	return measurements, nil
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: DB_NAME,
	}

	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func main() {
	// Make client
	fmt.Println("Connection to database")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: USER_NAME,
		Password: PASSWORD,
	})

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	//query := client.Query{
	//	Command:  fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", DB_NAME),
	//	Database: DB_NAME,
	//}
	//
	//c.Query(query)
	//query = client.Query{
	//	Command:  fmt.Sprintf("USE %s", DB_NAME),
	//	Database: DB_NAME,
	//}
	//c.Query(query)

	// Create a new point batch
	//bp, err := client.NewBatchPoints(client.BatchPointsConfig{
	//	Database:  DB_NAME,
	//	Precision: "s",
	//})
	//
	//if err != nil {
	//	log.Printf("Error: %s", err.Error())
	//}
	//
	//rand.Seed(11414)
	//for i := 0; i < 10; i++ {
	//	// Create a point and add to batch
	//	tags := map[string]string{"device_id": "abcd" + string(rand.Int())}
	//	fields := map[string]interface{}{
	//		VOLTAGE:     rand.Int() % 100,
	//		TEMPERATURE: rand.Int() % 100,
	//		CURRENT:     rand.Int() % 100,
	//		POWER:       rand.Int() % 100,
	//		GPIO:        rand.Int() % 100,
	//	}
	//	pt, err := client.NewPoint("data", tags, fields, time.Now())
	//
	//	if err != nil {
	//		log.Printf("Error: %s", err.Error())
	//	}
	//
	//	bp.AddPoint(pt)
	//}
	//
	//// Write the batch
	//c.Write(bp)
	fmt.Printf("Start querying data")
	result, err := queryDB(c, fmt.Sprintf("SELECT * FROM data"))

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	fmt.Println("Start processing data")
	measurements, err := parseMeasurements(result)
	fmt.Println(measurements)

	fmt.Printf("Start querying data")
	result, err = queryDB(c, fmt.Sprintf("SELECT * FROM data"))

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	fmt.Println("Start processing data")
	measurements, err = parseMeasurements(result)
	fmt.Println(measurements)
}
