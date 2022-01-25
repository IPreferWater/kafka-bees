package kafkabee

var (
	Stream Streaming
)

type Streaming interface {
	Produce(Data) error
	ProduceEuropeanBee(europeanBee) error
}
