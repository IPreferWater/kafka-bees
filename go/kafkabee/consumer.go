package kafkabee

import (
	"encoding/binary"
	"encoding/json"
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

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			return (fmt.Errorf("error readin msg %s\n", err))
		}

		//value
		schemaID := binary.BigEndian.Uint32(msg.Value[1:5])

		schema, err := schemaRegistryClient.GetSchema(int(schemaID))

		if err != nil {
			panic(fmt.Errorf("Error getting the schema with id '%d' %s", schemaID, err))
		}

		native, _, _ := schema.Codec().NativeFromBinary(msg.Value[5:])
		value, _ := schema.Codec().TextualFromNative(nil, native)

		detectedValue := DataValue{}
		if err := json.Unmarshal(value, &detectedValue); err != nil {
			panic(fmt.Errorf("error unmarshall the string %s err => %s\n", string(value), err))
		}

		//key

		schemaIDKey := binary.BigEndian.Uint32(msg.Key[1:5])

		schemaKey, err := schemaRegistryClient.GetSchema(int(schemaIDKey))

		if err != nil {
			return fmt.Errorf("Error getting the schema key with id '%d' %s", schemaIDKey, err)
		}

		nativeKey, _, _ := schemaKey.Codec().NativeFromBinary(msg.Key[5:])
		key, _ := schemaKey.Codec().TextualFromNative(nil, nativeKey)

		detectedKey := DataKey{}
		if err := json.Unmarshal(key, &detectedKey); err != nil {
			panic(fmt.Sprintf("error unmarshall the string %s err => %s\n", string(key), err))
		}

		guess(detectedValue, detectedKey)

		//fmt.Printf("key : %#v\n value : %#v\n***\n", detectedKey, detectedValue)
	}

}


