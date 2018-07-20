package main

import (
	"encoding/json"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

//Message regroupe les données à envoyer au serveur
type Message struct {
	Command    string
	Clock      time.Time
	DataString string
	DataInt    int
	Game       Game
	Keys       Keyword
}

//Send envoi un message
func (msg Message) Send(any interface{}) {
	m, _ := json.Marshal(msg)
	switch any.(type) {
	case sockjs.Session:
		session := any.(sockjs.Session)
		session.Send(string(m))
	case *js.Object:
		socket := any.(*js.Object)
		socket.Call("send", string(m))
	}
}

//Broadcast un message à tous les joueurs
func (msg Message) Broadcast() {
	for i := range msg.Game.Players {
		msg.Send(msg.Game.Players[i].Session)
	}
}
