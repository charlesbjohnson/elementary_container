package images

import (
	"net/http"
	"net/url"
	"os"

	"github.com/charlesbjohnson/elementary_container/headmaster"
	"github.com/gorilla/context"
)

func startHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	clientUrl := (&url.URL{
		Scheme: os.Getenv("CLIENT_URL_PROTOCOL"),
		Host:   os.Getenv("CLIENT_URL_HOST"),
		Path:   os.Getenv("CLIENT_URL_PATH"),
	}).String()

	serverUrl := (&url.URL{
		Scheme: os.Getenv("SERVER_URL_PROTOCOL"),
		Host:   os.Getenv("SERVER_URL_HOST"),
	}).String()

	context := struct {
		ClientUrl string
		ServerUrl string
	}{clientUrl, serverUrl}

	application.View.Plain(response, http.StatusOK, "execute.sh.tmpl", context)
}
