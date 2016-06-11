package storage

import (
	"log"
	"database/sql"
	"fmt"
)


// MySQL storage implementation.
type MySQLStorage struct {
	Database sql.DB
}

func NewMySQLStorage(dbName, user, password, host, port string) (Storage, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db, err := sql.Open("mysql", uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return nil, err
	}

	return MySQLStorage{
		db,
	}, nil
}


func (self *MySQLStorage) Save(m Measurement) {
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

func (self *MySQLStorage) GetDeviceById(uuid string) (Device, error) {
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

func (self *MySQLStorage) CreateDevice(uuid, ipAddress string) error {
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

