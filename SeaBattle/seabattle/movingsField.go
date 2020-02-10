package seabattle

func isInside(x, y int) bool {
	if (x >= 0 && x < SIZE) && (y >= 0 && y < SIZE) {
		return true
	}
	return false
}
func toUp(x, y int) (int, int) {
	if !isInside(x, y-1) {
		return -1, -1
	}
	return x, y - 1
}
func toDown(x, y int) (int, int) {
	if !isInside(x, y+1) {
		return -1, -1
	}
	return x, y + 1
}
func toLeft(x, y int) (int, int) {
	if !isInside(x-1, y) {
		return -1, -1
	}
	return x - 1, y
}
func toRight(x, y int) (int, int) {
	if !isInside(x+1, y) {
		return -1, -1
	}
	return x + 1, y
}
func move(x, y, direction int) (int, int) {
	if direction == UP {
		return toUp(x, y)
	}
	if direction == DOWN {
		return toDown(x, y)
	}
	if direction == LEFT {
		return toLeft(x, y)
	}
	if direction == RIGHT {
		return toRight(x, y)
	}
	return -1, -1
}
func incrDirection(x, y, direction int) int {

	for {
		if direction == RIGHT {
			direction = UP
		} else {
			direction += 1
		}
		a, b := move(x, y, direction)
		if a != -1 && b != -1 {
			break
		}
	}
	return direction

}
func reverseDirection(direction int) int {
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
