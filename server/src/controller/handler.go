package controller

import (
	"net/http"
)

var Post = func(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			f(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
