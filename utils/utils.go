package utils

import (
	"encoding/hex"
	"log"
	"time"
)

func convertByteToInt(in []byte) int {
	return (int(in[0])<<24 | int(in[1]))
}

func DecodeData(data string) (int64, int, int, int, int) {
	timestamp := time.Now().Unix()
	tmp, _ := hex.DecodeString(data[:4])
	gpio := convertByteToInt(tmp)
	log.Printf("gpio %d", gpio)

	tmp, _ = hex.DecodeString(data[4:8])
	voltage := convertByteToInt(tmp)
	log.Printf("voltage %d", voltage)

	tmp, _ = hex.DecodeString(data[8:12])
	power := convertByteToInt(tmp)
	log.Printf("power %d", power)

	tmp, _ = hex.DecodeString(data[12:16])
	temperature := convertByteToInt(tmp)
	log.Printf("temperature %d", temperature)

	log.Printf("timestamp %d gpio voltage %d power %d  temperature %d",
		timestamp, gpio, power, temperature)

	return timestamp, gpio, voltage, power, temperature
}
