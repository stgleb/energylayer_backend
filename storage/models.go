package storage

import "fmt"

type Measurement struct {
	Timestamp   int64 `json:"timestamp"`
	Gpio        int   `json:"gpio"`
	Voltage     int   `json:"voltage"`
	Current     int   `json:"current"`
	Power       int   `json:"power"`
	Temperature int   `json:"temperature"`
	DeviceId    int   `json:"device_id"`

	encoded []byte
	err     error
}

type Device struct {
	Id        int
	Uuid      string
	UserId    int
	IpAddress string
}

func (m Measurement) String() string {
	return fmt.Sprintf("<Measurement {timestamp: %v, gpio: %v, voltage: %v, power: %v, temperature: %v}>",
		m.Timestamp, m.Gpio, m.Voltage, m.Power, m.Temperature)
}

func (d Device) String() string {
	return fmt.Sprintf("<Device :{ id: %d, uuid: %s, user_id: %d, ip_address: %s}>",
		d.Id, d.Uuid, d.UserId, d.IpAddress)
}
