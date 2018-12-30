package core

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user/auth/local"
	"net/http"
)

type Extension struct {
	nibbler.NoOpExtension
	AuthExtension *local.Extension
}

func (s *Extension) AddRoutes(context *nibbler.Application) error {
	context.GetRouter().HandleFunc("/api/ok", s.AuthExtension.EnforceLoggedIn(func (w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "OK"}`))
	})).Methods("GET")
	return nil
}
