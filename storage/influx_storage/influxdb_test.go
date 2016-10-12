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

var (
	mockClient MockClient
	influx     InfluxDbStorage
)

func TestNewInfluxDBStorageError(t *testing.T) {
	storage, err := NewInfluxDBStorage(USER, PASSWORD, DB_NAME, ADDR, "_http")
	assert.Error(t, err)
	assert.Nil(t, storage)

}

func TestAddMeasurementToBatch(t *testing.T) {
	mockClient.Refresh()
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
	mockClient.Refresh()
	measurements, err := influx.queryDB("My command")
	assert.NoError(t, err)
	assert.NotNil(t, measurements)
	assert.Equal(t, len(measurements), 0)
	assert.Equal(t, 1, mockClient.Called["Query"])
}

func TestInfluxGetMeasurement(t *testing.T) {
	mockClient.Refresh()
	measurements, err := influx.GetMeasurements(10)
	assert.NoError(t, err)
	assert.NotNil(t, measurements)
	assert.Equal(t, len(measurements), 0)
	assert.Equal(t, 1, mockClient.Called["Query"])
}

func TestInfluxCreateMeasurements(t *testing.T) {
	mockClient.Refresh()
	measurements := []Measurement{Measurement{
		Voltage:     0,
		Power:       0,
		Temperature: 0,
		Current:     0,
	}, Measurement{
		Voltage:     0,
		Power:       0,
		Temperature: 0,
		Current:     0,
	}}
	err := influx.CreateMeasurements(measurements)
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.Called["Write"])
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
	mockClient.Refresh()
	// Mocking client functions
	mockClient.FakeClose = func() error {
		mockClient.Called["Close"] += 1
		return nil
	}

	mockClient.FakeWrite = func(bp client.BatchPoints) error {
		mockClient.Called["Write"] += 1
		return nil
	}

	mockClient.FakeQuery = func(q client.Query) (*client.Response, error) {
		mockClient.Called["Query"] += 1
		return &client.Response{
			[]client.Result{},
			"",
		}, nil
	}

	mockClient.FakePing = func(timeout time.Duration) (time.Duration, string, error) {
		mockClient.Called["Ping"] += 1
		return time.Millisecond * 1, "", nil
	}

	// Initializing influx object
	influx = InfluxDbStorage{
		DbName:   DB_NAME,
		UserName: USER,
		Password: PASSWORD,
		Addr:     ADDR,
		Client:   mockClient,
	}
	code := m.Run()
	os.Exit(code)
}
