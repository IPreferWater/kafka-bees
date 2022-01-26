package gui

import (
	"testing"
)

func TestGetNumberOfTilesToDraw(t *testing.T) {
	screenWidth, screenHeight := 1024, 600
	xExpected, yExpected := 16, 10
	xGot, yGot := getNumberOfTilesToDraw(screenWidth, screenHeight, 64)

	if xExpected != xGot || yExpected != yGot {
		t.Errorf("expected x:%d and y:%d but got x: %d y : %d \n", xExpected, yExpected, xGot, yGot)
	}
}
