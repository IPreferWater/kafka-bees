package main

import (
	"fmt"
	_ "image/png"

	"github.com/ipreferwater/kafka-bees/gui"
)

func main() {
	fmt.Println("go")

	/*if err := kafkabee.Init(); err != nil {
		log.Fatal(err)
	}

	if err := kafkabee.InitConsumer(); err != nil {
		log.Fatal(err)
	}*/

	gui.StartEbiten()
}
