package world

import (
	"goland/protocol/gen/go/command"
	"goland/server/src/presence"
	"github.com/EngoEngine/ecs"
	"goland/server/src/world/component"
	"goland/server/src/world/system"
	"math"
	"time"
)

type PrakritiInput struct {
	player string
	direction component.V2
	movement component.V2
}

func (p PrakritiInput) PlayerId() string {
	return p.player
}

type Prakriti struct {
	world *ecs.World
	tiles *Map
	players *system.PlayerSystem
}

type PrakritiPlayer struct {
	ecs.BasicEntity
	component.Transform
	component.Player
}

func ConvertToPI(playerId string, c *command.PresenceCommand) *PrakritiInput {
	// Bounds check -1 to 1 for vectors

	return &PrakritiInput{
		player: playerId,
		direction: component.V2{
			X: c.Direction.X,
			Y: c.Direction.Y,
		},
		movement: component.V2{
			X: float32(math.Min(1, math.Max(-1, float64(c.Movement.X)))),
			Y: float32(math.Min(1, math.Max(-1, float64(c.Movement.Y)))),
		},
	}
}

func (pr *Prakriti) Initialize() World {
	pr.world = &ecs.World{}
	pr.tiles = NewBorderedMap(80, 55)
	pr.players = &system.PlayerSystem{}

	var playerable *system.Playerable
	pr.world.AddSystemInterface(pr.players, playerable, nil)

	return pr
}


func (pr *Prakriti) Join(playerId string) {
	if exists := pr.players.Get(playerId); exists == nil {
		// If it doesn't exist, create a player entity
		player := &PrakritiPlayer{
			BasicEntity: ecs.NewBasic(),
			Transform:  component.Transform{
				X: float32(pr.tiles.Width / 2),
				Y: float32(pr.tiles.Height / 2),
				Width:  1,
				Height: 1,
			},
			Player: component.Player{
				PlayerId: playerId,
				LastSync: time.Now(),
			},
		}
		pr.world.AddEntity(player)
	}
}

func (pr *Prakriti) Leave(playerId string) {

}

func (pr *Prakriti) Input(dt float32, i Input) {
	if input := i.(*PrakritiInput); input != nil {
		m := pr.tiles

		if pe := pr.players.Get(i.PlayerId()); pe != nil {
			p := pe.GetPlayer()
			t := pe.GetTransform()

			if p.MoveElapsed > 0.033 {
				return
			}

			if input.movement.X == 0 && input.movement.Y == 0 {
				return
			}

			dx := t.X + ((input.movement.X * dt) * 10000)
			dy := t.Y + ((input.movement.Y * dt) * 10000)

			destination := pr.tiles.xy_idx(
				int(dx),
	 			int(dy),
			)

			if tile := pr.tiles.Tiles[destination]; tile != WALL {
				t.X = float32(math.Min(float64(m.Width-1), math.Max(0, float64(dx))))
				t.Y = float32(math.Min(float64(m.Height-1), math.Max(0, float64(dy))))
			}

			p.MoveElapsed = 0
		}
	}
}

func (pr *Prakriti) Update(dt float32) {
	pr.world.Update(dt)
}

func (pr *Prakriti) Sync(pp *presence.User) {
	now := time.Now()

	if pe := pr.players.Get(pp.Id); pe != nil {
		p := pe.GetPlayer()

		pp.Outbox <- pr.players.Sync(pe)
		p.LastSync = now
	}
	// Spectator / admin?
}

func (pr *Prakriti) Halt() {

}

