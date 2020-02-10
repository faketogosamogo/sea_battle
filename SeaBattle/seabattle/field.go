package seabattle

import (
	"encoding/json"
	"errors"
	"fmt"
)

type field struct {
	Coordinates coordinates
	Bot         *bot
}

func (f *field) initField() {
	f.Coordinates = coordinates{}
}

func (f field) getField() []byte {
	var data []byte
	if f.Bot == nil {
		data, _ = f.Coordinates.getCoordinatesJSON()
	} else {
		data, _ = f.Coordinates.getHiddenCoordinatesJSON()
	}
	return data
}
func (f *field) setField(data []byte) error {
	err := errors.New("У бота, не может быть ручного заполнения")
	if f.Bot != nil {
		return err
	}
	f.Coordinates = coordinates{}
	ships := make([]ship, 0)
	err = json.Unmarshal(data, &ships)
	if err != nil {
		return err
	}
	err = errors.New("Error of ships data")
	if len(ships) != 10 {
		return err
	}

	lenMap := make(map[int]int)
	for _, v := range ships {
		lenMap[v.Length]++
	}
	if lenMap[4] == 1 && lenMap[3] == 2 && lenMap[2] == 3 && lenMap[1] == 4 {
		for _, v := range ships {
			if !f.addShip(v.Length, v.X, v.Y, v.Direction) {
				return err
			}
		}
	} else {
		return err
	}
	return nil

}

func (f *field) getShot(x, y int) int {
	value := f.getValuePoint(x, y)
	if value == NOT_INSIDE {
		return NOT_INSIDE
	}
	if value == SHIP {
		f.setValuePoint(x, y, HIT)
		if f.isKilled(x, y) {
			return KILLED
		}
		return SHIP
	}
	if value == EMPTY {
		f.setValuePoint(x, y, MISS)
		return EMPTY
	}
	return value

}
func (f field) getValuePoint(x, y int) int {
	if !isInside(x, y) {
		return NOT_INSIDE
	}
	return f.Coordinates.getValue(x, y)
}

func (f *field) setValuePoint(x, y, value int) int {
	if !isInside(x, y) {
		return NOT_INSIDE
	}
	f.Coordinates.setValue(x, y, value)
	return INSIDE
}

func (f field) isEmpty(x, y int) bool {
	if f.Coordinates.getValue(x, y) == EMPTY {
		return true
	}
	return false
}
func (f field) isHit(x, y int) bool {
	if f.Coordinates.getValue(x, y) == HIT {
		return true
	}
	return false
}
func (f field) isMiss(x, y int) bool {
	if f.Coordinates.getValue(x, y) == MISS {
		return true
	}
	return false
}
func (f field) isShip(x, y int) bool {
	if f.Coordinates.getValue(x, y) == SHIP {
		return true
	}
	return false
}
func (f *field) isKilled(x, y int) bool {
	if !f.shipAroundPoint(x, y) { //работает только для корабля длины==1
		f.deletePoint(x, y)
		return true
	}

	direction := f.getShipDirection(x, y)
	if direction == -1 {
		return false
	}
	startX := x
	startY := y
	//fmt.Println(x, y, direction)
	for {
		if t1, t2 := move(startX, startY, direction); t1 == -1 && t2 == -1 {
			break
		}
		if f.isShip(move(startX, startY, direction)) {
			return false
		}
		if f.isEmpty(move(startX, startY, direction)) || f.isMiss(move(startX, startY, direction)) {
			break
		}
		if f.isHit(move(startX, startY, direction)) {
			startX, startY = move(startX, startY, direction)
		}
	}

	finishX := startX
	finishY := startY

	//fmt.Println(startX, startY)
	direction = reverseDirection(direction)

	for {
		if t1, t2 := move(finishX, finishY, direction); t1 == -1 && t2 == -1 {
			break
		}
		if f.isShip(move(finishX, finishY, direction)) {
			return false
		}
		if f.isEmpty(move(finishX, finishY, direction)) || f.isMiss(move(finishX, finishY, direction)) {
			break
		}
		if f.isHit(move(finishX, finishY, direction)) {
			finishX, finishY = move(finishX, finishY, direction)
		}
	}

	f.deleteShip(startX, startY, finishX, finishY, direction)

	return true
}
func (f field) hasPoints() bool {
	fmt.Println("HAS POINTS")
	for _, value := range f.Coordinates {
		if value == SHIP {
			return true
		}
	}
	fmt.Println("NET POINTS")
	return false
}
