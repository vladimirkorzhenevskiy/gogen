package adapter

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/base"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	adapter *model.Adapter
}

func New(adapter *model.Adapter) Generator {
	return Generator{
		adapter: adapter,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := base.Pipeline(templates).With(
		template.Go{
			Template:  "adapter.go.tmpl",
			Filepath:  g.adapter.Path(),
			Filename:  "adapter.go",
			Overwrite: false,
			Data:      g.adapter,
		},
		template.Go{
			Template:  "adapter_gen.go.tmpl",
			Filepath:  g.adapter.Path(),
			Filename:  "adapter_gen.go",
			Overwrite: true,
			Data:      g.adapter,
		},
	)

	return pipe.Execute()
}
