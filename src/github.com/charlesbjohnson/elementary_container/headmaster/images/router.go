package images

import (
	"github.com/charlesbjohnson/elementary_container/headmaster"
)

func Register(server headmaster.Server) {
	application := server.Application()

	router := application.Router.PathPrefix("/images").Subrouter()

	router.Path("/").Methods("GET").HandlerFunc(startHandler)
}
