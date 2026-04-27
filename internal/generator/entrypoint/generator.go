package entrypoint

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/core"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type TemplateContext struct {
	App        *model.App
	Entrypoint *model.Entrypoint
	Components model.Components
}

type Generator struct {
	app        *model.App
	entrypoint *model.Entrypoint
}

func New(
	app *model.App,
	entrypoint *model.Entrypoint,
) Generator {
	return Generator{
		app:        app,
		entrypoint: entrypoint,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := template.NewPipeline(core.RootFS(), template.MustSub(templates, "templates"))

	pipe = pipe.With(
		template.Go{
			Template:  "main_gen.go.tmpl",
			Filepath:  g.entrypoint.Path(),
			Filename:  "main_gen.go",
			Overwrite: true,
			Data: TemplateContext{
				App:        g.app,
				Entrypoint: g.entrypoint,
				Components: g.entrypoint.Components(),
			},
		},
	)

	return pipe.Execute()
}
