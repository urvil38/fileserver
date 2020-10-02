package main

import (
	"crypto/subtle"
	"net/http"
)

type auth struct {
	username string
	password string
	relm     string
}

func (a *auth) basicAuthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(a.username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(a.password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+a.relm+`"`)
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
