package stream

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"

	"github.com/golang-microservices/template/model"
)

// KafkaStore manage all redis action
type KafkaStore struct {
	Config *model.Service
}

// KafkaConsumerStore manage concrete consumer action
type KafkaConsumerStore struct {
	Session sarama.Consumer
	Config  *model.Service
}

// KafkaProducerStore manage concrete producer action
type KafkaProducerStore struct {
	Session sarama.SyncProducer
	Config  *model.Service
}

/*
	@configMappingKafka: Mapping between model.Service and KafkaStore for singleton pattern
	@sessionMappingKafkaProducer: Mapping between model.Service and KafkaProducerStore for singleton pattern
	@sessionMappingKafkaConsumer: Mapping between model.Service and KafkaConsumerStore for singleton pattern
*/
var (
	configMappingKafka          = make(map[string]*KafkaStore)
	sessionMappingKafkaProducer = make(map[string]*KafkaProducerStore)
	sessionMappingKafkaConsumer = make(map[string]*KafkaConsumerStore)
)

// NewKafkaStore function for config mapping
func NewKafkaStore(config *model.Service) StreamStore {
	hash := config.Hash()
	currentConfig := configMappingKafka[hash]

	if currentConfig == nil {
		currentConfig = &KafkaStore{config}
		configMappingKafka[hash] = currentConfig
	}

	return currentConfig
}

// NewConsumerStore is concrete factory for consumer
func (k *KafkaStore) NewConsumerStore() ConsumerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingKafkaConsumer[hash]

	if currentSession == nil {
		consumer, err := sarama.NewConsumer(k.Config.Kafka.Hosts, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println("Connected to Consumer Kafka Server")

		currentSession = &KafkaConsumerStore{consumer, k.Config}
		sessionMappingKafkaConsumer[hash] = currentSession
	}

	return currentSession
}

// NewProducerStore is concrete factory for producer
func (k *KafkaStore) NewProducerStore() ProducerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingKafkaProducer[hash]

	if currentSession == nil {
		config := sarama.NewConfig()
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer(k.Config.Kafka.Hosts, config)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Connected to Producer Kafka Server")

		currentSession = &KafkaProducerStore{producer, k.Config}
		sessionMappingKafkaProducer[hash] = currentSession
	}

	return currentSession
}

// Producer function will sent message based on topic and message
func (k *KafkaProducerStore) Producer(topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: k.Config.Kafka.Partition,
		Value:     sarama.StringEncoder(message),
	}
	_, _, err := k.Session.SendMessage(msg)
	if err != nil {
		log.Println(err)
	}
}

// Consumer function will receive message bases on topic
func (k *KafkaConsumerStore) Consumer(topic string, c chan string) {
	partitionList, err := k.Session.Partitions(topic)
	if err != nil {
		log.Println(err)
	}

	initialOffset := sarama.OffsetNewest
	go func() {
		for _, partition := range partitionList {
			pc, _ := k.Session.ConsumePartition(topic, partition, initialOffset)

			for message := range pc.Messages() {
				c <- fmt.Sprintf("%s|%s", string(message.Topic), string(message.Value))
			}
		}
	}()

	return
}
