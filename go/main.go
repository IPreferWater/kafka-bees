package main

import (
	"fmt"
	_ "image/png"

	"github.com/ipreferwater/kafka-bees/gui"
)

func main() {
	fmt.Println("go")
	//go kafkabee.Init()
	//go kafkabee.InitConsumer()

	gui.StartEbiten()
}
