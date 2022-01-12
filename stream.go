package main

import "github.com/IPreferWater/Bric-A-Brac-World/kafkabee"

func sendDetectionToStream(t InsectType, hiveID int, direction bool) {

	if t == EuropeanBee {
		kafkabee.Stream.Produce(fakeBee(hiveID, direction))
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
