package common

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/base"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	controller *model.Controller
}

func New(controller *model.Controller) Generator {
	return Generator{
		controller: controller,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := template.NewPipeline(base.RootFS(), template.MustSub(templates, "templates"))

	pipe = pipe.With(
		template.Go{
			Template:  "controller.go.tmpl",
			Filepath:  g.controller.Path(),
			Filename:  "controller.go",
			Overwrite: false,
			Data:      g.controller,
		},
		template.Go{
			Template:  "controller_gen.go.tmpl",
			Filepath:  g.controller.Path(),
			Filename:  "controller_gen.go",
			Overwrite: true,
			Data:      g.controller,
		},
	)

	return pipe.Execute()
}
