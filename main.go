package main

import (
	"./storage"
	"./utils"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

var db *sql.DB
var err error

func init() {
	log.Println("Initializing application")
	db, err = storage.StorageFactory("mysql")

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
	}

	log.Printf("Connected to database, verify connection")
	err = db.Ping()

	if err != nil {
		log.Fatalf("Error pinging to database %s", err.Error())
	}

	log.Println("Connection has been established successfully")
}

func Receiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device_id := vars["device_id"]
	data := vars["data_string"]
	ipAddr := strings.Split(r.RemoteAddr, ":")[0]
	log.Printf("Received data %s from device %s", data, device_id)

	log.Printf("Get device by uuid = %s", device_id)
	device, err := storage.GetDeviceById(db, ipAddr)

	if device.IpAddress != ipAddr {
		storage.UpdateDeviceIP(db, device_id, ipAddr)
	}

	var id int

	if err != nil {
		log.Printf("Device not found %s . Try to create new one", err.Error())
		id, err = storage.CreateDevice(db, device_id, ipAddr)

		if err != nil {
			log.Printf("Error while creating new device %s", err.Error())
		}
	} else {
		id = device.Id
	}

	timestamp, gpio, voltage, power, temperature := utils.DecodeData(data)
	m := storage.Measurement{
		DeviceId:    id,
		Timestamp:   timestamp,
		Gpio:        gpio,
		Voltage:     voltage,
		Power:       power,
		Temperature: temperature,
	}

	// Save measurement to database .
	err = storage.CreateMeasurement(db, m)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(m.String())
	w.WriteHeader(http.StatusCreated)
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
