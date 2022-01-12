package kafkabee

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/riferrei/srclient"
)

type KafkaStreaming struct {
	producer    *kafka.Producer
	scClient    *srclient.SchemaRegistryClient
	valueSchema *srclient.Schema
	keySchema   *srclient.Schema
	config      kafkaConfig
}

type kafkaConfig struct {
	kafkaUrl          string
	schemaRegistryUrl string
	securityProtocol  string

	pathCaPem       string
	pathServiceCert string
	pathServiceKey  string

	topic           string
	schemaNameValue string
	schemaNameKey   string
}

func Init() {
	conf := getConfig()

	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers": conf.kafkaUrl,
		"security.protocol": conf.securityProtocol,
	}

	if conf.securityProtocol != "PLAINTEXT" {
		kafkaConfigMap.SetKey("ssl.ca.location", conf.pathCaPem)
		kafkaConfigMap.SetKey("ssl.certificate.location", conf.pathServiceCert)
		kafkaConfigMap.SetKey("ssl.key.location", conf.pathServiceKey)
	}

	p, err := kafka.NewProducer(kafkaConfigMap)
	if err != nil {
		panic(fmt.Sprintf("error creating producer %s", err))
	}

	// this will check the status of the sent messages
	/*go func() {
		for event := range p.Events() {
			switch ev := event.(type) {
			case *kafka.Message:
				m := ev
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Error delivering the message '%s'\n error %s\n", m.Key, ev.TopicPartition.Error)
				}
				// uncoment this part to debug your succes
				else {
					fmt.Printf("Message delivered successfully!\n key : %s\n headers : %s\n opaque : %s\n timestamp: %s\n value : %s \n offet : %s \n partition : %d\n topic : %s\n", m.Key, m.Headers, m.Opaque, m.Timestamp, m.Value, m.TopicPartition.Offset, m.TopicPartition.Partition, m.TopicPartition.Topic)
				}
			}
		}
	}()*/
	defer p.Close()

	schemaRegistryClient := srclient.CreateSchemaRegistryClient(conf.schemaRegistryUrl)

	Stream = KafkaStreaming{
		producer:    p,
		scClient:    schemaRegistryClient,
		valueSchema: getSchema(schemaRegistryClient, conf.schemaNameValue),
		keySchema:   getSchema(schemaRegistryClient, conf.schemaNameKey),
	}

}

func (k KafkaStreaming) Produce(d Data) error {

	recordValue := getValueByte(k.valueSchema, d.DataValue)
	recordKey := getValueByte(k.keySchema, d.DataKey)

	errProduce := k.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &k.config.topic, Partition: -1},
			Value:          recordValue,
			Key:            recordKey,
		}, nil)

	if errProduce != nil {
		panic(fmt.Sprintf("errProduce %s", errProduce))
	}

	return nil
}

func getConfig() kafkaConfig {
	return kafkaConfig{
		kafkaUrl:          getEnvAndWarnIfMissing("KAFKA_URL"),
		schemaRegistryUrl: getEnvAndWarnIfMissing("SCHEMA_REGISTRY_URL"),
		securityProtocol:  getEnvAndWarnIfMissing("SECURITY_PROTOCOL"),
		pathCaPem:         getEnvAndWarnIfMissing("CERTS_PATH_CA_PEM"),
		pathServiceCert:   getEnvAndWarnIfMissing("CERTS_PATH_SERVICE_CERT"),
		pathServiceKey:    getEnvAndWarnIfMissing("CERTS_PATH_SERVICE_KEY"),
		topic:             getEnvAndWarnIfMissing("TOPIC"),
		schemaNameValue:   getEnvAndWarnIfMissing("SCHEMA_NAME_VALUE"),
		schemaNameKey:     getEnvAndWarnIfMissing("SCHEMA_NAME_KEY"),
	}
}

func getEnvAndWarnIfMissing(s string) string {
	v := os.Getenv(s)

	if v == "" {
		log.Printf("WARN : %s env is not set", s)
	}

	return v
}
