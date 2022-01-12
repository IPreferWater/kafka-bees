package kafkabee

var (
	Stream Streaming
)

type Data struct {
	DataValue
	DataKey
}
type DataValue struct {
	Colors   map[string]float64
	Size     float64
	HasWings bool
}

type DataKey struct {
	HiveID int
	//true = in; false = out;
	Direction bool
}

type Streaming interface {
	//TODO generics ?
	Produce(Data) error
}
