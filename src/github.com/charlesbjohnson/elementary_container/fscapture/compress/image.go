package compress

import (
	"github.com/charlesbjohnson/elementary_container/fscapture/image"
	"os"
)

type Image struct {
	*image.Image
	targetPath string
	targetFile *os.File
}

func New(inputPath string) *Image {
	return &Image{Image: image.New(inputPath)}
}

func (image *Image) Capture(outputPath string) error {
	if err := image.Image.Capture(outputPath); err != nil {
		return err
	}

	if err := image.create(".gz"); err != nil {
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
