package main

import (
	"time"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

//Game regroupe la liste des joueurs
type Game struct {
	ID      string
	Clock   time.Time
	Players []Player
}

//NewGame nouvelle instance de jeux
func NewGame(id string) Game {
	return Game{ID: id, Clock: time.Now()}
}

//AddPlayer ajoute un joueur
func (game *Game) AddPlayer(session sockjs.Session) {
	game.Players = append(game.Players, NewPlayer(session))
	msg := Message{Command: "connect", DataString: session.ID(), Game: *game}
	msg.Send(session)
}

//RemovePlayer supprime un joueur
func (game *Game) RemovePlayer(session sockjs.Session) {
	for i := range game.Players {
		if game.Players[i].ID == session.ID() {
			game.Players = append(game.Players[:i], game.Players[i+1:]...)
			break
		}
	}
}

//Play met à jour l'état du jeux
func (game *Game) Play() {
	for i := range game.Players {
		game.Players[i].Move()
		game.Players[i].Fire()
	}
}

//Sync synchronise les joueurs
func (game Game) Sync() {
	msg := Message{Clock: time.Now(), Game: game}
	msg.Broadcast()
}

//Sleep Interval entre chaque mise à jour de l'état du jeux
func (game *Game) Sleep() {
	dt := time.Second/40 - time.Now().Sub(game.Clock)
	time.Sleep(dt)
	game.Clock = time.Now()
}
