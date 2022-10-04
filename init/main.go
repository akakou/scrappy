package main

import (
	"core"
	"log"
)

func main() {
	err := core.Setup()

	if err != nil {
		log.Fatalf("%v: ", err)
	}
}
