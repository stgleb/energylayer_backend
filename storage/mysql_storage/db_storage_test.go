package storage

import (
	. "../../storage"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"
)

func newStorage(uri, dbType string) (*sql.DB, error) {
	db, err := sql.Open(dbType, uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return db, err
	}

	return db, nil
}

func getDB() (DatabaseStorage, error) {
	var uri string = ":memory:"
	var dbType string = "sqlite3"

	storage, err := newStorage(uri, dbType)

	if err != nil {
		return DatabaseStorage{storage}, err
	}

	return DatabaseStorage{storage}, nil
}

func createMeasurement() Measurement {
	// Test data retrieved.
	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     rand.Int(),
		Power:       rand.Int(),
		Temperature: rand.Int(),
	}

	return m
}

func TestCreateMeasurement(t *testing.T) {
	db, err := getDB()
	assert.NoError(t, err)

	// Create table measurement
	_, err = db.Exec("CREATE TABLE measurement( `id` INTEGER PRIMARY KEY,`tag` varchar(64) DEFAULT NULL,`gpio` int(11) DEFAULT NULL,`voltage` int(11) DEFAULT NULL,`power` int(11) DEFAULT NULL,`temperature` int(11) DEFAULT NULL,`timestamp` int(11) DEFAULT NULL,`device_id` int(11) DEFAULT NULL);")
	assert.NoError(t, err)

	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     10,
		Power:       20,
		Temperature: 30,
	}

	err = db.CreateMeasurement(m)
	assert.NoError(t, err)
}

func TestGetMeasurements(t *testing.T) {
	db, err := getDB()
	assert.NoError(t, err)
	// Create table measurement
	_, err = db.Exec("CREATE TABLE measurement( `id` INTEGER PRIMARY KEY,`tag` varchar(64) DEFAULT NULL,`gpio` int(11) DEFAULT NULL,`voltage` int(11) DEFAULT NULL,`power` int(11) DEFAULT NULL,`temperature` int(11) DEFAULT NULL,`timestamp` int(11) DEFAULT NULL,`device_id` int(11) DEFAULT NULL);")
	measurements, _ := db.GetMeasurements(4)

	log.Printf("%v", db)
	// Test that data is empty
	assert.Equal(t, 0, len(measurements))

	m1 := createMeasurement()
	err = db.CreateMeasurement(m1)
	m2 := createMeasurement()
	err = db.CreateMeasurement(m2)
	m3 := createMeasurement()
	err = db.CreateMeasurement(m3)
	m4 := createMeasurement()
	err = db.CreateMeasurement(m4)
	assert.NoError(t, err)

	measurements, err = db.GetMeasurements(4)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(measurements))
}

func TestGetMeasurementsByDevice(t *testing.T) {
	db, err := getDB()
	assert.NoError(t, err)
	// Create table measurement
	_, err = db.Exec("CREATE TABLE measurement( `id` INTEGER PRIMARY KEY,`tag` varchar(64) DEFAULT NULL,`gpio` int(11) DEFAULT NULL,`voltage` int(11) DEFAULT NULL,`power` int(11) DEFAULT NULL,`temperature` int(11) DEFAULT NULL,`timestamp` int(11) DEFAULT NULL,`device_id` int(11) DEFAULT NULL);")
	measurements, _ := db.GetMeasurements(4)

	log.Printf("%v", db)
	// Test that data is empty
	assert.Equal(t, 0, len(measurements))

	m1 := createMeasurement()
	m1.DeviceId = 2
	err = db.CreateMeasurement(m1)
	m2 := createMeasurement()
	m2.DeviceId = 2
	err = db.CreateMeasurement(m2)
	m3 := createMeasurement()
	m3.DeviceId = 3
	err = db.CreateMeasurement(m3)
	m4 := createMeasurement()
	m4.DeviceId = 2
	err = db.CreateMeasurement(m4)
	assert.NoError(t, err)

	measurements, err = db.GetMeasurementsByDevice(2, 4)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(measurements))
}

func TestCreateDevice(t *testing.T) {
	db, err := getDB()
	assert.NoError(t, err)
	// Create table devices
	_, err = db.Exec("CREATE TABLE device(`id` INTEGER PRIMARY KEY,`uuid` varchar(64) UNIQUE NOT NULL,`user_id` int(11) DEFAULT NULL,`ip_addr` varchar(40) DEFAULT NULL);")
	assert.NoError(t, err)
	// Insert first device
	err = db.CreateDevice("abcd", "127.0.0.1", 1)
	assert.NoError(t, err)
	// Check that device with the same uuid will fail
	err = db.CreateDevice("abcd", "127.0.0.1", 1)
	assert.Error(t, err)
}

func TestGetDeviceById(t *testing.T) {
	db, err := getDB()
	assert.NoError(t, err)
	// Create table devices
	_, err = db.Exec("CREATE TABLE device(`id` INTEGER PRIMARY KEY,`uuid` varchar(64) UNIQUE NOT NULL,`user_id` int(11) DEFAULT NULL,`ip_addr` varchar(40) DEFAULT NULL);")
	assert.NoError(t, err)

	uuid := "abcd"
	err = db.CreateDevice(uuid, "127.0.0.1", 1)
	assert.NoError(t, err)

	device, err := db.GetDeviceById(uuid)
	assert.Equal(t, uuid, device.Uuid)
	assert.NoError(t, err)
}
