package main

type InsectType int

const (
    EuropeanBee InsectType = iota
    AsianWasp
)


type Hive struct {
	ID int
	position  coordinate
	beesCount int
	beesToAdd int
	beesToRemove int
	//beesToCome []Bee
	//beesToGo []Bee
	insectsToCome map[InsecType][]Insect
	insectsToGo map[InsecType][]Insect
	hiveEntry coordinate
	hiveExit coordinate
}

type Insect struct {
	position coordinate

}

type InsecType int

const (
    Bee InsecType = iota
    Wasp
)


type coordinate struct {
	x float64
	y float64
}
