package fsmatch

import (
	"os"
	"path"
	"path/filepath"
)

func Match(pattern, directory string) ([]string, error) {
	matches := make([]string, 0)

	matchWalker := func(walkedPath string, walkedInfo os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if walkedInfo.IsDir() {
			return nil
		}

		ok, err := path.Match(pattern, path.Base(walkedPath))
		if err != nil {
			return err
		}

		if ok {
			matches = append(matches, walkedPath)
		}

		return nil
	}

	if err := filepath.Walk(directory, matchWalker); err != nil {
		return make([]string, 0), err
	}

	return matches, nil
}
