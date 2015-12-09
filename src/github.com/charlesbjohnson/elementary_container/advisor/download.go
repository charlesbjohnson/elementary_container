package advisor

import (
	"io"
	"os"
)

type downloader struct {
	strategy func(io.WriterAt) error
}

func (application *Application) download(downloader *downloader, downloadPath string) error {
	file, err := os.Create(downloadPath)
	if err != nil {
		return err
	}

	if err := downloader.strategy(file); err != nil {
		return err
	}

	application.Log.WithField("file", file.Name()).Info("image downloaded")

	return nil
}
