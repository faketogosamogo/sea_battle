package main

import (
	"math/rand"
	"time"

	"./seabattle"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	h := seabattle.HttpRest{}
	h.StartServer()
}
