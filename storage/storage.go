package storage

// Storage interface for storing and retrieving data.
type Storage interface {
	CreateMeasurement(m Measurement) error
	GetMeasurement(count int) ([]Measurement, error)
	GetDeviceById(uuid string) (Device, error)
	CreateDevice(uuid, ipAddress string, userId int) error
}
