package headmaster

import (
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
)

type Application struct {
	Log        *logrus.Logger
	Middleware *negroni.Negroni
	Router     *mux.Router
	View       *View
}

type Server interface {
	Application() *Application
	Run()
}

func New(log *logrus.Logger) (*Application, error) {
	application := &Application{}

	_, currentFile, _, _ := runtime.Caller(0)
	directory := path.Join(path.Dir(currentFile), "views")

	view, err := NewView(directory)
	if err != nil {
		return nil, err
	}

	middleware := negroni.New()
	middleware.Use(negroni.NewRecovery())

	logrusMiddleware := negronilogrus.NewMiddleware()
	logrusMiddleware.Logger.Out = os.Stdout
	middleware.Use(logrusMiddleware)

	middleware.UseHandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		context.Set(request, "server", application)
	})

	router := mux.NewRouter().StrictSlash(true)
	middleware.UseHandler(router)

	application.Log = logrusMiddleware.Logger
	application.Middleware = middleware
	application.Router = router
	application.View = view

	return application, nil
}

func (application *Application) Application() *Application {
	return application
}

func (application *Application) Run() {
	host := os.Getenv("SERVER_URL_HOST")
	application.Log.WithField("host", host).Info("starting server")
	http.ListenAndServe(host, application.Middleware)
}
