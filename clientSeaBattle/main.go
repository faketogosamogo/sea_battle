package main

import (
	"./clientTitanik"
)

func main() {
	client := clientTitanik.ClientTitanik{}
	client.InitClient()

	client.StartGame()

}
