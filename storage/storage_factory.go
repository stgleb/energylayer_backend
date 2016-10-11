package storage

import (
	influx "./influx_storage"
	mysql "./mysql_storage"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"log"
)

const (
	USER     = "root"
	PASSWORD = "1234"
	DATABASE = "energylayer"
	PORT     = "3306"
	HOST     = "localhost"
)

func StorageFactory(storageType string) (Storage, error) {
	var uri string
	var dbType string

	switch storageType {
	case "mysql":
		uri = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASSWORD, HOST, PORT, DATABASE)
		dbType = "mysql"
		storage, err := newStorage(uri, dbType)

		if err != nil {
			return mysql.DatabaseStorage{storage}, err
		}

		return mysql.DatabaseStorage{storage}, nil
	case "memory":
		uri = ":memory:"
		dbType = "sqlite3"

		storage, err := newStorage(uri, dbType)

		if err != nil {
			return mysql.DatabaseStorage{storage}, err
		}

		return mysql.DatabaseStorage{storage}, nil

	case "influx":
		uri = "influx"
		dbType = "influx"
		return influx.NewInfluxDBStorage("", "", "", "", "")
	}

	return nil, errors.New("No such db drivers")
}

func newStorage(uri, dbType string) (*sql.DB, error) {
	db, err := sql.Open(dbType, uri)

	if err != nil {
		log.Printf("Error during connecting to database %s", err.Error())
		return db, err
	}

	return db, nil
}
