package stream

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/golang-common-packages/template/model"
)

// RabbitMQStore manage all redis action
type RabbitMQStore struct {
	Config *model.Service
}

// RabbitMQProducerStore manage concrete producer action
type RabbitMQProducerStore struct {
	Channel *amqp.Channel
	Config  *model.RabbitMQ
}

// RabbitMQConsumerStore manage concrete consumer action
type RabbitMQConsumerStore struct {
	Channel *amqp.Channel
	Config  *model.RabbitMQ
}

/*
	@configMappingRabbitMQ: Mapping between model.Service and RabbitMQStore for singleton pattern
	@sessionMappingRabbitMQProducer: Mapping between model.Service and RabbitMQProducerStore for singleton pattern
	@sessionMappingRabbitMQConsumer: Mapping between model.Service and RabbitMQConsumerStore for singleton pattern
*/
var (
	configMappingRabbitMQ          = make(map[string]*RabbitMQStore)
	sessionMappingRabbitMQProducer = make(map[string]*RabbitMQProducerStore)
	sessionMappingRabbitMQConsumer = make(map[string]*RabbitMQConsumerStore)
)

// NewRabbitMQStore function for config mapping
func NewRabbitMQStore(config *model.Service) StreamStore {
	hash := config.Hash()
	currentConfig := configMappingRabbitMQ[hash]

	if currentConfig == nil {
		currentConfig = &RabbitMQStore{config}
		configMappingRabbitMQ[hash] = currentConfig
	}

	return currentConfig
}

// NewConsumerStore function for concrete factory for consumer
func (k *RabbitMQStore) NewConsumerStore() ConsumerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingRabbitMQConsumer[hash]

	if currentSession == nil {

		currentSession = &RabbitMQConsumerStore{RabbitMQDial(k.Config), &k.Config.RabbitMQ}
		sessionMappingRabbitMQConsumer[hash] = currentSession
	}

	return currentSession
}

// NewProducerStore function for concrete factory for consumer
func (k *RabbitMQStore) NewProducerStore() ProducerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingRabbitMQProducer[hash]

	if currentSession == nil {

		currentSession = &RabbitMQProducerStore{RabbitMQDial(k.Config), &k.Config.RabbitMQ}
		sessionMappingRabbitMQProducer[hash] = currentSession
	}

	return currentSession
}

// RabbitMQDial function will establish RabbitMQ connection
func RabbitMQDial(config *model.Service) *amqp.Channel {
	conn, err := amqp.Dial("amqp://" + config.RabbitMQ.User + ":" + config.RabbitMQ.Password + "@" + config.RabbitMQ.Host)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	fmt.Println("Connected to RabbitMQ Server")
	return ch
}

// Producer function will sent message based on topic and message
func (c *RabbitMQProducerStore) Producer(topic, message string) {
	err := c.Channel.ExchangeDeclare(
		topic,   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	err = c.Channel.Publish(
		topic,          // exchange
		c.Config.Route, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Sent: %s", message)
}

// Consumer function will receive message bases on topic
func (r *RabbitMQConsumerStore) Consumer(topic string, c chan string) {
	err := r.Channel.ExchangeDeclare(
		topic,   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	q, err := r.Channel.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = r.Channel.QueueBind(
		q.Name,         // queue name
		r.Config.Route, // routing key
		topic,          // exchange
		false,
		nil)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}

	msgs, err := r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Recieve message: %s", d.Body)
			c <- fmt.Sprintf("%s|%s", topic, string(d.Body))
		}
	}()

	return
}
