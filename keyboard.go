package main

//Keyword regroupe les commandes du jeux
type Keyword struct {
	Z, Q, S, D, Space int
	MouseX, MouseY    float64
}

//OnKeyDown retourne les commandes préssé
func (key *Keyword) OnKeyDown(keyCode int) {
	if keyCode == 90 {
		key.Z = 1
	} else if keyCode == 81 {
		key.Q = 1
	} else if keyCode == 83 {
		key.S = 1
	} else if keyCode == 68 {
		key.D = 1
	} else if keyCode == 32 {
		key.Space = 1
	}
}

//OnKeyUp retourne les commandes relaché
func (key *Keyword) OnKeyUp(keyCode int) {
	if keyCode == 90 {
		key.Z = 0
	} else if keyCode == 81 {
		key.Q = 0
	} else if keyCode == 83 {
		key.S = 0
	} else if keyCode == 68 {
		key.D = 0
	} else if keyCode == 32 {
		key.Space = 0
	}
}
