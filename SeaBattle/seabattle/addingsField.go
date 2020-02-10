package seabattle

import (
	"math/rand"
)

func (f field) checkPoint(x, y int) bool {
	if !isInside(x, y) {
		return false
	}
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !isInside(x+j, y+i) {
				continue
			}
			if !f.isEmpty(x+j, y+i) && !f.isMiss(x+j, y+i) {
				return false
			}
		}
	}
	return true
}
func (f field) checkArea(lenght, x, y, direction int) bool {
	if direction == UP {
		if y-lenght+1 < 0 {
			return false
		}
		for i := y; i > y-lenght; i-- {
			if !f.checkPoint(x, i) {
				return false
			}
		}
		return true
	}
	if direction == DOWN {
		if y+lenght-1 > 9 {
			return false
		}
		for i := y; i < y+lenght; i++ {
			if !f.checkPoint(x, i) {
				return false
			}
		}
		return true
	}
	if direction == LEFT {
		if x-lenght+1 < 0 {
			return false
		}
		for i := x; i > x-lenght; i-- {
			if !f.checkPoint(i, y) {
				return false
			}
		}
		return true
	}
	if direction == RIGHT {
		if x+lenght-1 > 9 {
			return false
		}
		for i := x; i < x+lenght; i++ {
			if !f.checkPoint(i, y) {
				return false
			}
		}
		return true
	}
	return false
}
func (f *field) addShip(lenght, x, y, direction int) bool {
	if direction > 4 || direction < 1 {
		return false
	}
	if !f.checkPoint(x, y) {
		return false
	}
	if !f.checkArea(lenght, x, y, direction) {
		return false
	}

	if direction == UP {
		for i := y; i > y-lenght; i-- {
			f.setValuePoint(x, i, SHIP)
		}
	}
	if direction == DOWN {
		for i := y; i < y+lenght; i++ {
			f.setValuePoint(x, i, SHIP)
		}
	}
	if direction == LEFT {
		for i := x; i > x-lenght; i-- {
			f.setValuePoint(i, y, SHIP)

		}
	}
	if direction == RIGHT {
		for i := x; i < x+lenght; i++ {
			f.setValuePoint(i, y, SHIP)
		}
	}
	return true
}
func (f *field) generateField() {
	f.Coordinates = coordinates{}
	k := 1
	for i := LEN_SHIP; i > 0; i-- {
		for j := 0; j < k; j++ {
			x := rand.Intn(SIZE)
			y := rand.Intn(SIZE)
			direction := rand.Intn(4)
			if !f.addShip(i, x, y, direction) {
				j--
			}
		}
		k++
	}
}
