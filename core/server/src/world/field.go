package world

import (
	"github.com/EngoEngine/ecs"
	"goland/protocol/gen/go/command"
	"goland/server/src/presence"
	"goland/server/src/world/component"
	"goland/server/src/world/system"
	"math"
	"time"
)

type FieldInput struct {
	player    string
	direction component.V2
	movement  component.V2
}

func (f FieldInput) PlayerId() string {
	return f.player
}

type Field struct {
	world   *ecs.World
	tiles   *Map
	players *system.PlayerSystem
}

type FieldPlayer struct {
	ecs.BasicEntity
	component.Transform
	component.Player
}

func ConvertToPI(playerId string, c *command.HIDCommand) *FieldInput {
	// Bounds check -1 to 1 for vectors

	return &FieldInput{
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

func (f *Field) Initialize() World {
	f.world = &ecs.World{}
	f.tiles = NewBorderedMap(80, 55)
	f.players = &system.PlayerSystem{}

	var playerable *system.Playerable
	f.world.AddSystemInterface(f.players, playerable, nil)

	return f
}

func (f *Field) Join(playerId string) {
	if exists := f.players.Get(playerId); exists == nil {
		// If it doesn't exist, create a player entity
		player := &FieldPlayer{
			BasicEntity: ecs.NewBasic(),
			Transform: component.Transform{
				X:      float32(f.tiles.Width / 2),
				Y:      float32(f.tiles.Height / 2),
				Width:  1,
				Height: 1,
			},
			Player: component.Player{
				PlayerId: playerId,
				LastSync: time.Now(),
			},
		}
		f.world.AddEntity(player)
	}
}

func (pr *Field) Leave(playerId string) {

}

func (pr *Field) Input(dt float32, i Input) {
	if input := i.(*FieldInput); input != nil {
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

func (f *Field) Update(dt float32) {
	f.world.Update(dt)
}

func (f *Field) Sync(pp *presence.User) {
	now := time.Now()

	if pe := f.players.Get(pp.Id); pe != nil {
		p := pe.GetPlayer()

		pp.Outbox <- f.players.Sync(pe)
		p.LastSync = now
	}
	// Spectator / admin?
}

func (f *Field) Halt() {

}
