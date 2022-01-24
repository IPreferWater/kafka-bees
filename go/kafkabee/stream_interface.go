package kafkabee

var (
	Stream Streaming
)

type Streaming interface {
	//TODO generics ?
	Produce(Data) error
	ProduceEuropeanBee(europeanBee) error
}
