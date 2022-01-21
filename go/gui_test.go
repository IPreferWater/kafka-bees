package main

import (
	"testing"
)

//TODO we should get the distance instead to test every situations
func TestGetCloserFromHive(t *testing.T) {
	hiveEntryX, hiveEntryY := float64(110), float64(140)
	xStart, yStart := float64(100), float64(200)

	insect := Insect{
		position: coordinate{
			x: float64(xStart),
			y: float64(yStart),
		},
	}
	getCloserFromHive(&insect, hiveEntryX, hiveEntryY)

	if insect.position.x <= xStart || insect.position.y <= yStart {
		t.Errorf("x & y should increased but xStart %f current x %f & yStart %f current y %f\n", xStart, insect.position.x, yStart, insect.position.y)
	}

}
