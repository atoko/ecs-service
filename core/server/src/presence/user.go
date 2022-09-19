package presence

type User struct {
	Id     string
	Last   int
	Outbox chan interface{}
}
