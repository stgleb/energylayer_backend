package storage

import (
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type KafkaStorage struct {
	*sql.DB
	MeasurementProducer sarama.AsyncProducer
}

func NewKafkaStorage() (KafkaStorage, error) {
	return KafkaStorage{}, nil
}

func (kafka KafkaStorage) CreateMeasurement(m Measurement) error {
	return nil
}

func (kafka KafkaStorage) GetMeasurement(count int) ([]Measurement, error) {
	return nil, nil
}

func (kafka KafkaStorage) GetDeviceById(uuid string) (Device, error) {
	return nil, nil
}

func (kafka KafkaStorage) CreateDevice(uuid, ipAddress string, userId int) error {
	return nil
}
