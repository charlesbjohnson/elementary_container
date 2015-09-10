package bootstrap

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/charlesbjohnson/elementary_container/headmaster"
	"github.com/charlesbjohnson/elementary_container/headmaster/images"
	"github.com/joho/godotenv"
)

func Run() {
	logger := logrus.New()
	logger.Out = os.Stdout

	if err := godotenv.Load(); err != nil {
		logger.WithField("file", ".env").Fatal(err)
	}

	server, err := headmaster.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	images.Register(server)
	server.Run()
}
