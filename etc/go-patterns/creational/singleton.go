package main

import (
	"fmt"
	"sync"
)

type singleton map[string]string

var (
	once sync.Once
	instance singleton
)

func NewMap() singleton {
	once.Do(func(){
		instance = make(singleton)
	})

	return instance
}

func main() {
	s := NewMap()
	s["this"] = "that"
	s2 := NewMap()
	fmt.Println("This is ", s2["this"])
}