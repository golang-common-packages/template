package stream

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/golang-microservices/template/model"
)

// MQTTStore manage all mqtt action
type MQTTStore struct {
	Config *model.Service
}

// MQTTProducerStore manage concrete producer action
type MQTTProducerStore struct {
	Client mqtt.Client
	Config *model.MQTT
}

// MQTTConsumerStore manage concrete consumer action
type MQTTConsumerStore struct {
	Client mqtt.Client
	Config *model.MQTT
}

/*
	@configMappingMQTT: Mapping between model.Service and MQTTStore for singleton pattern
	@sessionMappingMQTTProducer: Mapping between model.Service and MQTTProducerStore for singleton pattern
	@sessionMappingMQTTConsumer: Mapping between model.Service and MQTTConsumerStore for singleton pattern
*/
var (
	configMappingMQTT          = make(map[string]*MQTTStore)
	sessionMappingMQTTProducer = make(map[string]*MQTTProducerStore)
	sessionMappingMQTTConsumer = make(map[string]*MQTTConsumerStore)
)

// NewMQTTStore function for config mapping
func NewMQTTStore(config *model.Service) StreamStore {
	hash := config.Hash()
	currentConfig := configMappingMQTT[hash]

	if currentConfig == nil {
		currentConfig = &MQTTStore{config}
		configMappingMQTT[hash] = currentConfig
	}

	return currentConfig
}

// NewConsumerStore is concrete factory for consumer
func (k *MQTTStore) NewConsumerStore() ConsumerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingMQTTConsumer[hash]

	if currentSession == nil {
		client := connect(k.Config.MQTT.PrefixSub+strconv.FormatInt(time.Now().Unix(), 10), k.Config.MQTT.Host)
		log.Println("Connected to Consumer MQTT Server")

		currentSession = &MQTTConsumerStore{client, &k.Config.MQTT}
		sessionMappingMQTTConsumer[hash] = currentSession
	}

	return currentSession
}

// NewProducerStore is concrete factory for consumer
func (k *MQTTStore) NewProducerStore() ProducerStore {
	hash := k.Config.Hash()
	currentSession := sessionMappingMQTTProducer[hash]

	if currentSession == nil {
		client := connect(k.Config.MQTT.PrefixPub+strconv.FormatInt(time.Now().Unix(), 10), k.Config.MQTT.Host)
		log.Println("Connected to Producer MQTT Server")

		currentSession = &MQTTProducerStore{client, &k.Config.MQTT}
		sessionMappingMQTTProducer[hash] = currentSession
	}

	return currentSession
}

// connect function will establish MQTT connection
func connect(clientID string, uri string) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

// createClientOptions function will create client option to connect to MQTT Server
func createClientOptions(clientID string, uri string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetClientID(clientID)
	return opts
}

// Producer function will sent message based on topic and message
func (m *MQTTProducerStore) Producer(topic, message string) {
	m.Client.Publish(topic, m.Config.QoS, false, message)
}

// Consumer function will receive message bases on topic
func (m *MQTTConsumerStore) Consumer(topic string, c chan string) {
	m.Client.Subscribe(topic, m.Config.QoS, func(client mqtt.Client, msg mqtt.Message) {
		c <- fmt.Sprintf("%s|%s", msg.Topic(), string(msg.Payload()))
	})
}
