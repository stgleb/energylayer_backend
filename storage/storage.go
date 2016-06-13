package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// Storage interface for storing and retrieving data.
type Storage interface {
	CreateMeasurement(m Measurement) error
	GetMeasurement(count int) ([]Measurement, error)
	GetDeviceById(uuid string) (Device, error)
	CreateDevice(uuid, ipAddress string, userId int)
}

const (
	USER     = "root"
	PASSWORD = "1234"
	DATABASE = "energylayer"
	PORT     = "3306"
	HOST     = "localhost"
)

type DatabaseStorage struct {
	Database sql.DB
}

func StorageFactory(storageType string) (*sql.DB, error) {
	var uri string
	var dbType string

	switch storageType {
	case "mysql":
		uri = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASSWORD, HOST, PORT, DATABASE)
		dbType = "mysql"
	case "memory":
		uri = ":memory:"
		dbType = "sqlite3"
	}

	storage, err := NewStorage(uri, dbType)

	if err != nil {
		return storage, err
	}

	return storage, nil
}

func NewStorage(uri, dbType string) (*sql.DB, error) {
	db, err := sql.Open(dbType, uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return db, err
	}

	return db, nil
}

func CreateMeasurement(Database *sql.DB, m Measurement) error {
	tx, err := Database.Begin()
	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT measurement(voltage, power, temperature, device_id)" +
		"VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	stmt.Exec(m.Voltage, m.Power, m.Temperature, m.DeviceId)
	err = tx.Commit()

	if err != nil {
		log.Printf("Error while commiting transaction message: %s", err.Error())
		return err
	}

	log.Printf("Data has been inserted sucessfully")
	return nil
}

func GetMeasurements(Database *sql.DB, count int) ([]Measurement, error) {
	var result []Measurement
	result = make([]Measurement, count)

	tx, err := Database.Begin()
	defer tx.Rollback()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return nil, err
	}

	rows, err := tx.Query("SELECT timestamp, voltage, power, temperature, device_id"+
		" from measurement ORDER BY DESC LIMIT ?", count)
	defer rows.Close()

	if err != nil {
		log.Printf("Error while executing query %s", err.Error())
		return nil, err
	}

	var timestamp int64
	var voltage int
	var power int
	var temperature int
	var device_id int

	for rows.Next() {
		err := rows.Scan(&timestamp, &voltage, &power, &temperature, &device_id)

		if err != nil {
			log.Printf("Error while extracting value from row %s", err.Error())
			return nil, err
		}

		m := Measurement{
			Timestamp:   timestamp,
			Voltage:     voltage,
			Power:       power,
			Temperature: temperature,
			DeviceId:    device_id,
		}

		result = append(result, m)
	}

	return result, nil
}

func GetDeviceById(Database *sql.DB, uuid string) (Device, error) {
	tx, err := Database.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return Device{}, err
	}
	defer tx.Rollback()

	var ipAddr string
	var id int

	row := tx.QueryRow("select id, uuid, ip_addr from device where uuid = ?", 1)
	err = row.Scan(&id, &uuid, &ipAddr)

	if err != nil {
		log.Printf("Error while reading data from the row %s", err.Error())
		return Device{}, err
	}

	return Device{
		Id:        id,
		Uuid:      uuid,
		IpAddress: ipAddr,
	}, nil
}

func UpdateDeviceIP(Database *sql.DB, uuid, ipAddr string) error {
	tx, err := Database.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}
	defer tx.Rollback()
	log.Printf("Update device uuid = %s ip address to %s", uuid, ipAddr)

	_, err = tx.Exec("update device set ip_addr = ? where uuid = ?", ipAddr, uuid)

	if err != nil {
		log.Printf("Error while updating device %s ip address to %s %s", uuid, ipAddr, err.Error())
		return err
	}

	tx.Commit()
	return nil
}

func CreateDevice(Database *sql.DB, uuid, ipAddress string) (int, error) {
	tx, err := Database.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return -1, err
	}
	defer tx.Rollback()

	log.Printf("Inserting new device with uuid %s", uuid)
	_, err = tx.Exec("INSERT device(uuid, ip_addr) VALUES(?, ?)", uuid, ipAddress)

	if err != nil {
		log.Printf("Error while inserting device: %s", err.Error())
		return -1, err
	}

	var id int

	row := tx.QueryRow("select id from device where uuid = ?", 1)
	err = row.Scan(&id)

	if err != nil {
		log.Printf("Error while reading data from the row %s", err.Error())
		return -1, err
	}

	tx.Commit()
	return id, nil
}
