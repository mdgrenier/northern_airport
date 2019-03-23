package main

import (
	"github.com/gorilla/sessions"
)

func getClient(s *sessions.Session) Client {
	val := s.Values["client"]
	var client = Client{}
	client, ok := val.(Client)

	if !ok {
		return Client{Authenticated: false}
	}
	return client
}
