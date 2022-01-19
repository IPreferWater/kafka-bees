package kafkabee

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/riferrei/srclient"
)

func getSchema(client *srclient.SchemaRegistryClient, schemaName string) (*srclient.Schema, error) {
	schema, err := client.GetLatestSchema(schemaName)
	if err != nil {
		return nil, fmt.Errorf("error getLatestSchema %s", err)
	}

	fmt.Printf("get latest %s id = %d\n", schemaName, schema.ID())
	if schema == nil {
		//TODO I need to check this part
		fmt.Printf("schema %s is nil !\n", schemaName)
		schemaBytes, errReadFile := ioutil.ReadFile("complexType.avsc")
		if errReadFile != nil {
			return nil, fmt.Errorf("error read file %s", errReadFile)
		}
		schema, err = client.CreateSchema(schemaName, string(schemaBytes), srclient.Avro)
		if err != nil {
			return nil, fmt.Errorf("Error creating the schema %s", err)
		}
	}

	return schema, nil
}

func getValueByte(schema *srclient.Schema, toMarshall interface{}) []byte {

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	value, errJ := json.Marshal(toMarshall)
	if errJ != nil {
		panic(fmt.Sprintf("err jsonMarshall %s", errJ))
	}

	native, _, err := schema.Codec().NativeFromTextual(value)
	if err != nil {
		panic(fmt.Sprintf("err NativeFromTextual %s toMarshall %s\n", err, value))
	}
	valueBytes, err := schema.Codec().BinaryFromNative(nil, native)

	if err != nil {
		panic(fmt.Sprintf("err BinaryFromNative %s", err))
	}

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	return recordValue
}
