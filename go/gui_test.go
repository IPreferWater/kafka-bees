package main

import (
	"math"
	"testing"
)

func TestGetCloserFromHive(t *testing.T) {
	hiveEntryX, hiveEntryY := float64(110), float64(140)
	positionsToTest := [4]coordinate{
		{
			x: 80,
			y: 120,
		},
		{
			x: 145,
			y: 112,
		},
		{
			x: 105,
			y: 160,
		},
		{
			x: 180,
			y: 180,
		},
	}

	
	for _, coordinateToTest := range positionsToTest {
		insect := Insect{
			position: coordinateToTest,
		}
		getCloserFromHive(&insect, hiveEntryX, hiveEntryY)

		distanceStart := getDistanceBeetweenTwoPoints(coordinateToTest.x, coordinateToTest.y, hiveEntryX, hiveEntryY)
		distanceGot := getDistanceBeetweenTwoPoints(insect.position.x, insect.position.y, hiveEntryX, hiveEntryY)

		if distanceStart < distanceGot {
			t.Errorf("the insect was %f pixels away and now %f, he should be closer\n", distanceStart, distanceGot)
		}
	}

}

func getDistanceBeetweenTwoPoints(x1, y1, x2, y2 float64) float64 {
	//d(A,B)=√(x2−x1)2+(y2−y1)2
	toSquare := math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)
	return math.Sqrt(toSquare)
}

func TestGetNumberOfTilesToDraw(t *testing.T){
	screenWidth, screenHeight  := 1024, 600
	xExpected, yExpected := 16,10 
	xGot ,yGot := getNumberOfTilesToDraw(screenWidth, screenHeight,64)

	if xExpected != xGot || yExpected != yGot {
		t.Errorf("expected x:%d and y:%d but got x: %d y : %d \n", xExpected, yExpected, xGot, yGot)
	}
}

func TestGetBeeIndexToBeKilled(t *testing.T){
	arr := []Insect{
{
	position:  coordinate{
		x: 110,
		y: 120,
	},
},
{
	position:  coordinate{
		x: 160,
		y: 170,
	},
},
	}
	expected := 1
	got := getInsectInRange(arr,155,165,5)

	if expected != got {
		t.Errorf("expected index %d but got %d", expected, got)
	}

	expected = 1
	got = getInsectInRange(arr,164,174,5)

	if expected != got {
		t.Errorf("expected index %d but got %d", expected, got)
	}

	expected = -1
	got = getInsectInRange(arr,164,176,5)

	if expected != got {
		t.Errorf("expected index %d but got %d", expected, got)
	}
}
