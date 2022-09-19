package session

import (
	"goland/server/src/config"
	"goland/server/src/presence"
	"goland/server/src/world"
	"reflect"
	"time"
)
import "goland/protocol/gen/go/command"

type State struct {
	Inbox    chan interface{}
	commands chan interface{}
	closed   chan bool
	players  []string
	world    world.World
	lastTime time.Time
}

type PresenceCommand struct {
	Id      string
	Command *command.HIDCommand
}

type JoinCommand struct {
	id string
}

type LeaveCommand struct {
	id string
}

func (s *State) Join(id string) {
	for i := range s.players {
		sp := s.players[i]
		if sp == id {
			return
		}
	}

	s.commands <- JoinCommand{id: id}
}

func (s *State) Leave(id string) {
	index := -1
	for i := range s.players {
		sp := s.players[i]
		if sp == id {
			index = i
			break
		}
	}

	if index >= 0 {
		s.commands <- LeaveCommand{id: id}
	}
}

var timestep = time.Millisecond * 33 // 30 hz
func (s *State) Load(selfId string) {
	if s.closed != nil {
		s.closed <- true
		close(s.Inbox)
		close(s.commands)
	}

	s.closed = make(chan bool)
	s.Inbox = make(chan interface{}, 8)
	s.commands = make(chan interface{}, 16)
	s.lastTime = time.Now()

	go func() {
		for v := range s.Inbox {
			config.StaticLoggers.Info.Printf(
				"%s: recieved %s, typeof %s",
				selfId,
				v,
				reflect.TypeOf(v),
			)
			s.commands <- v
		}
	}()

	go func() {
		for {
			dt := float32(time.Now().Sub(s.lastTime).Seconds())
			select {
			case input, valid := <-s.commands:
				if valid {
					switch input := input.(type) {
					case JoinCommand:
						for i := range s.players {
							sp := s.players[i]
							if sp == input.id {
								return
							}
						}

						s.players = append(s.players, input.id)
						s.world.Join(input.id)
					case LeaveCommand:
						index := -1
						for i := range s.players {
							sp := s.players[i]
							if sp == input.id {
								index = i
								break
							}
						}

						if index >= 0 {
							s.players[index] = s.players[len(s.players)-1]
							s.players = s.players[:len(s.players)-1]
						}
						s.world.Leave(input.id)
					case PresenceCommand:
						s.world.Input(dt, world.ConvertToPI(
							input.Id,
							input.Command,
						))
					default:
						config.StaticLoggers.Info.Printf("%s: Default case: Value %v", selfId)
					}
				}
			default:
			}

			select {
			case _ = <-s.closed:
				goto exit
			case _ = <-time.After(timestep):
				s.world.Update(dt)
				s.lastTime = time.Now()
			}

			for i := range s.players {
				id := s.players[i]
				p := presence.LocalStore.ById(id)

				if p == nil {
					continue
				}

				s.world.Sync(p)
			}
		}
	exit:
		s.world.Halt()
		// Send players message
		config.StaticLoggers.Info.Printf("Exiting state loop")
	}()
}

func (s *State) Unload() {
	close(s.Inbox)
	close(s.commands)
	s.closed <- true
	s.closed = nil
}
