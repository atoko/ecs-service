package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"goland/server/src/config"
	"goland/server/src/controller"
	"goland/server/src/server"
	"goland/server/src/session"
	"net/http"
	"time"
)

func main() {
	PORT := 9999
	loggers := config.StaticLoggers

	mux := mux.NewRouter()
	mux.HandleFunc("/v1/session", session.CreateHandler)
	mux.Handle("/v1/session/{sessionId}/ws", server.WebsocketHandler)
	withAuth := &controller.AuthHeaderParser{
		Handler: mux,
	}
	withContext := &controller.GolandContextMiddleware{
		Handler: withAuth,
		Context: &controller.GolandContext{
			Log: loggers,
		},
	}
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        withContext,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    180 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	loggers.System.Printf("Initializing server on port %d", PORT)
	loggers.System.Fatal(s.ListenAndServe())
}
