package world

import "goland/server/src/presence"

type Input interface {
	PlayerId() string
}

type World interface {
	Initialize() World
	Halt()
	Join(playerId string)
	Leave(playerId string)
	Sync(player *presence.User)
	Input(dt float32, i Input)
	Update(dt float32)
}

func New() World {
	return (&Prakriti{}).Initialize()
}
