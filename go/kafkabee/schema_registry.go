package kafkabee

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/riferrei/srclient"
)

func getSchema(client *srclient.SchemaRegistryClient, schemaName string) (*srclient.Schema, error) {
	schema, err := client.GetLatestSchema(schemaName)
	if err != nil {
		return nil, fmt.Errorf("error getting latestSchema %s => %s", schemaName, err)
	}

	if schema == nil {
		return nil, fmt.Errorf("schema %s doesn't exist => %s", schemaName, err)
	}

	return schema, nil
}

// transform interface to an array of bytes with the avro schema
func encodeDataFromSchema(schema *srclient.Schema, toMarshall interface{}) ([]byte, error) {

	//get a json format
	value, errJ := json.Marshal(toMarshall)
	if errJ != nil {
		return nil, fmt.Errorf("error transform %s to json with marshal => %s", toMarshall, errJ)
	}

	//json to Datum format
	native, _, err := schema.Codec().NativeFromTextual(value)
	if err != nil {
		return nil, fmt.Errorf("error transforming json Data %s to Datum with the provided schema => %s", value, err)
	}

	// in kafka, the 4 first bytes stock the schemaID to be retrieved with the schema-registry
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	var recordValue []byte
	//add kafka magic byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)

	valueBytes, err := schema.Codec().BinaryFromNative(recordValue, native)

	if err != nil {
		return nil, fmt.Errorf("error transforming Datum to []bytes with the provided schema => %s", err)
	}

	return valueBytes, nil
}

//TODO duplicated code
func decodeDataValueFromSchema(b []byte, scClient *srclient.SchemaRegistryClient) (DataValue, error){
	detectedValue := DataValue{}

	byteValue, err := decodeFromSchema(b, scClient)
	if err != nil {
		return detectedValue, err
	}
	
	// unmarshall json to given pointer
	if err := json.Unmarshal(byteValue, &detectedValue); err != nil {
		return detectedValue, fmt.Errorf("error unmarshall the value %s to the given pointer => %s", string(byteValue), err)
	}
	return detectedValue, nil
}

func decodeDataKeyFromSchema(b []byte, scClient *srclient.SchemaRegistryClient) (DataKey, error){
	detectedKey := DataKey{}

	byteKey, err := decodeFromSchema(b, scClient)
	if err != nil {
		return detectedKey, err
	}
	
	// unmarshall json to given pointer
	if err := json.Unmarshal(byteKey, &detectedKey); err != nil {
		return detectedKey, fmt.Errorf("error unmarshall the value %s to the given pointer => %s", string(byteKey), err)
	}
	return detectedKey, nil
}

func decodeFromSchema(valueB []byte, scClient *srclient.SchemaRegistryClient) ([]byte, error) {
	// extract schemaID from the message
	schemaID := binary.BigEndian.Uint32(valueB[1:5])

	// retrieve schema via ID in schema-registry
	schema, err := scClient.GetSchema(int(schemaID))

	if err != nil {
		return nil, fmt.Errorf("Error getting the schema with id '%d' %s", schemaID, err)
	}

	// []binary to Datum format
	native, _, err := schema.Codec().NativeFromBinary(valueB[5:])
	if err != nil {
		return nil, fmt.Errorf("Error binary to Datum =>  %s", err)
	}

	// Datum to json format
	json, err := schema.Codec().TextualFromNative(nil, native)
	if err != nil {
		return nil, fmt.Errorf("Error Datum to json => %s", err)
	}

	return json, nil

}
