package main

import (
	"./storage"
	"./utils"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

var (
	db  *storage.DatabaseStorage
)

func init() {
	log.Println("Initializing application")
	db, err := storage.StorageFactory("mysql")

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
	}

	log.Println("Connected to database, verify connection")
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
	device, err := db.GetDeviceById(device_id)

	if device.IpAddress != ipAddr {
		db.UpdateDeviceIP(device_id, ipAddr)
	}
	var id int

	if err != nil {
		log.Printf("Device not found %s . Try to create new one", err.Error())
		err = db.CreateDevice(device_id, ipAddr)

		if err != nil {
			log.Printf("Error while creating new device %s", err.Error())
		}

		device, err = db.GetDeviceById(ipAddr)
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
	err = db.CreateMeasurement(m)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(m)
	w.WriteHeader(http.StatusCreated)
}

//
//func UserDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	if len(vars["count"]) > 0 {
//		count := int(vars["count"])
//		log.Printf("Count: %d", count)
//	}
//
//	if interval, ok := vars["interval"]; ok {
//		log.Printf("Interval: %d", interval)
//	}
//}
//
//func DeviceMeasurements(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	device_id := vars["device_id"]
//	if len(vars["count"]) > 0 {
//		count := int(vars["count"])
//		log.Printf("Count: %d", count)
//	}
//
//	log.Printf("device_id: %s", device_id)
//
//	if interval, ok := vars["interval"]; ok {
//		log.Printf("Interval: %d", interval)
//	}
//}

func main() {
	//m := &storage.Measurement{
	//	Gpio: 1,
	//	Temperature: 1,
	//	Power: 1,
	//	Current: 1,
	//	Voltage: 1,
	//	DeviceId: 1,
	//	Timestamp: 11335,
	//}
	//
	//kafka, err := storage.NewKafkaStorage([]string{"localhost:9092"})
	//
	//if err != nil {
	//	log.Printf("Error while connecting to Kafka")
	//}
	//err = kafka.CreateMeasurement(m, "TutorialTopic")
	//
	//if err != nil {
	//	log.Printf("Error while creating measurement %s", err.Error())
	//} else {
	//	log.Printf("Success!!!")
	//}

	r := mux.NewRouter()
	r.HandleFunc("/rs/data/post/{device_id}/{data_string}", Receiver)
	//r.HandleFunc("/api/data/user/measurement/{count}/{interval}", UserDeviceMeasurements)
	//r.HandleFunc("/api/data/user/measurement/{count}", UserDeviceMeasurements)
	//r.HandleFunc("/api/measurement/{device_uuid}/count/{count}/{interval}", DeviceMeasurements)
	//r.HandleFunc("/api/measurement/{device_uuid}/count/{count}", DeviceMeasurements)

	port := *flag.String("port", "8000", "Port to listen")
	host := *flag.String("host", "0.0.0.0", "host name")
	addressString := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listen addr %s", addressString)
	log.Fatal(http.ListenAndServe(addressString, r))
}
