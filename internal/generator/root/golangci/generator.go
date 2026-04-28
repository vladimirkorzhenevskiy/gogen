package golangci

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	app *model.App
}

func New(app *model.App) Generator {
	return Generator{app: app}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := template.NewPipeline(template.MustSub(templates, "templates"))

	pipe.With(template.Text{
		Template:  ".golangci.yaml.tmpl",
		Filename:  ".golangci.yaml",
		Overwrite: false,
		Data:      g.app,
	})

	return pipe.Execute()
}
