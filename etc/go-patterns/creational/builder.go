package main

import (
	"fmt"
)

type Color string
type Speed int

const (
	BlueColor Color = "blue"
	GreenColor		= "green"
	RedColor		= "red"
)

type Car interface {
	Drive() string
	Stop() string
}

type car struct {
	topSpeed Speed
	color Color
}

func (c *car) Drive() string{
	return "Driving at speed : " + c.topSpeed.String()
}

func (c *car) Stop() string {
	return "Stopping a " + string(c.color) + " car."
}

func (s Speed) String() string {
	return fmt.Sprintf("%d", s)
}

type CarBuilder interface {
	Color(Color)	CarBuilder
	TopSpeed(Speed)	CarBuilder
	Build()			Car
}

type carBuilder struct {
	speedOpt Speed
	color Color
}

func (cb *carBuilder) Color(color Color) CarBuilder{
	cb.color = color
	return cb
}

func (cb *carBuilder) TopSpeed(speed Speed) CarBuilder {
	cb.speedOpt = speed
	return cb
}

func (cb *carBuilder) Build() Car {
	return &car {
		topSpeed: cb.speedOpt,
		color: cb.color,
	}
}

func New() CarBuilder {
	return &carBuilder{}
}

func main() {
	assembly := New()
	sportsCar := assembly.TopSpeed(250).Color(RedColor).Build()
	toyCar := assembly.TopSpeed(10).Color(BlueColor).Build()

	fmt.Println(sportsCar.Drive())
	fmt.Println(sportsCar.Stop())
	fmt.Println(toyCar.Drive())
	fmt.Println(toyCar.Stop())
}