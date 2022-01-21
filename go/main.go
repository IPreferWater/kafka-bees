package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ipreferwater/kafka-bees/kafkabee"
)

func main() {
	fmt.Println("go")
	go kafkabee.Init()
	go kafkabee.InitConsumer()

	g := &Game{
		hives: []Hive{
			{
				ID: 1,
				position: coordinate{
					x: 100,
					y: 100,
				},
				beesCount:     1000,
				beesToAdd:     1,
				beesToRemove:  1,
				waspsCount:    0,
				waspsToAdd:    1,
				waspsToRemove: 1,
				hiveEntry: coordinate{
					x: 152,
					y: 190,
				},
				hiveExit: coordinate{
					x: 180,
					y: 190,
				},
				insectsToCome: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
				insectsToGo: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
			}},
		mapCenterX: 9,
		mapCenterY: 6,
		worldSpeed: 3,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bees-World")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
