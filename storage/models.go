package storage

import "fmt"

type Measurement struct {
	Timestamp   int64
	Gpio        int
	Voltage     int
	Power       int
	Temperature int
	DeviceId    int
}

type Device struct {
	Id        int
	Uuid      string
	UserId    int
	IpAddress string
}

func (m Measurement) String() string {
	return fmt.Sprintf("{timestamp: %d, gpio: %d, voltage: %d, power: %d, temperature: %d}",
		m.Timestamp, m.Gpio, m.Voltage, m.Power, m.Temperature)
}
