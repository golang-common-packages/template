package model

import (
	"strconv"
	"strings"

	"github.com/golang-common-packages/template/common/util/hash"
)

// toString private function convert Service config to string
func (s Service) toString(t interface{}) string {
	switch t.(type) {
	case MongoDB:
		return s.Database.MongoDB.User + s.Database.MongoDB.Password + strings.Join(s.Database.MongoDB.Hosts, ",") + s.Database.MongoDB.DB + strings.Join(s.Database.MongoDB.Options, ",")
	case Postgres:
		return s.Database.Postgres.User + s.Database.Postgres.Password + strings.Join(s.Database.Postgres.Host, ",") + s.Database.Postgres.DBName + s.Database.Postgres.Port
	case Redis:
		return s.Database.Redis.Prefix + s.Database.Redis.Password + s.Database.Redis.Host + string(s.Database.Redis.DB)
	case Email:
		return s.Email.Host + s.Email.Port + s.Email.Username + s.Email.Password + s.Email.Key + s.Email.From + s.Email.Subject + s.Email.Message
	case Fluent:
		return s.Fluent.Host + strconv.Itoa(s.Fluent.Port) + s.Fluent.Prefix + s.Fluent.Tag
	case Sharepoint:
		return s.Database.Sharepoint.SiteURL + s.Database.Sharepoint.Username + s.Database.Sharepoint.Password
	case Kafka:
		return strings.Join(s.Kafka.Hosts, ",") + strconv.Itoa(int(s.Kafka.Partition))
	case RabbitMQ:
		return s.RabbitMQ.Host + s.RabbitMQ.User + s.RabbitMQ.Password + s.RabbitMQ.Route
	case PGO:
		return s.PGO.SpaceName
	case MQTT:
		return s.MQTT.Host + strconv.Itoa(int(s.MQTT.QoS)) + s.MQTT.PrefixPub + s.MQTT.PrefixSub
	}
	return ""
}

// Hash function return hash of Service config
func (s Service) Hash() string {
	hashclient := hash.Client{}
	return hashclient.SHA256(s.toString(s))
}
