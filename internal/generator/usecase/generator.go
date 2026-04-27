package usecase

import (
	"embed"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/core"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

type Generator struct {
	usecase *model.UseCase
}

func New(usecase *model.UseCase) Generator {
	return Generator{
		usecase: usecase,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	pipe := template.NewPipeline(core.RootFS(), template.MustSub(templates, "templates"))

	pipe = pipe.With(
		template.Go{
			Template:  "use_case.go.tmpl",
			Filepath:  g.usecase.Path(),
			Filename:  "use_case.go",
			Overwrite: false,
			Data:      g.usecase,
		},
		template.Go{
			Template:  "use_case_gen.go.tmpl",
			Filepath:  g.usecase.Path(),
			Filename:  "use_case_gen.go",
			Overwrite: true,
			Data:      g.usecase,
		},
	)

	return pipe.Execute()
}
