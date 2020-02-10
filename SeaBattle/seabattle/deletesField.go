package seabattle

func (f field) shipAroundPoint(x, y int) bool {
	//d1, d2 := -1, -1
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}
			if !isInside(x+j, y+i) {
				continue
			}
			if f.getValuePoint(x+j, i+y) == SHIP || f.getValuePoint(x+j, i+y) == HIT {
				return true
			}
		}
	}

	return false
}
func (f field) getShipDirection(x, y int) int {

	if X, Y := toUp(x, y); X != -1 && Y != -1 {
		if f.isShip(toUp(x, y)) {
			return -1
		}
		if f.isHit(toUp(x, y)) {
			return UP
		}
	}
	if X, Y := toDown(x, y); X != -1 && Y != -1 {
		if f.isShip(toDown(x, y)) {
			return -1
		}
		if f.isHit(toDown(x, y)) {
			return DOWN
		}
	}
	if X, Y := toLeft(x, y); X != -1 && Y != -1 {
		if f.isShip(toLeft(x, y)) {
			return -1
		}
		if f.isHit(toLeft(x, y)) {
			return LEFT
		}
	}
	if X, Y := toRight(x, y); X != -1 && Y != -1 {
		if f.isShip(toRight(x, y)) {
			return -1
		}
		if f.isHit(toRight(x, y)) {
			return RIGHT
		}
	}
	return -1
}
func (f field) reverseDirection(direction int) int {
	if direction == UP {
		return DOWN
	}
	if direction == DOWN {
		return UP
	}
	if direction == LEFT {
		return RIGHT
	}
	if direction == RIGHT {
		return LEFT
	}
	return direction
}
func (f *field) deletePoint(x, y int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !isInside(x+j, y+i) {
				continue
			}
			if !f.isHit(x+j, y+i) {
				f.setValuePoint(x+j, y+i, MISS)
			}
		}
	}
}
func (f *field) deleteShip(startX, startY, finishX, finishY, direction int) {
	//fmt.Println(startX, startY, finishX, finishY, direction)
	if direction == UP {
		for i := startY; i >= finishY; i-- {
			f.deletePoint(startX, i)
		}
	}
	if direction == DOWN {
		for i := startY; i <= finishY; i++ {
			f.deletePoint(startX, i)
		}
	}
	if direction == LEFT {
		for i := startX; i >= finishX; i-- {
			f.deletePoint(i, startY)
		}
	}
	if direction == RIGHT {
		for i := startX; i <= finishX; i++ {
			f.deletePoint(i, startY)
		}
	}

}
