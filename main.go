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

type Measurement struct {
	Timestamp   int64
	Gpio        int
	Voltage     int
	Power       int
	Temperature int
}

// Storage interface for storing and retrieving data.
type Storage interface {
	Save(m Measurement) error
	Load(count int) ([]Measurement, error)
}

// MySQL storage implementation.
type MySQLStorage struct {
	db sql.DB
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
