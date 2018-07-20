// +build !js

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

const port string = "pc:8081"

var game = NewGame("Game1")
var messages = make(map[string]Message)

//SockJsHandler !
func SockJsHandler(session sockjs.Session) {
	fmt.Println("Client connected : ", session.ID())
	game.AddPlayer(session)
	for {
		if m, err := session.Recv(); err == nil {
			msg := Message{}
			json.Unmarshal([]byte(m), &msg)
			lastMsg, present := messages[msg.DataString]
			if !present || lastMsg.Clock.Before(msg.Clock) {
				messages[msg.DataString] = msg
			}
			continue
		}
		break
	}
	fmt.Println("Client disconnected:", session.ID())
	game.RemovePlayer(session)
}

var templates = template.Must(template.ParseFiles("index.html"))

//HomeHandler !
func HomeHandler(w http.ResponseWriter, t *http.Request) {
	//currentTime := time.Now()
	//response := "Time : " + currentTime.String()
	//fmt.Fprintln(w, response, nil)
	//templates.ExecuteTemplate(w, "index.html", nil)
	templates.ExecuteTemplate(w, "index.html", nil)
}

//StaticFileHandler !
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {

	go func() {
		for {
			for i := range game.Players {
				if msg, ok := messages[game.Players[i].ID]; ok {
					game.Players[i].Clock = msg.Clock
					game.Players[i].Keys = msg.Keys
				}
			}
			game.Play()
			game.Sync()
			game.Sleep()
		}
	}()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/static/", StaticFileHandler)
	http.Handle("/socket/", sockjs.NewHandler("/socket", sockjs.DefaultOptions, SockJsHandler))

	fmt.Println("Starting server ...")
	http.ListenAndServe(port, nil)
}
