package advisor

import (
	"sync"

	"github.com/charlesbjohnson/elementary_container/fscapture"
)

func (application *Application) capture(capturable fscapture.Capturable) error {
	var wait sync.WaitGroup
	defer wait.Wait()

	wait.Add(1)
	go func() {
		defer wait.Done()

		for event := range capturable.FileCaptureEvents() {
			logger := application.Log.WithField("file", event.Path)

			if event.Captured {
				logger.Info(event.Message)
			} else {
				logger.Warn(event.Message)
			}
		}
	}()

	if err := capturable.Capture(); err != nil {
		return err
	}

	application.Log.WithField("file", capturable.Path()).Info("image capture succeeded")
	return nil
}
