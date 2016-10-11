package storage

import (
	. "../../storage"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DatabaseStorage struct {
	*sql.DB
}

func (db DatabaseStorage) CreateMeasurement(m Measurement) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO measurement(timestamp, voltage, power, temperature, device_id)" +
		"VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error while creating statement %s", err.Error())
		return err
	}

	_, err = stmt.Exec(m.Timestamp, m.Voltage, m.Power, m.Temperature, m.DeviceId)
	if err != nil {
		log.Printf("Error while execuring statement: %s", err.Error())
		return err
	}

	// Prepared statements should be closed before commit transaction
	stmt.Close()
	err = tx.Commit()

	if err != nil {
		log.Printf("Error while commiting transaction message: %s", err.Error())
		return err
	}

	log.Printf("Data has been inserted sucessfully")
	return nil
}

// TODO: remove code duplication in both methods
func (db DatabaseStorage) CreateMeasurements(measurements []Measurement) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO measurement(timestamp, voltage, power, temperature, device_id)" +
		"VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error while creating statement %s", err.Error())
		return err
	}

	for _, m := range measurements {
		_, err = stmt.Exec(m.Timestamp, m.Voltage, m.Power, m.Temperature, m.DeviceId)
		if err != nil {
			log.Printf("Error while execuring statement: %s", err.Error())
			return err
		}
	}
	// Prepared statements should be closed before commit transaction
	stmt.Close()
	err = tx.Commit()

	if err != nil {
		log.Printf("Error while commiting transaction message: %s", err.Error())
		return err
	}

	log.Printf("Data has been inserted sucessfully")
	return nil
}

func (db DatabaseStorage) GetMeasurements(count int) ([]Measurement, error) {
	result := make([]Measurement, 0, count)
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return nil, err
	}

	rows, err := tx.Query("SELECT timestamp, voltage, power, temperature, device_id"+
		" from measurement ORDER BY timestamp DESC LIMIT ?", count)

	if err != nil {
		log.Printf("Error while executing query %s", err.Error())
		return nil, err
	}

	var timestamp int64
	var voltage int
	var power int
	var temperature int
	var device_id int

	for rows.Next() {
		err := rows.Scan(&timestamp, &voltage, &power,
			&temperature, &device_id)
		if err != nil {
			log.Printf("Error while extracting value from row %s", err.Error())
			return nil, err
		}

		m := Measurement{
			Timestamp:   timestamp,
			Voltage:     voltage,
			Power:       power,
			Temperature: temperature,
			DeviceId:    device_id,
		}

		result = append(result, m)
	}
	rows.Close()
	tx.Commit()

	return result, nil
}

func (db DatabaseStorage) GetDeviceById(uuid string) (Device, error) {
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return Device{}, err
	}
	defer tx.Commit()

	var ipAddr string
	var id int

	row := tx.QueryRow("select id, uuid, ip_addr from device where uuid = ?", uuid)
	err = row.Scan(&id, &uuid, &ipAddr)

	if err != nil {
		log.Printf("Error while reading data from the row %s", err.Error())
		return Device{}, err
	}

	return Device{
		Id:        id,
		Uuid:      uuid,
		IpAddress: ipAddr,
	}, nil
}

func (db DatabaseStorage) UpdateDeviceIP(uuid, ipAddr string) error {
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}
	defer tx.Rollback()
	log.Printf("Update device uuid = %s ip address to %s", uuid, ipAddr)

	_, err = tx.Exec("update device set ip_addr = ? where uuid = ?", ipAddr, uuid)

	if err != nil {
		log.Printf("Error while updating device %s ip address to %s %s", uuid, ipAddr, err.Error())
		return err
	}

	tx.Commit()
	return nil
}

func (db DatabaseStorage) CreateDevice(uuid, ipAddress string, userId int) error {
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return err
	}
	defer tx.Rollback()

	log.Printf("Inserting new device with uuid %s", uuid)
	_, err = tx.Exec("INSERT INTO device(uuid, ip_addr, user_id) VALUES(?, ?, ?)", uuid,
		ipAddress,
		userId)

	if err != nil {
		log.Printf("Error while inserting device: %s", err.Error())
		return err
	}

	tx.Commit()

	return nil
}

func (db DatabaseStorage) GetMeasurementsByDevice(deviceId, count int) ([]Measurement, error) {
	result := make([]Measurement, 0, count)
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error while opening transaction message: %s", err.Error())
		return nil, err
	}

	rows, err := tx.Query("SELECT timestamp, voltage, power, temperature, device_id"+
		" from measurement WHERE device_id = ? ORDER BY timestamp DESC LIMIT ?", deviceId, count)

	if err != nil {
		log.Printf("Error while executing query %s", err.Error())
		return nil, err
	}

	var timestamp int64
	var voltage int
	var power int
	var temperature int
	var device_id int

	for rows.Next() {
		err := rows.Scan(&timestamp, &voltage, &power,
			&temperature, &device_id)
		if err != nil {
			log.Printf("Error while extracting value from row %s", err.Error())
			return nil, err
		}

		m := Measurement{
			Timestamp:   timestamp,
			Voltage:     voltage,
			Power:       power,
			Temperature: temperature,
			DeviceId:    device_id,
		}

		result = append(result, m)
	}
	rows.Close()
	tx.Commit()

	return result, nil
}
