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

func StorageFactory(storageType string) Storage {
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
		return nil
	}

	return storage
}

func NewStorage(uri, dbType string) (Storage, error) {
	db, err := sql.Open(dbType, uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return nil, err
	}

	return Storage{
		db,
	}, nil
}

func (self *DatabaseStorage) CreateMeasurement(m Measurement) {
	tx, err := self.Database.Begin()
	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT measurement(voltage, power, temperature, device_id)" +
		"VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	stmt.Exec(m.Voltage, m.Power, m.Temperature, m.DeviceId)
	err = tx.Commit()

	if err != nil {
		log.Printf("Error while commiting transaction message: %s", err.Error())
	}
}

func (self *DatabaseStorage) GetMeasurements(count int) ([]Measurement, error) {
	var result []Measurement
	result = make([]Measurement, count)

	tx, err := self.Database.Begin()
	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return nil, err
	}

	rows, err := tx.Query("SELECT voltage, power, temperature, device_id"+
		"from measurement ORDER BY DESC LIMIT ?", count)
	defer rows.Close()

	if err != nil {
		log.Printf("Error while executing query %s", err.Error())
		return nil, err
	}

	var voltage int
	var power int
	var temperature int
	var device_id int

	for rows.Next() {
		err := rows.Scan(&voltage, &power, &temperature, &device_id)

		if err != nil {
			log.Printf("Error while extracting value from row %s", err.Error())
			return nil, err
		}

		m := Measurement{
			Voltage:     voltage,
			Power:       power,
			Temperature: temperature,
			DeviceId:    device_id,
		}

		result = append(result, m)
	}

	return result, nil
}

func (self *DatabaseStorage) GetDeviceById(uuid string) (Device, error) {
	tx, err := self.Database.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return nil, err
	}
	defer tx.Rollback()

	var uuid string
	var ipAddr string
	var id int
	var user_id int

	row := tx.QueryRow("select id, uuid, user_id, ip_addr from users where uuid = ?", 1)
	err = row.Scan(&id, &uuid, &user_id, &ipAddr)

	if err != nil {
		log.Printf("Error while reading data from the row %s", err.Error())
		return nil, err
	}

	return Device{
		Id:        id,
		Uuid:      uuid,
		UserId:    user_id,
		IpAddress: ipAddr,
	}, nil
}

func (self *DatabaseStorage) CreateDevice(uuid, ipAddress string) error {
	tx, err := self.Database.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT device(uuid, ip_addr) VALUES(?, ?)", uuid, ipAddress)

	if err != nil {
		log.Printf("Error while inserting device: %s", err.Error())
		return err
	}

	tx.Commit()
	return nil
}
