package storage

// Storage interface for storing and retrieving data.
type Storage interface {
	CreateMeasurements(m []Measurement) error
	GetMeasurementsByDevice(int, int) ([]Measurement, error)
	GetMeasurements(count int) ([]Measurement, error)
	GetDeviceById(uuid string) (Device, error)
	CreateDevice(uuid, ipAddress string, userId int) error
}
