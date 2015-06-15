package compress

import (
	"compress/gzip"
	"io"
	"os"
)

func (image *Image) create(extension string) error {
	finalPath := image.Image.Path() + extension

	if _, err := os.Stat(finalPath); err == nil {
		return os.ErrExist
	}

	file, err := os.Create(finalPath)
	if err != nil {
		return err
	}

	image.targetPath = finalPath
	image.targetFile = file

	return nil
}

func (image *Image) pack() error {
	file, err := os.Open(image.Image.Path())
	if err != nil {
		return err
	}

	writer := gzip.NewWriter(image.targetFile)

	if _, err := io.Copy(writer, file); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

func (image *Image) finalize() error {
	if err := image.targetFile.Close(); err != nil {
		return err
	}

	if err := os.Remove(image.Image.Path()); err != nil {
		return err
	}

	return nil
}
