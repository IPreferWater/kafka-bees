package main

import (
	"github.com/ipreferwater/kafka-bees/kafkabee"
)

//import "github.com/ipreferwater/kafka-bees/kafkabee"

func sendDetectionToStream(t InsectType, hiveID int, direction bool) {

	if kafkabee.Stream == nil {
		//TODO can't send to stream
		return
	}
	if t == EuropeanBee {
		kafkabee.Stream.Produce(fakeBee(hiveID, direction))
		//fmt.Printf("kafka produce be direction %t hive %d\n", direction, hiveID)
		return
	}
}

func fakeBee(hiveID int, direction bool) kafkabee.Data {
	mapColors := make(map[string]float64, 2)
	percentageBrown := randomNumberBeetween(60, 80)
	percentageBlack := 100 - percentageBrown
	mapColors["brown"] = percentageBrown
	mapColors["black"] = percentageBlack

	dataValue := kafkabee.DataValue{
		Colors:   mapColors,
		Size:     randomNumberBeetween(13, 16),
		HasWings: true,
	}

	dataKey := kafkabee.DataKey{
		HiveID:    hiveID,
		Direction: direction,
	}

	return kafkabee.Data{
		DataValue: dataValue,
		DataKey:   dataKey,
	}

}

func fakeAsianWasp(hiveID int, direction bool) kafkabee.Data {
	mapColors := make(map[string]float64, 2)
	percentageOrange := randomNumberBeetween(5, 20)
	percentageBlack := 100 - percentageOrange
	mapColors["orange"] = percentageOrange
	mapColors["black"] = percentageBlack

	dataValue := kafkabee.DataValue{
		Colors:   mapColors,
		Size:     randomNumberBeetween(18, 23),
		HasWings: true,
	}

	dataKey := kafkabee.DataKey{
		HiveID:    hiveID,
		Direction: direction,
	}

	return kafkabee.Data{
		DataValue: dataValue,
		DataKey:   dataKey,
	}

}
