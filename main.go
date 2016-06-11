package main

import (
	"encoding/hex"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"fmt"
	"flag"
)

const (
	USER = "root"
	PASSWORD = "1234"
	DATABASE = "energylayer"
	PORT = "3306"
	HOST = "localhost"
)

type Measurement struct {
	Timestamp   int64
	Gpio        int
	Voltage     int
	Power       int
	Temperature int
	DeviceId   int
}

type Device struct {
	Id int
	Uuid string
	UserId int
	IpAddress string
}

// Storage interface for storing and retrieving data.
type Storage interface {
	SaveMeasurement(m Measurement) error
	LoadMeasurements(count int) ([]Measurement, error)
	GetDeviceById(uuid string) (Device, error)
	CreateDevice(uuid, ipAddress string , userId int)
}

// MySQL storage implementation.
type MySQLStorage struct {
	Database sql.DB
}

func NewMySQLStorage(dbName ,user, password string) (Storage, error) {
	host := "localhost"
	port := "3306"
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db, err  := sql.Open("mysql", uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return nil, err
	}

	return MySQLStorage{
		db,
	}, nil
}

func StorageFactory(storageType string) Storage {
	switch storageType {
	case "mysql":
		storage, err := NewMySQLStorage(storageType, USER, PASSWORD)

		if err != nil {
			return nil
		}

		return storage
	}

	return nil
}

func (self *MySQLStorage) Save(m Measurement) {
	tx ,err := self.Database.Begin()
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
	tx ,err := self.Database.Begin()

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
		log.Printf("Error while reading data from the row %s",err.Error())
		return nil, err
	}

	return Device{
		Id:id,
		Uuid: uuid,
		UserId: user_id,
		IpAddress:ipAddr,
	}, nil
}

func (self *MySQLStorage) CreateDevice(uuid, ipAddress string) error{
	tx ,err := self.Database.Begin()

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

func (m Measurement) String() string {
	return fmt.Sprintf("{timestamp: %d, gpio: %d, voltage: %d, power: %d, temperature: %d}",
		m.Timestamp, m.Gpio, m.Voltage, m.Power, m.Temperature)
}

func DecodeData(data string) Measurement {
	timestamp := time.Now().Unix()
	tmp, _ := hex.DecodeString(data[:4])
	gpio := int(tmp[0])
	tmp, _ = hex.DecodeString(data[4:8])
	voltage := int(tmp[0])
	tmp, _ = hex.DecodeString(data[8:12])
	power := int(tmp[0])
	tmp, _ = hex.DecodeString(data[12:16])
	temperature := int(tmp[0])

	return Measurement{
		Timestamp: timestamp,
		Gpio: gpio,
		Voltage: voltage,
		Power: power,
		Temperature: temperature,
	}
}

func StoreData(m Measurement) {

}

func Receiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device_id := vars["device_id"]
	data := vars["data_string"]
	log.Printf("Received data %s from device %s", data, device_id)

	m := DecodeData(data)
	log.Printf("Measurement %v", m)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rs/data/post/{device_id}/{data_string}", Receiver)

	port := *flag.String("port", "8000", "Port to listen")
	host := *flag.String("host", "0.0.0.0", "host name")
	addressString := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listen addr %s", addressString)

	log.Fatal(http.ListenAndServe(addressString, r))
}
