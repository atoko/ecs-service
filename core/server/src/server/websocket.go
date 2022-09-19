package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"goland/server/src/controller"
	"goland/server/src/presence"
	"goland/server/src/session"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"time"
)
import "goland/protocol/gen/go/command"

const REQUEST_TIMEOUT = 180 * time.Second

var WebsocketHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	loggers := controller.GetGolandContext(r.Context()).Log
	principal := controller.AuthPrincipalFromContext(r.Context())
	if principal == nil {
		invalidAuthResponse, _ := json.Marshal(map[string]map[string]string{"error": {"code": "AUTH_INVALID"}})
		w.WriteHeader(401)
		w.Write(invalidAuthResponse)
		return
	}

	vars := mux.Vars(r)
	sessionId, ok := vars["sessionId"]
	if !ok {
		sessionNotFound, _ := json.Marshal(map[string]map[string]string{"error": {"code": "SESSION_ID_REQUIRED"}})
		w.WriteHeader(404)
		w.Write(sessionNotFound)
		return
	}
	loggers.Warn.Printf("Connection from %s to session: %s", principal.ProfileId, sessionId)

	state, err := session.LocalStore.Get(sessionId)
	if err != nil {
		sessionNotFound, _ := json.Marshal(map[string]map[string]string{"error": {"code": "SESSION_NOT_FOUND"}})
		w.WriteHeader(404)
		w.Write(sessionNotFound)
		return
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: true,
	})
	if err != nil {
		loggers.Error.Printf("Unable to upgrade connection to websocket: %s", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	user := &presence.User{
		Id:     principal.ProfileId,
		Outbox: make(chan interface{}, 128),
	}

	presence.LocalStore.Connect(user)
	state.Join(principal.ProfileId)

	ctx, cancel := context.WithTimeout(r.Context(), REQUEST_TIMEOUT)
	defer cancel()

	go func() {
		for v := range user.Outbox {
			err := wsjson.Write(r.Context(), c, v)
			if err != nil {
				go func() {
					for _ = range user.Outbox {
					}
				}()
				return
			}
		}
	}()

	reason := "TERMINATED"
	for {
		ctx, cancel = context.WithTimeout(r.Context(), REQUEST_TIMEOUT)
		var v command.HIDCommand
		err = wsjson.Read(ctx, c, &v)

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway ||
			websocket.CloseStatus(err) == websocket.StatusInvalidFramePayloadData {
			reason = string(websocket.CloseStatus(err))
			break
		}

		if errors.Is(err, context.Canceled) {
			reason = "CANCELED"
			break
		}
		if errors.Is(err, context.DeadlineExceeded) {
			reason = "TIMEOUT"
			cancel()
			break
		}
		if err != nil {
			loggers.Error.Printf("Error reading message: %s", err)

			if strings.Contains(strings.ToLower(err.Error()), "closed") {
				reason = "CLOSED"
				cancel()
				break
			}
		} else {
			if v.Elapsed == 0 {
				presence.HeartbeatRecieve(user)
			} else {
				state.Inbox <- session.PresenceCommand{
					Id:      user.Id,
					Command: &v,
				}
			}
		}
		cancel()
	}

	presence.LocalStore.Disconnect(user)

	cancel()
	if reason != "TERMINATED" {
		loggers.Warn.Printf("Connection closed: %s", reason)
	}
	err = c.Close(websocket.StatusNormalClosure, reason)
	if err != nil {
		return
	}
})
