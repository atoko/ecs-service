package component

import "time"

type Player struct {
	PlayerId string
	MoveElapsed float32
	LastSync time.Time
}

func (t *Player) GetPlayer() *Player {
	return t
}

type PlayerFace interface {
	GetPlayer() *Player
}