// +build js

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

var player *js.Object
var background *js.Object
var bullet *js.Object
var game = Game{}
var id string
var latency int
var latencyTimer time.Time
var keyboard = Keyword{}
var messages []Message

func main() {

	background = js.Global.Get("document").Call("createElement", "img")
	background.Set("src", "Sprites/background.jpg")
	player = js.Global.Get("document").Call("createElement", "img")
	player.Set("src", "Sprites/character.png")
	bullet = js.Global.Get("document").Call("createElement", "img")
	bullet.Set("src", "Sprites/bullet.png")

	host := js.Global.Get("location").Get("host").String()
	console.Log("host", host)
	sockjs := js.Global.Call("SockJS", "http://pc:8081/socket/")

	canvas := js.Global.Get("document").Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")
	ctx.Set("font", "10px Georgia")
	ctx.Set("fillStyle", "black")

	sockjs.Set("onopen", func() {
		console.Log("Connected to server.")
	})

	sockjs.Set("onmessage", func(e *js.Object) {

		msg := Message{}
		json.Unmarshal([]byte(e.Get("data").String()), &msg)
		if msg.Command == "connect" {
			id = msg.DataString
			game = msg.Game
			console.Log(fmt.Sprintf("Session ID %s", id))
		} else {
			game = msg.Game
			if time.Now().After(latencyTimer.Add(time.Second)) {
				latency = int(time.Now().Sub(msg.Clock) / time.Millisecond)
				latencyTimer = time.Now()
			}
		}
	})

	sockjs.Set("onclose", func() {
		console.Log("Disconnected from server.")
	})

	js.Global.Get("document").Set("onkeydown", func(e *js.Object) {
		e.Call("preventDefault")
		keyboard.OnKeyDown(e.Get("keyCode").Int())
		canvas.Set("onmousemove", func(e *js.Object) {
			e.Call("preventDefault")
			rect := canvas.Call("getBoundingClientRect")
			canvasX := canvas.Get("width").Float()
			canvasY := canvas.Get("height").Float()

			MouseX := e.Get("clientX").Float() - rect.Get("left").Float()
			scaleX := rect.Get("right").Float() - rect.Get("left").Float()
			posX := MouseX / scaleX * canvasX

			MouseY := e.Get("clientY").Float() - rect.Get("top").Float()
			scaleY := rect.Get("bottom").Float() - rect.Get("top").Float()
			posY := MouseY / scaleY * canvasY
			keyboard.MouseX = posX
			keyboard.MouseY = posY
			console.Log("mousse", keyboard.MouseX, keyboard.MouseY)

		})
	})

	js.Global.Get("document").Set("onkeyup", func(e *js.Object) {
		e.Call("preventDefault")
		keyboard.OnKeyUp(e.Get("keyCode").Int())
	})

	for {
		if game.ID == "" {
			console.Log("sleep")
			time.Sleep(time.Second)
			continue
		}
		for i := range game.Players {
			if game.Players[i].ID == id {
				leftover := make([]Message, len(messages))
				for j := range messages {
					if messages[j].Clock.After(game.Players[i].Clock) {
						game.Players[i].Move()
						leftover = append(leftover, messages[j])
					}
				}
				messages = leftover
				break
			}
		}
		msg := Message{Clock: time.Now(), DataString: id, Keys: keyboard}
		msg.Send(sockjs)
		messages = append(messages, msg)
		ctx.Call("clearRect", 0, 0, canvas.Get("width"), canvas.Get("height"))
		//ctx.Call("drawImage", background, 0, 0)
		for i := range game.Players {
			console.Log("player", game.Players[i].X, game.Players[i].Y)
			game.Players[i].Draw(ctx, player)
			for j := range game.Players[i].Bullets {
				game.Players[i].Bullets[j].Draw(ctx, bullet)
			}
		}
		ctx.Call("fillText", fmt.Sprintf("Latency %dms", latency), 950, 490)
		game.Sleep()
	}
}
