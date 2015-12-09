package images

import (
	"github.com/charlesbjohnson/elementary_container/headmaster"
)

func Register(server headmaster.Server) {
	application := server.Application()

	router := application.Router.PathPrefix("/images").Subrouter()

	router.Path("/").Methods("GET").HandlerFunc(startHandler)
	router.Path("/").Methods("POST").HandlerFunc(createHandler)
	router.Path("/commit/").Methods("POST").HandlerFunc(commitHandler)
	router.Path("/add/").Methods("GET").HandlerFunc(addHandler)
	router.Path("/poll/").Methods("GET").HandlerFunc(pollHandler)
}
