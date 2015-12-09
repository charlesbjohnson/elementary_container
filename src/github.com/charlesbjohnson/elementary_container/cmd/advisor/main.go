package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/charlesbjohnson/elementary_container/advisor"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stdout

	application := advisor.New(logger)
	application.Run()
}
