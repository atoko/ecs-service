package session

import (
	"errors"
	"github.com/google/uuid"
	"goland/server/src/world"
	"sync"
)

type Store struct {
	table  map[string]*State
	active sync.Map
}

var LocalStore = &Store{
	table:  make(map[string]*State),
	active: sync.Map{},
}

func (ss *Store) Create() (string, *State, error) {
	id := uuid.New().String()
	state := &State{
		players: make([]string, 4),
		world: world.New(),
	}

	ss.table[id] = state
	return id, state, nil
}

func (ss *Store) Get(id string) (*State, error) {
	if s, ok := ss.table[id]; ok {
		if _, loaded := ss.active.LoadOrStore(id, true); loaded == false {
			s.Load(id)
		}
		return s, nil
	} else {
		return nil, errors.New("SESSION_NOT_FOUND")
	}
}
