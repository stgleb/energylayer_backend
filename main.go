package main

import (
	"./storage"
	"encoding/hex"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)




func DecodeData(data string) storage.Measurement {
	timestamp := time.Now().Unix()
	tmp, _ := hex.DecodeString(data[:4])
	gpio := int(tmp[0])
	tmp, _ = hex.DecodeString(data[4:8])
	voltage := int(tmp[0])
	tmp, _ = hex.DecodeString(data[8:12])
	power := int(tmp[0])
	tmp, _ = hex.DecodeString(data[12:16])
	temperature := int(tmp[0])

	return storage.Measurement{
		Timestamp:   timestamp,
		Gpio:        gpio,
		Voltage:     voltage,
		Power:       power,
		Temperature: temperature,
	}
}

func StoreData(m storage.Measurement) {

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
