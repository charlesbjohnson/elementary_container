package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/charlesbjohnson/elementary_container/fscapture"
	"github.com/charlesbjohnson/elementary_container/fscapture/compress"
	"os"
	"runtime"
	"strings"
	"sync"
)

const RootDirectory = "/home/vagrant/test"

var ExcludePaths = []string{
	"/dev",
	"/proc",
	"/sys",
	"/tmp",
	"/run",
	"/mnt",
	"/media",
	"/lost+found",
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	image := compress.New(RootDirectory)

	if err := createImage(image); err != nil {
		log.Fatal(err)
	}
}

func createImage(capturable fscapture.Capturable) error {
	capturable.RegisterFileCaptureHook("excludedPaths", pathExcluded)

	var wait sync.WaitGroup
	wait.Add(1)
	defer wait.Wait()

	go func() {
		defer wait.Done()

		for event := range capturable.FileCaptureEvents() {
			logger := log.WithField("file", event.Path)

			if event.Captured {
				logger.Info(event.Message)
			} else {
				log.Warn(event.Message)
			}
		}
	}()

	if err := capturable.Capture("/tmp/advisor/image"); err != nil {
		return err
	}

	log.WithField("file", capturable.Path()).Info("image capture succeeded")
	return nil
}

func pathExcluded(path string, info os.FileInfo) bool {
	for _, excludedPath := range ExcludePaths {
		if strings.HasPrefix(path, excludedPath) {
			return false
		}
	}

	return true
}
