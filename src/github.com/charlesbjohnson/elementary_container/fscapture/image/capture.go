package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charlesbjohnson/elementary_container/fscapture"
)

func (image *Image) create() error {
	file, err := os.Create(image.targetPath)
	if err != nil {
		return err
	}

	image.targetFile = file

	return nil
}

func (image *Image) pack() error {
	image.writer = tar.NewWriter(image.targetFile)

	walk := func(path string, info os.FileInfo, err error) error {
		return image.walkPath(path, info, err)
	}

	if err := filepath.Walk(image.inputPath, walk); err != nil {
		return err
	}

	return nil
}

func (image *Image) walkPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		image.emit(path, err.Error(), info, false)
		return filepath.SkipDir
	}

	if info.IsDir() {
		return nil
	}

	for _, handler := range image.handlers {
		if include := handler.Hook(path, info); !include {
			image.emit(path, fmt.Sprintf("file ignored by handler %s", handler.Id), info, false)
			return nil
		}
	}

	if err := image.packPath(path, info); err != nil {
		image.emit(path, "file pack failed", info, false)
	}

	image.emit(path, "file pack succeeded", info, true)

	return nil
}

func (image *Image) packPath(path string, info os.FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    path,
		Mode:    int64(info.Mode()),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if err := image.writer.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(image.writer, file); err != nil {
		return err
	}

	return nil
}

func (image *Image) emit(path, message string, info os.FileInfo, captured bool) {
	image.wait.Add(1)

	go func() {
		defer image.wait.Done()
		image.fileCaptureEvents <- fscapture.FileCaptureEvent{path, message, info, captured}
	}()
}

func (image *Image) finalize() error {
	if err := image.writer.Close(); err != nil {
		return err
	}

	return nil
}

func (image *Image) endEmit() {
	if image.fileCaptureEventsClosed {
		return
	}

	image.wait.Wait()
	close(image.fileCaptureEvents)

	image.fileCaptureEventsClosed = true
}
