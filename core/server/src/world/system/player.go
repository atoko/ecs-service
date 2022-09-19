package system

import (
	"github.com/EngoEngine/ecs"
	"goland/server/src/world/component"
)

type playerEntity struct {
	*ecs.BasicEntity
	*component.Player
	*component.Transform
}

type Playerable interface {
	ecs.BasicFace
	component.PlayerFace
	component.TransformFace
}

type PlayerSystem struct {
	entities map[string]playerEntity
}

func (ps *PlayerSystem) New(world *ecs.World)  {
	ps.entities = map[string]playerEntity{}
}

func (ps *PlayerSystem) Get(id string) Playerable {
	if p, ok := ps.entities[id]; ok {
		return p
	}

	return nil
}

func (ps *PlayerSystem) AddByInterface(o ecs.Identifier) {
	obj := o.(Playerable)
	ps.Add(obj.GetBasicEntity(),
		obj.GetPlayer(),
		obj.GetTransform(),
	)
}

func (ps *PlayerSystem) Add(
	basic *ecs.BasicEntity,
	player *component.Player,
	space *component.Transform,
) {

	entity := playerEntity{
		basic,
		player,
		space,
	}

	ps.entities[player.PlayerId] = entity
}

func (ps *PlayerSystem) Remove(basic ecs.BasicEntity) {
	index := ""
	for i, entity := range ps.entities {
		if entity.ID() == basic.ID() {
			index = i
		}
	}

	if index != "" {
		delete(ps.entities, index)
	}
}

func (ps *PlayerSystem) Update(dt float32) {
	for _, entity := range ps.entities {
		if entity.MoveElapsed < 90.0 {
			entity.MoveElapsed += dt
		}
	}
}

type PlayerSystemSync struct {
	PlayerTransforms map[string]component.V2
}

func (ps *PlayerSystem) Sync(pe Playerable) PlayerSystemSync {
	message := PlayerSystemSync{PlayerTransforms: map[string]component.V2{}}
	for id, entity := range ps.entities {
		t := entity.Transform
		// calculate viewshed
		message.PlayerTransforms[id] = component.V2{
			X: t.X,
			Y: t.Y,
		}
	}

	return message
}