package seabattle

import (
	"encoding/json"
)

const (
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

type ship struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction"`
	Length    int `json:"length"`
}

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type pointValue struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Value int `json:"value"`
}

type coordinates map[point]int

func (c *coordinates) getCoordinatesJSON() ([]byte, error) {
	pointsValue := make([]pointValue, 0)
	for key, value := range *c {
		pointsValue = append(pointsValue, pointValue{key.X, key.Y, value})
	}
	return json.Marshal(pointsValue)
}
func (c *coordinates) setCoordinatesJSON(data []byte) error {
	//*c = coordinates{}
	pointsValue := make([]pointValue, 0)
	err := json.Unmarshal(data, &pointsValue)
	if err != nil {
		return err
	}
	*c = coordinates{}
	for _, value := range pointsValue {
		(*c)[point{value.X, value.Y}] = value.Value
	}
	return nil
}

func (c *coordinates) getHiddenCoordinatesJSON() ([]byte, error) {
	pointsValue := make([]pointValue, 0)
	for key, value := range *c {
		if value != SHIP {
			pointsValue = append(pointsValue, pointValue{key.X, key.Y, value})
		}
	}
	return json.Marshal(pointsValue)

}

func (c *coordinates) setValue(x, y, value int) {
	(*c)[point{x, y}] = value
}
func (c coordinates) getValue(x, y int) int {
	return c[point{x, y}]
}

/*эту функцию перенести
func (c *Coordinates) GetHiddenCoordinatesJSON() ([]byte, error) {
	pointsValue := make([]pointValue, 0)
	for key, value := range *c {
		if value != ship {
			pointsValue = append(pointsValue, pointValue{key.X, key.Y, value})
		}
	}
	return json.Marshal(pointsValue)
	Coordinates[Point{1,1}]
}
*/
