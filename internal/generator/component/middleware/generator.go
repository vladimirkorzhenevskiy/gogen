package middleware

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/base"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	middleware *model.Middleware
}

func New(middleware *model.Middleware) Generator {
	return Generator{
		middleware: middleware,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := base.Pipeline(templates).With(
		template.Go{
			Template:  "middleware.go.tmpl",
			Filepath:  g.middleware.Path(),
			Filename:  "middleware.go",
			Overwrite: false,
			Data:      g.middleware,
		},
		template.Go{
			Template:  "middleware_gen.go.tmpl",
			Filepath:  g.middleware.Path(),
			Filename:  "middleware_gen.go",
			Overwrite: true,
			Data:      g.middleware,
		},
	)

	return pipe.Execute()
}
