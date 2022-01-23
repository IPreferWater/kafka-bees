package kafkabee

import "fmt"

func guess(v DataValue, k DataKey) {
	if isEuropeanBee(v) {
		eB := europeanBee{
			HiveID:    k.HiveID,
			Size:      int(v.Size),
			Direction: k.Direction,
		}
		Stream.ProduceEuropeanBee(eB)
		return
	}

	if isAsianWasp(v) {
		fmt.Printf("send alert asian wasp on hiveID %d\n", k.HiveID)
	}
}

func isEuropeanBee(v DataValue) bool {

	if v.Size < 13 || v.Size > 16 {
		return false
	}

	if !v.HasWings {
		return false
	}

	if v.Colors["brown"] < 60 || v.Colors["brown"] > 80 && v.Colors["black"] < 20 || v.Colors["black"] > 40 {
		return false
	}
	return true
}

func isAsianWasp(v DataValue) bool {

	if v.Size < 18 || v.Size > 23 {
		return false
	}

	if !v.HasWings {
		return false
	}

	if v.Colors["orange"] < 5 || v.Colors["orange"] > 20 && v.Colors["black"] < 80 || v.Colors["black"] > 95 {
		return false
	}
	return true
}
