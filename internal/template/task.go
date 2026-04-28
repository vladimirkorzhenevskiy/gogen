package template

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"text/template"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Task interface {
	Execute(tpl *template.Template) (model.File, error)
}

type Go struct {
	Template  string
	Filepath  string
	Filename  string
	Overwrite bool
	Data      any
}

func (t Go) Execute(tpl *template.Template) (model.File, error) {
	var buf bytes.Buffer

	if err := tpl.ExecuteTemplate(&buf, t.Template, t.Data); err != nil {
		return model.File{}, err
	}

	content, err := format.Source(buf.Bytes())
	if err != nil {
		return model.File{}, fmt.Errorf("format error: %w", err)
	}

	return model.NewFile(path.Join(t.Filepath, t.Filename), content, t.Overwrite), nil
}

type Text struct {
	Template  string
	Filepath  string
	Filename  string
	Overwrite bool
	Data      any
}

func (t Text) Execute(tpl *template.Template) (model.File, error) {
	var buf bytes.Buffer

	if err := tpl.ExecuteTemplate(&buf, t.Template, t.Data); err != nil {
		return model.File{}, err
	}

	return model.NewFile(path.Join(t.Filepath, t.Filename), buf.Bytes(), t.Overwrite), nil
}
