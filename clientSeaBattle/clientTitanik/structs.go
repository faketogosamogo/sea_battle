package clientTitanik

import (
	"encoding/json"
	"fmt"
)

const (
	BOT_WIN   = 3
	USER_WIN  = 2
	USER_SHOT = 1
	BOT_SHOT  = 4
	BOT_KILL  = 5
	USER_KILL = 6

	HIT        = -1 //если здесь произошло попадание по кораблю
	EMPTY      = 0  //если на клетке ничего нет и не было
	NOT_INSIDE = -2
	INSIDE     = -3
	SHIP       = 1 //если здесь находится корабль
	MISS       = 2 //если здесь ничего не находилось
	KILLED     = 3

	UP    = 1
	DOWN  = 2
	LEFT  = 3
	RIGHT = 4

	LEN_SHIP = 4
	SIZE     = 10
)

type field struct {
	Coordinates Coordinates
}

func isInside(x, y int) bool {
	if (x >= 0 && x < SIZE) && (y >= 0 && y < SIZE) {
		return true
	}
	return false
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
func (f *field) setValuePoint(x, y, value int) int {
	if !isInside(x, y) {
		return NOT_INSIDE
	}
	f.Coordinates.setValue(x, y, value)
	return INSIDE
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
func (f field) printField() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			v := f.Coordinates.getValuePoint(j, i)
			if v == 0 {
				fmt.Print("* ")
			} else {
				fmt.Print("+ ")
			}
		}
		fmt.Println()
	}
}
func (c *Coordinates) setValue(x, y, value int) {
	(*c)[Point{x, y}] = value
}
func (c Coordinates) getValue(x, y int) int {
	return c[Point{x, y}]
}

type RoundResponse struct {
	State      int          `json:"state"`
	UserResult pointValue   `json: "userResult"`
	BotResults []pointValue `json:"botResults"`
}
type Ship struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction"`
	Length    int `json:"length"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type pointValue struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Value int `json:"value"`
}
type Coordinates map[Point]int

func (c *Coordinates) GetCoordinatesJSON() ([]byte, error) {
	pointsValue := make([]pointValue, 0)
	for key, value := range *c {
		pointsValue = append(pointsValue, pointValue{key.X, key.Y, value})
	}
	return json.Marshal(pointsValue)
}
func (c *Coordinates) SetCoordinatesJSON(data []byte) error {
	*c = Coordinates{}
	pointsValue := make([]pointValue, 0)
	err := json.Unmarshal(data, &pointsValue)
	if err != nil {
		return err
	}
	*c = Coordinates{}
	for _, value := range pointsValue {
		(*c)[Point{value.X, value.Y}] = value.Value
	}
	return nil
}
func (c *Coordinates) GetHiddenCoordinatesJSON() ([]byte, error) {
	pointsValue := make([]pointValue, 0)
	for key, value := range *c {
		if value != SHIP {
			pointsValue = append(pointsValue, pointValue{key.X, key.Y, value})
		}
	}
	return json.Marshal(pointsValue)
}

func (c Coordinates) getValuePoint(x, y int) int {
	return c[Point{x, y}]
}

//true если точка внутри поля!
func (c Coordinates) isInside(x, y int) bool {
	if (x < SIZE && x >= 0) && (y < SIZE && y >= 0) {
		return true
	}
	return false
}

///запись и чтение из файла записать
