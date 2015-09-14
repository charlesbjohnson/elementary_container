package plain

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/unrolled/render"
)

type Plain struct {
	render.Head
	Name      string
	Templates *template.Template
}

func (p Plain) Render(w http.ResponseWriter, binding interface{}) error {
	out := new(bytes.Buffer)

	if err := p.Templates.ExecuteTemplate(out, p.Name, binding); err != nil {
		return err
	}

	p.Head.Write(w)
	out.WriteTo(w)

	return nil
}
