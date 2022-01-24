package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/ipreferwater/kafka-bees/gui"
	"github.com/ipreferwater/kafka-bees/kafkabee"
)

func main() {
	fmt.Println("go")

	if err := kafkabee.Init(); err != nil {
		log.Fatal(err)
	}

	if err := kafkabee.InitConsumer(); err != nil {
		log.Fatal(err)
	}

	gui.StartEbiten()
}
