package presence

import "goland/server/src/config"

var HeartbeatRecieve = func(user *User) {
	config.StaticLoggers.Info.Printf("%s sent heartbeat", user.Id)
	user.Outbox <- "!"
	LocalStore.Heartbeat(user)
}
