package seabattle

import (
	"bytes"
	"encoding/gob"
	"os"

)

const (
	USER2_WIN   = 3
	USER1_WIN  = 2
	USER1_SHOT = 1
	USER2_SHOT  = 4
	USER2_KILL = 5
	USER1_KILL = 6
)

type roundResponse struct {
	State      int          `json:"state"`
	User1Results []pointValue   `json:userResult`
	User2Results []pointValue `json:"botResults"`
}

func (r *roundResponse) init() {
	r.State = -1
	r.User1Results = make([]pointValue, 0)
	r.User2Results = make([]pointValue, 0)
}

type referee struct {
	State         int
	Field1     field
	Field2      field
	RoundResponse roundResponse
	TypeEnemy bool //false - bot, true-player
	NumberPlayer int
}

func (r *referee) init() {
	r.Field1.initField()
	r.Field2.initField()
	r.Field2.Bot = new(bot)
	r.Field2.Bot.initBot()
}
func (r *referee) saveGame() error {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	enc.Encode(&r)

	file, err := os.Create("GameSave.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(buffer.Bytes())
	return nil
}
func (r *referee) loadGame() error {

	r.Field1 = field{}
	r.Field2 = field{}

	file, err := os.Open("GameSave.txt")

	if err != nil {
		return err
	}

	dec := gob.NewDecoder(file)
	err = dec.Decode(r)

	if err != nil {
		return err
	}
	return nil
}
func (r *referee) userRound(x, y int) {
	f:= &r.Field2
	resp:= &r.RoundResponse.User1Results
	if r.TypeEnemy{
		if r.NumberPlayer==2{
			f = &r.Field1
			resp = &r.RoundResponse.User2Results
		}
	}
	result := f.getShot(x, y)
	*resp = append(*resp, pointValue{x,y,result})
	//r.RoundResponse.User1Result = pointValue{x,y,result}
	if result== KILLED && !f.hasPoints(){
		r.State = USER1_WIN
		return
	}
	if result == SHIP || result == HIT || result == MISS || result==KILLED{		
		r.State = USER1_SHOT
	} else {
		r.State = USER2_SHOT
	}
}
func (r *referee) botRound() {
	x, y := r.Field2.Bot.makeShot()
	result := r.Field1.getShot(x, y)
	r.Field2.Bot.setResult(result)
	r.RoundResponse.User2Results = append(r.RoundResponse.User2Results,pointValue{x,y,result})
	if result== KILLED && !r.Field1.hasPoints(){
		r.State= USER2_WIN
		return
	}
	if result == SHIP || result == HIT || result == MISS || result==KILLED{
	
		r.State = USER2_SHOT
	} else {
		r.State = USER1_SHOT
		return
	}
	r.botRound()
}
func (r *referee) roundWithBot(x, y int) {
	r.RoundResponse.init()
	defer r.saveGame()
	r.userRound(x, y)
	r.RoundResponse.State = r.State
	if r.State == USER1_SHOT || r.State == USER1_WIN || r.State == USER1_KILL{
		return
	}
	r.botRound()
	r.RoundResponse.State = r.State
}
