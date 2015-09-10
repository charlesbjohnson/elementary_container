package headmaster

import (
	"net/http"
	"text/template"

	"github.com/charlesbjohnson/elementary_container/fstemplate"
	"github.com/charlesbjohnson/elementary_container/render/plain"
	"github.com/unrolled/render"
)

type View struct {
	*render.Render
	templates *template.Template
}

func NewView(directory string) (*View, error) {
	templates, err := fstemplate.New("*.tmpl", directory)
	if err != nil {
		return nil, err
	}

	return &View{Render: render.New(), templates: templates}, nil
}

func (view *View) Plain(w http.ResponseWriter, status int, name string, binding interface{}) {
	head := render.Head{
		ContentType: "text/plain; charset=utf-8",
		Status:      status,
	}

	p := plain.Plain{
		Head:      head,
		Name:      name,
		Templates: view.templates,
	}

	view.Render.Render(w, p, binding)
}
