package influxdb_storage

import (
	. "../../storage"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"os"
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

var (
	fakeClient FakeClient
	influx     InfluxDbStorage
)

func (fake FakeClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Millisecond * 1, "", nil
}

func (fake FakeClient) Write(bp client.BatchPoints) error {
	return nil
}

func (fake FakeClient) Query(q client.Query) (*client.Response, error) {
	return &client.Response{
		[]client.Result{},
		"",
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
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.DbName,
		Precision: PRECISION,
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
	measurements, err := influx.GetMeasurements(10)
	assert.NoError(t, err)
	assert.NotNil(t, measurements)
	assert.Equal(t, len(measurements), 0)
}

func TestInfluxGetDeviceById(t *testing.T) {
	device, err := influx.GetDeviceById("abcd")
	assert.Equal(t, 0, device.Id)
	assert.Error(t, err)
}

func TestInfluxCreateDevice(t *testing.T) {
	err := influx.CreateDevice("123", "", 1)
	assert.Error(t, err)
}

func TestMain(m *testing.M) {
	// Initializing influx object
	influx = InfluxDbStorage{
		DbName:   DB_NAME,
		UserName: USER,
		Password: PASSWORD,
		Addr:     ADDR,
		Client:   fakeClient,
	}
	code := m.Run()
	os.Exit(code)
}
