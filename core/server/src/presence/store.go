package presence

import "time"

type PresenceStore interface {
	Connect(user *User) (string, *User, error)
	Disconnect(user *User) (string, *User, error)
}

type Store struct {
	table map[string]*User
}

var LocalStore = &Store{
	table: make(map[string]*User),
}

func (ss *Store) ById(id string) *User {
	return ss.table[id]
}

func (ss *Store) Connect(user *User) (string, *User, error) {
	// Last one in wins
	ss.table[user.Id] = user
	return user.Id, user, nil
}

func (ss *Store) Heartbeat(user *User)  {
	if s, ok := ss.table[user.Id]; ok {
		s.Last = time.Now().Second()
	}
}

func (ss *Store) Disconnect(user *User) {
	ss.table[user.Id] = nil
}
