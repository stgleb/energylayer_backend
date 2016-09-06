package storage

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestStorageFactory(t *testing.T) {
	_, err := StorageFactory("wrong")
	assert.NotNil(t, err)

	db, err := StorageFactory("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.NoError(t, err)
}

func TestCreateMeasurement(t *testing.T) {
	db, _ := StorageFactory("memory")
	// Create table measurement
	_, err := db.Exec("CREATE TABLE measurement( `id` int(11) NOT NULL,`tag` varchar(64) DEFAULT NULL,`gpio` int(11) DEFAULT NULL,`voltage` int(11) DEFAULT NULL,`power` int(11) DEFAULT NULL,`temperature` int(11) DEFAULT NULL,`timestamp` int(11) DEFAULT NULL,`device_id` int(11) DEFAULT NULL,PRIMARY KEY (`id`));")
	assert.NoError(t, err)

	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     10,
		Power:       20,
		Temperature: 30,
	}

	err = CreateMeasurement(db, m)
	assert.NoError(t, err)
}

func TestGetMeasurements(t *testing.T) {
	db, err := StorageFactory("memory")
	assert.NoError(t, err)
	// Create table measurement
	_, err = db.Exec("CREATE TABLE measurement( `id` int(11) NOT NULL,`tag` varchar(64) DEFAULT NULL,`gpio` int(11) DEFAULT NULL,`voltage` int(11) DEFAULT NULL,`power` int(11) DEFAULT NULL,`temperature` int(11) DEFAULT NULL,`timestamp` int(11) DEFAULT NULL,`device_id` int(11) DEFAULT NULL,PRIMARY KEY (`id`));")
	measurements, _ := GetMeasurements(db, 4)

	log.Printf("%v", db)
	// Test that data is empty
	assert.Equal(t, 0, len(measurements))

	// Test data retrieved.
	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     10,
		Power:       20,
		Temperature: 30,
	}

	err = CreateMeasurement(db, m)
	m.Timestamp += 1
	err = CreateMeasurement(db, m)
	m.Timestamp += 1
	err = CreateMeasurement(db, m)
	m.Timestamp += 1
	err = CreateMeasurement(db, m)
	assert.NoError(t, err)

	measurements, _ = GetMeasurements(db, 2)
	// FIXME: measurements aren't stored, so len(measurements) == 0
	assert.Equal(t, 2, len(measurements))
}

func TestCreateDevice(t *testing.T) {
	db, err := StorageFactory("memory")
	assert.NoError(t, err)
	// Create table devices
	_, err = db.Exec("CREATE TABLE device(`id` INTEGER PRIMARY KEY,`uuid` varchar(64) UNIQUE NOT NULL,`user_id` int(11) DEFAULT NULL,`ip_addr` varchar(40) DEFAULT NULL);")
	assert.NoError(t, err)
	// Insert first device
	err = CreateDevice(db, "abcd", "127.0.0.1")
	assert.NoError(t, err)
	// Check that device with the same uuid will fail
	err = CreateDevice(db, "abcd", "127.0.0.1")
	assert.Error(t, err)
}

func TestGetDeviceById(t *testing.T) {
	db, err := StorageFactory("memory")
	assert.NoError(t, err)
	// Create table devices
	_, err = db.Exec("CREATE TABLE device(`id` INTEGER PRIMARY KEY,`uuid` varchar(64) UNIQUE NOT NULL,`user_id` int(11) DEFAULT NULL,`ip_addr` varchar(40) DEFAULT NULL);")
	assert.NoError(t, err)

	uuid := "abcd"
	err = CreateDevice(db, uuid, "127.0.0.1")
	assert.NoError(t, err)

	device, err := GetDeviceById(db, uuid)
	assert.Equal(t, uuid, device.Uuid)
	assert.NoError(t, err)
}
