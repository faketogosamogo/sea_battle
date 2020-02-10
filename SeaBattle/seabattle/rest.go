package seabattle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type HttpRest struct {
	Referee referee
	Started bool
	
}
func (h *HttpRest) StartServer() {
	h.Referee.init()
	http.HandleFunc("/startGame", h.startGame)
	http.HandleFunc("/userField", h.getUserField)
	http.HandleFunc("/botField", h.getBotField)
	http.HandleFunc("/makeShot", h.makeShotBot)
	http.ListenAndServe("localhost:8080", nil)
}
func (h *HttpRest) startGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.Referee.Field2.generateField()
		choise := r.URL.Query().Get("choise")
		if choise != "auto" && choise != "manual" {
			http.Error(w, "Неверный параметр", 400)
			return
		}
		h.Referee.Field2.Bot.initBot()
		if choise == "auto" {
			h.Referee.Field1.generateField()
			h.Referee.Field2.Bot.loadData()
			//	fmt.Fprintln(w, "Совершайте ход!") ////////
		}
		if choise == "manual" {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", 400)
				return
			}

			err = h.Referee.Field1.setField(data)
			if err != nil {
				http.Error(w, "Неверные данные", 422)
				return
			}

			file, err := os.Create("manualFilling.txt")
			if err != nil {
				http.Error(w, "Ошибка создания файла", 500)
				return
			}
			jsonCoordinates, _ := h.Referee.Field1.Coordinates.getCoordinatesJSON()
			_, err = file.Write(jsonCoordinates)

			if err != nil {
				http.Error(w, "Ошибка записи в файл", 500)
				return
			}
			h.Referee.Field2.Bot.regulationData()
			//h.Started = true
		}
		h.Started = true
	default:
		http.Error(w, "POST", 405)
		return
	}
}

func (h *HttpRest) getUserField(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !h.Started {
			http.Error(w, "Игра не начата!", 409)
			return
		}
		JSON := h.Referee.Field1.getField()
		fmt.Fprintln(w, string(JSON))
	default:
		http.Error(w, "GET", 405)
		return
	}
}
func (h *HttpRest) getBotField(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !h.Started {
			http.Error(w, "Игра не начата!", 409)
			return
		}
		JSON := h.Referee.Field2.getField()
		fmt.Fprintln(w, string(JSON))
	default:
		http.Error(w, "GET", 405)
		return
	}
}
func (h *HttpRest) makeShotBot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if !h.Started {
			http.Error(w, "Игра не начата!", 409)
			return
		}
		p := new(point)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Переданы неверные данные", 422)
			return
		}
		err = json.Unmarshal(body, p)		
		if err != nil {
			fmt.Println("332")
			http.Error(w, "Переданы неверные данные", 422)
			return
		}
		x := p.X
		y := p.Y
		if !isInside(x, y) {
			http.Error(w, "Переданы неверные данные(Выход за пределы поля)", 422)
			return
		}
		h.Referee.roundWithBot(x, y)
		js,_:= json.Marshal(h.Referee.RoundResponse)
		w.Write(js)
	default:
		http.Error(w, "POST", 405)
		return
	}
}
func (h *HttpRest)makeShotPlayer(w http.ResponseWriter, r *http.Request){
	//здесь првоерки на номер игрока
	switch r.Method {
	case "POST":
		if !h.Started {
			http.Error(w, "Игра не начата!", 409)
			return
		}
		p := new(pointValue)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Переданы неверные данные", 422)
			return
		}
		err = json.Unmarshal(body, p)
		if err != nil {
			fmt.Println("332")
			http.Error(w, "Переданы неверные данные", 422)
			return
		}
		x := p.X
		y := p.Y
		if !isInside(x, y) {
			http.Error(w, "Переданы неверные данные(Выход за пределы поля)", 422)
			return
		}
		h.Referee.NumberPlayer = p.Value
		h.Referee.userRound(x, y)
		//здесь запись




		js,_:= json.Marshal(h.Referee.RoundResponse)
		w.Write(js)
	default:
		http.Error(w, "POST", 405)
		return
	}
}
