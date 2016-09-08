package storage

import (
	"database/sql"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type KafkaStorage struct {
	*sql.DB
	MeasurementProducer sarama.AsyncProducer
}

func newAccessLogProducer(brokerList []string) (sarama.AsyncProducer, error) {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal        // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionGZIP      // Compress messages
	config.Producer.Flush.Frequency = 1000 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer, nil
}

func NewKafkaStorage(brokerList []string) (KafkaStorage, error) {
	accessLogProducer, err := newAccessLogProducer(brokerList)

	if err != nil {
		log.Printf("Error while creating Access logger %s", err.Error())
		return KafkaStorage{}, err
	}

	return KafkaStorage{
		MeasurementProducer: accessLogProducer,
	}, nil
}

func (kafka KafkaStorage) Close() error {
	if err := kafka.MeasurementProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
		return err
	}

	return nil
}

func (kafka KafkaStorage) CreateMeasurement(m *Measurement, topicName string) error {
	kafka.MeasurementProducer.Input() <- &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(m.Timestamp),
		Value: m,
	}

	return nil
}

func (kafka KafkaStorage) GetMeasurement(count int) ([]Measurement, error) {
	return nil, nil
}

func (kafka KafkaStorage) GetDeviceById(uuid string) (Device, error) {
	return Device{}, nil
}

func (kafka KafkaStorage) CreateDevice(uuid, ipAddress string, userId int) error {
	return nil
}
