package main

type Hive struct {
	ID int
	position  coordinate
	beesCount int
	beesToAdd int
	beesToRemove int
	//TODO create type
	waspsCount int
	waspsToAdd int
	waspsToRemove int
	//beesToCome []Bee
	//beesToGo []Bee
	insectsToCome map[InsecType][]Insect
	insectsToGo map[InsecType][]Insect
	hiveEntry coordinate
	hiveExit coordinate
}

type Insect struct {
	position coordinate
	waspState WaspStateType
}

type InsecType int

const (
    Bee InsecType = iota
    Wasp
)

type WaspStateType int

const (
    Approching WaspStateType = iota
    Hunting
	Leaving
)

type coordinate struct {
	x float64
	y float64
}
