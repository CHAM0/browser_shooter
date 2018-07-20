package main

import (
	"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

//Player regroupe pes attributs du joueur
type Player struct {
	ID      string
	Clock   time.Time
	Session sockjs.Session
	X, Y    int
	Keys    Keyword
	Bullets []Bullet
}

//NewPlayer créer un nouveau joueur à partir d'une session
func NewPlayer(session sockjs.Session) Player {
	return Player{ID: session.ID(), Session: session}
}

//Move deplace le joueur
func (p *Player) Move() {
	p.X += (-p.Keys.Q + p.Keys.D) * 10
	p.Y += (-p.Keys.Z + p.Keys.S) * 10
}

//Draw dessine le joueur
func (p *Player) Draw(ctx *js.Object, img *js.Object) {
	ctx.Call("drawImage", img, p.X, p.Y)
}

//Fire joueur tire
func (p *Player) Fire() {
	if p.Keys.Space == 1 {
		if len(p.Bullets) > 0 {
			if p.Bullets[len(p.Bullets)-1].cpt < 65 {
				p.Bullets = append(p.Bullets, NewBullet(p))
			}
		} else {
			p.Bullets = append(p.Bullets, NewBullet(p))
		}
	}
	for i := range p.Bullets {
		p.Bullets[i].Update()
		if p.Bullets[i].cpt <= 0 {
			p.Bullets = append(p.Bullets[:i], p.Bullets[i+1:]...)
			break
		}
	}
}

//Bullet balle
type Bullet struct {
	Angle     float64
	X, Y      float64
	XVelocity float64
	YVelocity float64
	cpt       int
}

//NewBullet créer une balle
func NewBullet(p *Player) Bullet {

	angle := math.Atan2(p.Keys.MouseX-float64(p.X), p.Keys.MouseY-float64(p.Y))
	//angle = math.Atan(angle)
	xVelocity := math.Sin(angle)
	yVelocity := math.Cos(angle)

	return Bullet{
		Angle:     angle,
		X:         float64(p.X),
		Y:         float64(p.Y),
		XVelocity: xVelocity,
		YVelocity: yVelocity,
		cpt:       70,
	}
}

//Update deplace la balle
func (b *Bullet) Update() {
	b.X += b.XVelocity * 15
	b.Y += b.YVelocity * 15
	b.cpt--
}

//Draw dessine la balle
func (b *Bullet) Draw(ctx *js.Object, img *js.Object) {
	ctx.Call("drawImage", img, b.X, b.Y)
}
