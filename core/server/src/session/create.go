package session

import (
	"encoding/json"
	"goland/server/src/controller"
	"net/http"
)
import dto "goland/protocol/gen/go/session"

var create = func(w http.ResponseWriter, r *http.Request) {
	loggers := controller.GetGolandContext(r.Context()).Log
	id, _, err := LocalStore.Create()
	if err != nil {
		loggers.Error.Printf("Error while creating session")
		return
	}
	response, _ := json.Marshal(dto.CreateResponse{
		Data: &dto.CreateResponse_Session{
			Id: id,
		},
	})
	w.Write(response)
}

var CreateHandler = controller.Post(create)
