package advisor

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/charlesbjohnson/elementary_container/fscapture/compress"
)

const (
	inputPath    = "/home/vagrant/test"
	outputPath   = "/tmp/advisor/image"
	artifactPath = "/tmp/advisor/artifact"
)

var excludePaths = []string{
	"/dev",
	"/proc",
	"/sys",
	"/tmp",
	"/run",
	"/mnt",
	"/media",
	"/lost+found",
}

type Application struct {
	Log *logrus.Logger
}

func New(log *logrus.Logger) *Application {
	return &Application{Log: log}
}

func (application *Application) Run() {
	image := compress.New(inputPath, outputPath)

	image.RegisterFileCaptureHook("ExcludePaths", isPathIncluded)

	if err := application.capture(image); err != nil {
		application.Log.WithField("file", image.Path()).Fatal(err)
	}

	downloader, err := application.upload(os.Getenv("SERVER_URL"), image.File())
	if err != nil {
		application.Log.Fatal(err)
	}

	// TODO long poll the server to find out when to download the thing

	if err := application.download(downloader, artifactPath); err != nil {
		application.Log.Fatal(err)
	}

	image.Close()
}

func isPathIncluded(path string, info os.FileInfo) bool {
	result := true

	for _, excludedPath := range excludePaths {
		if strings.HasPrefix(path, excludedPath) {
			result = false
			break
		}
	}

	return result
}
