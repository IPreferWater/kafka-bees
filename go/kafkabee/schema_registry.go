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
func getValueByte(schema *srclient.Schema, toMarshall interface{}) ([]byte, error) {

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
