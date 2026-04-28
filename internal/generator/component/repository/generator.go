package repository

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/base"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	repository *model.Repository
}

func New(repository *model.Repository) Generator {
	return Generator{
		repository: repository,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := base.Pipeline(templates).With(
		template.Go{
			Template:  "repository.go.tmpl",
			Filepath:  g.repository.Path(),
			Filename:  "repository.go",
			Overwrite: false,
			Data:      g.repository,
		},
		template.Go{
			Template:  "repository_gen.go.tmpl",
			Filepath:  g.repository.Path(),
			Filename:  "repository_gen.go",
			Overwrite: true,
			Data:      g.repository,
		},
	)

	return pipe.Execute()
}
