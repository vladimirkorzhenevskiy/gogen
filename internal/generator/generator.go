package generator

import (
	"fmt"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/nop"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/entrypoint"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/root"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Generator interface {
	Generate() ([]model.File, error)
}

type Pipeline []Generator

func (p Pipeline) With(generators ...Generator) Pipeline {
	return append(p, generators...)
}

func (p Pipeline) Generate() ([]model.File, error) {
	var res []model.File

	for _, generator := range p {
		files, err := generator.Generate()
		if err != nil {
			return nil, fmt.Errorf("%T: generate: %w", generator, err)
		}

		res = append(res, files...)
	}

	return res, nil
}

func App(app *model.App) Generator {
	return Pipeline{
		Root(app),
		Code(app),
		Build(app),
	}
}

func Root(app *model.App) Generator {
	return Pipeline{
		root.Gitignore(app),
		root.GolangCI(app),
		root.GoMod(app),
	}
}

func Code(app *model.App) Generator {
	return Pipeline{
		Pipeline(slicesutil.Map(app.Entrypoints(), func(item *model.Entrypoint) Generator {
			return entrypoint.New(app, item)
		})),
		Pipeline(slicesutil.Map(app.Components(), func(item model.Component) Generator {
			return component.New(app, item)
		})),
	}
}

func Build(app *model.App) Generator {
	return nop.New()
}
