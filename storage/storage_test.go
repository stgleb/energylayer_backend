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
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.Nil(t, err)
}

func TestCreateMeasurement(t *testing.T) {
	db, _ := StorageFactory("memory")

	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     10,
		Power:       20,
		Temperature: 30,
	}

	err := CreateMeasurement(db, m)
	assert.Nil(t, err)
}

func TestGetMeasurements(t *testing.T) {
	db, err := StorageFactory("memory")
	assert.Nil(t, err)
	measurements, _ := GetMeasurements(db, 4)

	log.Printf("%v", db)
	// Test that data is empty
	assert.Equal(t, len(measurements), 0)

	// Test data retrieved.
	m := Measurement{
		Timestamp:   time.Now().Unix(),
		Voltage:     10,
		Power:       20,
		Temperature: 30,
	}

	_ = CreateMeasurement(db, m)
	_ = CreateMeasurement(db, m)
	_ = CreateMeasurement(db, m)
	_ = CreateMeasurement(db, m)

	measurements, _ = GetMeasurements(db, 2)
	assert.Equal(t, len(measurements), 2)

}

func TestCreateDevice(t *testing.T) {
	db, _ := StorageFactory("memory")
	id, err := CreateDevice(db, "abcd", "127.0.0.1")
	assert.Nil(t, err)
	assert.NotNil(t, id)

	id, err = CreateDevice(db, "abcd", "127.0.0.1")
	assert.Equal(t, id, -1)
	assert.NotNil(t, err)
}

func TestGetDeviceById(t *testing.T) {
	db, _ := StorageFactory("memory")
	uuid := "abcd"
	id, err := CreateDevice(db, uuid, "127.0.0.1")
	assert.Nil(t, err)
	assert.NotNil(t, id)

	device, err := GetDeviceById(db, uuid)
	assert.Equal(t, device.Uuid, uuid)
	assert.Nil(t, err)
}
