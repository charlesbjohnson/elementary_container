package compress

import (
	"compress/gzip"
	"io"
	"os"
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
	writer := gzip.NewWriter(image.targetFile)

	if _, err := io.Copy(writer, image.Image.File()); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

func (image *Image) finalize() error {
	if err := image.Image.File().Close(); err != nil {
		return err
	}

	if err := os.Remove(image.Image.Path()); err != nil {
		return err
	}

	return nil
}
