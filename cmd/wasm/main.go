package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"syscall/js"
	"time"
)

type match struct {
	players           []string
	c                 http.Client
	playerOutputValue js.Value
	playerInputValue  js.Value
}

func (m *match) addPlayer() js.Func {
	addPlayersFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		// player add logic
		playerName := args[0].String()
		fmt.Println("adding player", playerName)
		m.appendPlayer(playerName)
		players := m.beautify()
		m.playerOutputValue.Set("value", players)
		return players
	})
	return addPlayersFunc
}

func (m *match) beautify() string {
	var players = ""
	for _, v := range m.players {
		players += v + "\n"
	}
	return players
}

func (m *match) appendPlayer(name string) {
	u, err := url.Parse(fmt.Sprintf("https://cricket-match.herokuapp.com/getplayer?name=%v", name))
	if err != nil {
		fmt.Println("error parsing url")
	}
	r := http.Request{
		Method: "GET",
		URL:    u,
	}
	r.Header.Add("Access-Control-Allow-Origin", "*")
	res, err := m.c.Do(&r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error while getting player", err)
	}
	json.Unmarshal(b, &m.players)
}

func (m *match) getPlayers() {
	_, err := url.Parse("https://cricket-match.herokuapp.com/getplayer")
	if err != nil {
		fmt.Println("error parsing url")
	}
	r, err := http.NewRequest(http.MethodGet, "https://cricket-match.herokuapp.com/getplayer", nil)
	if err != nil {
		fmt.Println("error creating new request", err)
	}
	// r := http.Request{
	// 	Method: "GET",
	// 	URL:    u,
	// }
	r.Header.Add("Access-Control-Allow-Origin", "*")
	res, err := m.c.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error while getting player", err)
	}
	json.Unmarshal(b, &m.players)
}

func initMatch() *match {
	var err error
	m := match{players: make([]string, 0)}
	m.c = http.Client{
		Timeout: 15 * time.Second,
	}
	if err != nil {
		fmt.Println("error opening a file", err)
	}
	// get do elements
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		_ = map[string]interface{}{
			"error": "Unable to get document object",
		}
		// return result
		return nil
	}
	playerOutputArea := jsDoc.Call("getElementById", "nameoutput")
	if !playerOutputArea.Truthy() {
		_ = map[string]interface{}{
			"error": "Unable to get output text area",
		}
		return nil
	}
	playerInputArea := jsDoc.Call("getElementById", "nameinput")
	if !playerInputArea.Truthy() {
		result := map[string]interface{}{
			"error": "Unable to get output text area",
		}
		fmt.Println(result)
	}
	m.playerInputValue = playerInputArea
	m.playerOutputValue = playerOutputArea
	return &m
}

func main() {
	m := initMatch()
	m.getPlayers()
	fmt.Println("Go Web Assembly")
	js.Global().Set("addPlayer", m.addPlayer())
	<-make(chan bool)
}
