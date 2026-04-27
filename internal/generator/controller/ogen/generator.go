package ogen

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/core"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/ogen/server"
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
	res, err := server.Generate(g.controller.Spec())
	if err != nil {
		return nil, err
	}

	pipe := template.NewPipeline(core.RootFS(), template.MustSub(templates, "templates"))

	pipe = pipe.With(
		template.Go{
			Template:  "controller.go.tmpl",
			Filepath:  g.controller.Path(),
			Filename:  "controller.go",
			Overwrite: true,
			Data:      g.controller,
		},
		template.Go{
			Template:  "controller_gen.go.tmpl",
			Filepath:  g.controller.Path(),
			Filename:  "controller_gen.go",
			Overwrite: false,
			Data:      g.controller,
		},
	)

	files, err := pipe.Execute()
	if err != nil {
		return nil, err
	}

	res = append(res, files...)

	return res, nil
}
