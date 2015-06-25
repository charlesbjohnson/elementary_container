package compress

import (
	"os"

	"github.com/charlesbjohnson/elementary_container/fscapture/image"
)

const extension = ".gz"

type Image struct {
	*image.Image
	targetPath string
	targetFile *os.File
}

func New(inputPath, outputPath string) *Image {
	image := image.New(inputPath, outputPath)

	return &Image{
		Image:      image,
		targetPath: image.Path() + extension,
	}
}

func (image *Image) Exists() bool {
	result := true

	if _, err := os.Stat(image.targetPath); err != nil {
		result = false
	}

	return result
}

func (image *Image) Capture() error {
	if image.Exists() {
		image.Image.Close()
		return os.ErrExist
	}

	if err := image.Image.Capture(); err != nil {
		return err
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
	if err := image.Image.Close(); err != nil {
		return err
	}

	if image.targetFile != nil {
		if err := image.targetFile.Close(); err != nil {
			return err
		}
	}

	return nil
}
