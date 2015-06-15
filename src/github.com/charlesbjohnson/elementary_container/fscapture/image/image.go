package image

import (
	"archive/tar"
	"github.com/charlesbjohnson/elementary_container/fscapture"
	"os"
	"sync"
)

type fileCaptureHookRegistration struct {
	Id   string
	Hook fscapture.FileCaptureHook
}

type Image struct {
	inputPath         string
	outputPath        string
	targetPath        string
	targetFile        *os.File
	writer            *tar.Writer
	handlers          []fileCaptureHookRegistration
	fileCaptureEvents chan fscapture.FileCaptureEvent
	wait              sync.WaitGroup
}

func New(inputPath string) *Image {
	return &Image{
		inputPath:         inputPath,
		fileCaptureEvents: make(chan fscapture.FileCaptureEvent),
	}
}

func (image *Image) RegisterFileCaptureHook(id string, hook fscapture.FileCaptureHook) fscapture.Capturable {
	image.handlers = append(image.handlers, fileCaptureHookRegistration{id, hook})
	return image
}

func (image *Image) FileCaptureEvents() <-chan fscapture.FileCaptureEvent {
	return image.fileCaptureEvents
}

func (image *Image) Capture(outputPath string) error {
	defer func() {
		image.wait.Wait()
		close(image.fileCaptureEvents)
	}()

	image.outputPath = outputPath

	if err := image.create(".tar"); err != nil {
		return err
	}

	if err := image.pack(); err != nil {
		return err
	}

	if err := image.finalize(); err != nil {
		return err
	}

	return nil
}

func (image *Image) Path() string {
	return image.targetPath
}
