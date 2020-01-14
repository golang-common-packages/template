package stream

import (
	"github.com/golang-microservices/template/model"
)

// ProducerStore is abstract producer
type ProducerStore interface {
	Producer(topic, message string)
}

// ConsumerStore is abstract consumer
type ConsumerStore interface {
	Consumer(topic string, c chan string)
}

// StreamStore is abstract factory
type StreamStore interface {
	NewProducerStore() ProducerStore
	NewConsumerStore() ConsumerStore
}

// define type of stream package
const (
	KAFKA = iota
	MQTT
	RABBITMQ
)

// NewStreamClient function for Factory Pattern
func NewStreamClient(actionType int, config *model.Service) StreamStore {
	switch actionType {
	case KAFKA:
		return NewKafkaStore(config)
	case RABBITMQ:
		return NewRabbitMQStore(config)
	case MQTT:
		return NewMQTTStore(config)
	}
	return nil
}
