package image

import (
	"archive/tar"
	"os"
	"sync"

	"github.com/charlesbjohnson/elementary_container/fscapture"
)

const extension = ".tar"

type fileCaptureHookRegistration struct {
	Id   string
	Hook fscapture.FileCaptureHook
}

type Image struct {
	inputPath               string
	targetPath              string
	targetFile              *os.File
	writer                  *tar.Writer
	handlers                []fileCaptureHookRegistration
	fileCaptureEvents       chan fscapture.FileCaptureEvent
	wait                    sync.WaitGroup
	fileCaptureEventsClosed bool
}

func New(inputPath, outputPath string) *Image {
	return &Image{
		inputPath:         inputPath,
		targetPath:        outputPath + extension,
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

func (image *Image) Exists() bool {
	result := true

	if _, err := os.Stat(image.targetPath); err != nil {
		result = false
	}

	return result
}

func (image *Image) Capture() error {
	defer image.endEmit()

	if image.Exists() {
		return os.ErrExist
	}

	if err := image.create(); err != nil {
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

func (image *Image) File() *os.File {
	return image.targetFile
}

func (image *Image) Close() error {
	if image.targetFile != nil {
		if err := image.targetFile.Close(); err != nil {
			return err
		}
	}

	image.endEmit()

	return nil
}
