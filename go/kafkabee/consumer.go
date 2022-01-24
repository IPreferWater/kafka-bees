package kafkabee

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/riferrei/srclient"
)

func InitConsumer() error {

	conf := getConfig()

	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers": conf.kafkaUrl,
		"security.protocol": conf.securityProtocol,
		"group.id":          "consumer-detected",
		"auto.offset.reset": "latest",
	}

	if conf.securityProtocol != "PLAINTEXT" {
		kafkaConfigMap.SetKey("ssl.ca.location", conf.pathCaPem)
		kafkaConfigMap.SetKey("ssl.certificate.location", conf.pathServiceCert)
		kafkaConfigMap.SetKey("ssl.key.location", conf.pathServiceKey)
	}

	c, err := kafka.NewConsumer(kafkaConfigMap)
	if err != nil {
		return (fmt.Errorf("can't create consumer %s", err))
	}

	err = c.SubscribeTopics([]string{conf.topic}, nil)

	schemaRegistryClient := srclient.CreateSchemaRegistryClient(conf.schemaRegistryUrl)
	go startConsuming(c, schemaRegistryClient)
	return nil

}

func startConsuming(c *kafka.Consumer, schemaRegistryClient *srclient.SchemaRegistryClient) {
	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Println(fmt.Errorf("error readin msg %s\n", err))
			continue
		}


		detectedValue, err := decodeDataValueFromSchema(msg.Value, schemaRegistryClient)
		if err != nil {
			fmt.Println(err)
			continue
		}

		detectedKey, err := decodeDataKeyFromSchema(msg.Key, schemaRegistryClient)
		if err != nil {
			fmt.Println(err)
			continue
		}

		guess(detectedValue, detectedKey)
	}
}
