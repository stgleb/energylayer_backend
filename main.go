package main

import (
	"./storage"
	"./utils"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"database/sql"
)

var db *sql.DB
var err error

func init() {
	log.Println("Initializing application")
	db, err =  storage.StorageFactory("mysql")

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
	log.Printf("Received data %s from device %s", data, device_id)

	timestamp, gpio, voltage, power, temperature := utils.DecodeData(data)
	m :=  storage.Measurement{
		Timestamp: timestamp,
		Gpio: gpio,
		Voltage: voltage,
		Power: power,
		Temperature: temperature,
	}

	log.Println(m.String())
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
