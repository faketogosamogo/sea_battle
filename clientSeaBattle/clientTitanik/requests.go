package clientTitanik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	makeShotURL  = "http://127.0.0.1:8080/makeShot"
	userFieldURL = "http://127.0.0.1:8080/userField"
	botFieldURL  = "http://127.0.0.1:8080/botField"
	startGameURL = "http://127.0.0.1:8080/startGame?choise="
	loadGameURL  = "http://127.0.0.1:8080/loadGame"
)

type ClientTitanik struct {
	BotCoordinates    Coordinates
	ClientCoordinates Coordinates
}

func (c *ClientTitanik) InitClient() {
	c.BotCoordinates = Coordinates{}
	c.ClientCoordinates = Coordinates{}
}

func (c *ClientTitanik) StartGame() {
	var choise int
	f := field{Coordinates{}}
	fmt.Println("1:Автозаполнение вашего поля/2:Ручное заполнение вашего поля/3:Загрузка последней игры")
	fmt.Scan(&choise)
	fmt.Println(startGameURL + "auto")
	if choise == 1 {
		_, err := http.Post(startGameURL+"auto", "application/json", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if choise == 2 {
		ships := make([]Ship, 0)
		var x, y, dir  int
		k := 1
		for i := LEN_SHIP; i > 0; i-- {
			for j := 0; j < k; j++ {
				fmt.Println("Введите x,y,направление (UP = 1, DOWN  = 2, LEFT  = 3, RIGHT = 4) для корабля с длиной: ", i)
				_, err := fmt.Scan(&x, &y, &dir)
				if err != nil {
					i--
					continue
				}
				if !f.addShip(i, x, y, dir) {
					j--
					continue
				}
				ships = append(ships, Ship{x, y, dir, i})
				f.printField()
			}
			k++
		}
		data, err := json.Marshal(&ships)
		if err != nil {
			fmt.Println(err)
		}

		_, err = http.Post(startGameURL+"manual", "application/json", bytes.NewBuffer(data))

		if err != nil {
			fmt.Println(err)
			c.StartGame()
			return
		}
	} else if choise == 3 {
		_, err := http.Post(loadGameURL, "application/json", nil)
		if err != nil {
			fmt.Println(err)
			c.StartGame()
		}
	} else {
		return
	}
	c.getUserCoordinates()
	c.getBotCoordinates()
	c.printFields()
	c.makeShot()
}

func (c *ClientTitanik) getBotCoordinates() {
	resp, err := http.Get(botFieldURL)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = c.BotCoordinates.SetCoordinatesJSON(body)
	if err != nil {
		fmt.Println(err)
	}

}
func (c *ClientTitanik) getUserCoordinates() {
	resp, err := http.Get(userFieldURL)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = c.ClientCoordinates.SetCoordinatesJSON(body)
	if err != nil {
		fmt.Println(err)
	}
}
func (c ClientTitanik) printField(coordinates Coordinates) {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			value := coordinates.getValuePoint(j, i)
			if value == EMPTY {
				fmt.Print("*  ")
			} else if value == SHIP {
				fmt.Print("+  ")
			} else if value == HIT {
				fmt.Print("-  ")
			} else if value == MISS {
				fmt.Print("/  ")
			}
		}
		fmt.Println()
	}
}

func (c ClientTitanik) printFields() {
	fmt.Println("Ваше поле!")
	c.printField(c.ClientCoordinates)
	fmt.Println("//////////////////////////")
	fmt.Println("Поле бота!")
	fmt.Println("//////////////////////////")
	c.printField(c.BotCoordinates)
}
func (c *ClientTitanik) makeShot() {
	var x, y int
	fmt.Println("Введите x,y")
	fmt.Scan(&x, &y)

	p := new(Point)
	p.X = x
	p.Y = y
	jsPoint, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Неверные координаты!")
		c.makeShot()
		return
	}

	buffer := bytes.Buffer{}
	buffer.Write(jsPoint)

	resp, err := http.Post(makeShotURL, "application/json", &buffer)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		fmt.Println("Попробуйте заново походить!")
		c.makeShot()
	}
	roundResp := new(RoundResponse)

	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	if err != nil {
		fmt.Println(err)
		fmt.Println("Попробуйте заново походить!")
		c.makeShot()
	}

	err = json.Unmarshal(body, roundResp)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Попробуйте заново походить!")
		c.makeShot()
		return
	}
	if roundResp.State == BOT_WIN {
		fmt.Println("Вас выиграл бот!")
		c.printField(c.ClientCoordinates)
		c.StartGame()
		return
	}
	if roundResp.State == USER_WIN {
		fmt.Println("Вы выиграли!")
		c.StartGame()
		return
	}

	c.getBotCoordinates()
	c.getUserCoordinates()

	fmt.Println("Результат вашего хода!")

	fmt.Println(roundResp.UserResult)

	if len(roundResp.BotResults) != 0 {
		fmt.Println("Ходы бота!")
		for _, value := range roundResp.BotResults {
			fmt.Println(value)
		}
	}
	c.printFields()
	fmt.Println("STATE: ", roundResp.State)
	c.makeShot()
}
