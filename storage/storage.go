package storage

// Storage interface for storing and retrieving data.
type Storage interface {
	SaveMeasurement(m Measurement) error
	LoadMeasurements(count int) ([]Measurement, error)
	GetDeviceById(uuid string) (Device, error)
	CreateDevice(uuid, ipAddress string, userId int)
}

const (
	USER     = "root"
	PASSWORD = "1234"
	DATABASE = "energylayer"
	PORT     = "3306"
	HOST     = "localhost"
)

func StorageFactory(storageType string) Storage {
	switch storageType {
	case "mysql":
		storage, err := NewMySQLStorage(DATABASE, USER, PASSWORD, HOST, PORT)

		if err != nil {
			return nil
		}

		return storage
	case "sqlite:memory":
		// TODO: add in memory sqlite.
	}

	return nil
}