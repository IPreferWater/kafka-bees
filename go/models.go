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
	beesToCome []Bee
	beesToGo []Bee
}

type Bee struct {
	position coordinate
}

type coordinate struct {
	x float64
	y float64
}
