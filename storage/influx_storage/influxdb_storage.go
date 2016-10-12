package influxdb_storage

import (
	. "../../storage"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

const (
	DEVICE_ID = "device_id"
	PRECISION = "s"
)

type InfluxDbStorage struct {
	DbName   string // energylayer
	UserName string // root
	Password string // 1234
	Addr     string // http://localhost:8086
	Client   client.Client
}

func NewInfluxDBStorage(userName, password, dbName, addr, clientType string) (*InfluxDbStorage, error) {
	var cl client.Client
	var err error

	switch clientType {
	case "http":
		cl, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     addr,
			Username: userName,
			Password: password,
		})
		if err != nil {
			return nil, err
		}
		// TODO: add other types of clients
	default:
		return nil, errors.New("Wrong client type")
	}

	// Create database if it not exists
	query := client.Query{
		Command:  fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName),
		Database: dbName,
	}
	_, err = cl.Query(query)
	if err != nil {
		return nil, err
	}

	query = client.Query{
		Command:  fmt.Sprintf("USE %s", dbName),
		Database: dbName,
	}
	// TODO: research whether use statement is needed in Go client
	_, err = cl.Query(query)
	if err != nil {
		return nil, err
	}

	_, _, err = cl.Ping(time.Second * 10)
	if err != nil {
		return nil, err
	}

	return &InfluxDbStorage{
		DbName:   dbName,
		UserName: userName,
		Password: password,
		Addr:     addr,
		Client:   cl,
	}, nil
}

func (influx InfluxDbStorage) AddMeasurementToBatch(m Measurement, batch client.BatchPoints) error {
	tags := map[string]string{"device_id": "device_id"}
	fields := map[string]interface{}{
		CURRENT:     m.Current,
		POWER:       m.Power,
		VOLTAGE:     m.Voltage,
		TEMPERATURE: m.Temperature,
		GPIO:        m.Gpio,
	}
	point, err := client.NewPoint("measurement", tags, fields, time.Now())

	if err != nil {
		log.Printf("Error during creating point from %s", m)
		return err
	}
	batch.AddPoint(point)

	return nil
}

func (influx InfluxDbStorage) CreateMeasurements(measurements []Measurement) error {
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.DbName,
		Precision: PRECISION,
	})

	if err != nil {
		return err
	}

	for _, measurement := range measurements {
		err := influx.AddMeasurementToBatch(measurement, batch)

		if err != nil {
			return err
		}
	}
	err = influx.Client.Write(batch)
	if err != nil {
		return err
	}

	return nil
}

func (influx InfluxDbStorage) CreateMeasurement(m Measurement) error {
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.DbName,
		Precision: PRECISION,
	})

	if err != nil {
		return err
	}

	err = influx.AddMeasurementToBatch(m, batch)
	if err != nil {
		return err
	}

	return nil
}

func (influx InfluxDbStorage) queryDB(command string) ([]client.Result, error) {
	var res []client.Result

	q := client.Query{
		Command:  command,
		Database: influx.DbName,
	}

	if response, err := influx.Client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (influx InfluxDbStorage) parseMeasurements(result []client.Result) ([]Measurement, error) {
	measurements := make([]Measurement, 0, 100)

	for _, r := range result {
		for _, row := range r.Series {
			var m map[string]int

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
					DeviceId:    m[DEVICE_ID],
					Timestamp:   timestamp,
				}

				measurements = append(measurements, measurement)
			}
		}
	}

	return measurements, nil
}

func (influx InfluxDbStorage) GetMeasurements(count int) ([]Measurement, error) {
	cmd := fmt.Sprintf("SELECT * FROM measurement")
	result, err := influx.queryDB(cmd)

	if err != nil {
		return nil, err
	}
	measurements, err := influx.parseMeasurements(result)

	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func (influx InfluxDbStorage) GetDeviceById(uuid string) (Device, error) {
	return Device{}, errors.New("Not implemented")
}

func (influx InfluxDbStorage) CreateDevice(uuid, ipAddress string, userId int) error {
	return errors.New("Not implemented")
}

func (influx InfluxDbStorage) GetMeasurementsByDevice(device_id, count int) ([]Measurement, error) {
	return nil, nil
}
