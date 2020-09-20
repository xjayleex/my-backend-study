package main

type Elevator struct {
	Id int
	Floor int
	Passengers []Call
	Status string
}

type Call struct {
	Id int
	Timestamp int
	Start int
	End int
}

type Start struct {
	Token string
	Timestamp int
	Elevators []Elevator
}

type Command struct {
	ElevatorID int
	Command string
	CallIDs []int
}