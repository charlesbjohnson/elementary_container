package fstemplate

// TODO how can this be made generic to text or html templates?

import (
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/charlesbjohnson/elementary_container/fsmatch"
)

func New(pattern, directory string) (*template.Template, error) {
	rootTemplate := template.New(directory)

	matches, err := fsmatch.Match(pattern, directory)
	if err != nil {
		return rootTemplate, err
	}

	for _, path := range matches {
		relative, err := filepath.Rel(directory, path)
		if err != nil {
			continue
		}

		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		template := rootTemplate.New(relative)
		if _, err := template.Parse(string(buffer)); err != nil {
			continue
		}
	}

	return rootTemplate, nil
}
