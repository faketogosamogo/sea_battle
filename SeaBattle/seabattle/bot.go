package seabattle

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
)

type bot struct {
	//Field              field
	CheckedCoordinates coordinates
	//RepetitivePoints   //частоповторяющиеся точки
	RepetitivePoints []pointValue
	X                int
	Y                int
	Direction        int
	TrueAttempts     int //количество верных попыток
	IterPoints       int //чтобы двигаться по RepetitivePoints
}

func (b *bot) initBot() {
	b.CheckedCoordinates = coordinates{}
	b.RepetitivePoints = make([]pointValue, 0)
	b.X = -1
	b.Y = -1
	b.Direction = 1

}

func (b *bot) loadData() {
	RepetitivePointsB := coordinates{}
	file, err := ioutil.ReadFile("botData.txt")
	if err != nil {
		fmt.Println(err)
	}
	RepetitivePointsB.setCoordinatesJSON(file)
	for index, value := range RepetitivePointsB {
		b.RepetitivePoints = append(b.RepetitivePoints, pointValue{index.X, index.Y, value})
	}
	b.sortRepetitivePoints()
	fmt.Println(b.RepetitivePoints)
	
}
func (b *bot) sortRepetitivePoints() {
	sort.Slice(b.RepetitivePoints, func(i, j int) bool {
		return b.RepetitivePoints[i].Value > b.RepetitivePoints[j].Value
	})
}
func (b *bot) regulationData() {
	RepetitivePointsMF := coordinates{}

	file, err := ioutil.ReadFile("manualFilling.txt")
	if err != nil {
		return
	}
	RepetitivePointsMF.setCoordinatesJSON(file)
	RepetitivePointsB := coordinates{}
	file, err = ioutil.ReadFile("botData.txt")
	if err != nil {
		fmt.Println(err)
	}
	RepetitivePointsB.setCoordinatesJSON(file)
	for index, _ := range RepetitivePointsMF {
		RepetitivePointsB[index]++
	}
	for index, value := range RepetitivePointsB {
		b.RepetitivePoints = append(b.RepetitivePoints, pointValue{index.X, index.Y, value})
	}

	wFile, err := os.Create("botData.txt")
	if err != nil {
		return
	}
	jsData, _ := RepetitivePointsB.getCoordinatesJSON()
	wFile.Write(jsData)

}
func (b *bot) setResult(result int) {
	b.CheckedCoordinates[point{b.X, b.Y}] += 1
	if result == KILLED {
		b.X, b.Y = -1, -1
		b.TrueAttempts = 0
		return
	}

	if result != SHIP && b.TrueAttempts == 0 {
		b.X, b.Y = -1, -1
		return
	}

	if result == SHIP && b.TrueAttempts == 0 {
		b.TrueAttempts += 1
		b.Direction = incrDirection(b.X, b.Y, b.Direction)
		b.X, b.Y = move(b.X, b.Y, b.Direction)
		return
	}

	if result != SHIP && b.TrueAttempts == 1 {
		b.X, b.Y = move(b.X, b.Y, reverseDirection(b.Direction)) //двигаемся в обратном направлении

		b.Direction = incrDirection(b.X, b.Y, b.Direction)
		b.X, b.Y = move(b.X, b.Y, b.Direction)
		return
	}
	if result == SHIP && b.TrueAttempts > 0 {
		b.TrueAttempts += 1
		b.X, b.Y = move(b.X, b.Y, b.Direction)
		return
	}
	if result != SHIP && b.TrueAttempts > 1 {
		b.Direction = reverseDirection(b.Direction) //разворачиваемся
		for i := 0; i <= b.TrueAttempts; i++ {
			b.X, b.Y = move(b.X, b.Y, b.Direction)
		}
		return
	}

}
func (b *bot) makeShot() (int, int) {
	var x, y int

	if b.X == -1 && b.Y == -1 {
		for {
			if b.IterPoints < len(b.RepetitivePoints) {
				x = b.RepetitivePoints[b.IterPoints].X
				y = b.RepetitivePoints[b.IterPoints].Y
				b.CheckedCoordinates[point{x, y}]++
				b.IterPoints += 1
				break
			} else {

				x = rand.Intn(SIZE)
				y = rand.Intn(SIZE)
			}
			if b.CheckedCoordinates[point{x, y}] > 0 {

			} else {
				b.CheckedCoordinates[point{x, y}]++
				break
			}
		}
		b.X, b.Y = x, y
		return x, y
	} else {
		return b.X, b.Y
	}

}
