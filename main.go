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

var (
	db  *sql.DB
	err error
)

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
		err = storage.CreateDevice(db, device_id, ipAddr)

		if err != nil {
			log.Printf("Error while creating new device %s", err.Error())
		}

		device, err = storage.GetDeviceById(db, ipAddr)
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

	log.Println(m)
	w.WriteHeader(http.StatusCreated)
}

func UserDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := int(vars["count"])

	log.Printf("Count: %d", count)
	if interval, ok := vars["interval"]; ok {
		log.Printf("Interval: %d", interval)
	}
}

func DeviceMeasurements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device_id := vars["device_id"]
	count := int(vars["count"])

	log.Printf("device_id: %s", device_id)
	log.Printf("Count: %d", count)

	if interval, ok := vars["interval"]; ok {
		log.Printf("Interval: %d", interval)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rs/data/post/{device_id}/{data_string}", Receiver)
	r.HandleFunc("/api/data/user/measurement/{count}/{interval}", UserDeviceMeasurements)
	r.HandleFunc("/api/data/user/measurement/{count}", UserDeviceMeasurements)
	r.HandleFunc("/api/measurement/{device_uuid}/count/{count}/{interval}", DeviceMeasurements)
	r.HandleFunc("/api/measurement/{device_uuid}/count/{count}", DeviceMeasurements)

	port := *flag.String("port", "8000", "Port to listen")
	host := *flag.String("host", "0.0.0.0", "host name")
	addressString := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listen addr %s", addressString)

	log.Fatal(http.ListenAndServe(addressString, r))
}
