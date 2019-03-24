package main

import (
	"github.com/gorilla/sessions"
)

// GetClient - return the signed in user stored in a session cookie
func GetClient(s *sessions.Session) Client {
	val := s.Values["client"]
	var client = Client{}
	client, ok := val.(Client)

	if !ok {
		return Client{Authenticated: false}
	}
	return client
}
