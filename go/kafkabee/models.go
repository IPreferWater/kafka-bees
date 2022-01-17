package kafkabee

type Data struct {
	DataValue
	DataKey
}
type DataValue struct {
	Colors   map[string]float64 `json:"colors"`
	Size     float64            `json:"size"`
	HasWings bool               `json:"has_wings"`
}

type DataKey struct {
	HiveID int `json:"hive_id"`
	//true = in; false = out;
	Direction bool `json:"direction"`
}

type europeanBee struct {
	HiveID    int  `json:"hive_id"`
	Size      int  `json:"size"`
	Direction bool `json:"direction"`
}
