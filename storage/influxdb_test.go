package storage

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	DB_NAME  = "test"
	USER     = "test"
	PASSWORD = "1234"
	ADDR     = "http://localhost:8086"
)

type FakeClient struct{}

var fakeClient FakeClient

var influx InfluxDbStorage = InfluxDbStorage{
	DbName:   DB_NAME,
	UserName: USER,
	Password: PASSWORD,
	Addr:     ADDR,
	Client:   fakeClient,
}

func (fake FakeClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Millisecond * 1, "", nil
}

func (fake FakeClient) Write(bp client.BatchPoints) error {
	return nil
}

func (fake FakeClient) Query(q client.Query) (*client.Response, error) {
	return &client.Response{
		[]client.Result{},
		nil,
	}, nil
}

func (fake FakeClient) Close() error {
	return nil
}

func TestAddMeasurementToBatch(t *testing.T) {
	m := Measurement{
		DeviceId:    1,
		Temperature: 10,
		Current:     21,
		Voltage:     33,
		Power:       44,
		Gpio:        10,
	}
	batch, err := client.NewBatchPoints(client.HTTPConfig{
		Addr:     influx.Addr,
		Username: influx.UserName,
		Password: influx.Password,
	})
	assert.NoError(t, err)
	influx.AddMeasurementToBatch(m, batch)
	assert.Equal(t, len(batch.Points()), 1)
	assert.Equal(t, batch.Database(), influx.DbName)
}

func TestQueryDB(t *testing.T) {
	measurements, err := influx.queryDB("My command")
	assert.NoError(t, err)
	assert.NotNil(t, measurements)
	assert.Equal(t, len(measurements), 0)
}

func TestInfluxGetMeasurement(t *testing.T) {
	measurements, err := influx.GetMeasurement(10)
	assert.NoError(t, err)
	assert.NotNil(t, measurements)
	assert.Equal(t, len(measurements), 0)
}

func TestInfluxGetDeviceById(t *testing.T) {
	device, err := influx.GetDeviceById("abcd")
	assert.Nil(t, device)
	assert.Error(t, err)
}

func TestInfluxCreateDevice(t *testing.T) {
	err := influx.CreateDevice("123", "", 1)
	assert.Error(t, err)
}
