package generator

import (
	"fmt"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/entrypoint"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/gomod"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/nop"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Generator interface {
	Generate() ([]model.File, error)
}

type Func func() ([]model.File, error)

func (f Func) Generate() ([]model.File, error) {
	return f()
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

func Gitignore(module *model.App) Generator {
	return nop.New()
}

func GoMod(app *model.App) Generator {
	return gomod.New(app)
}

func GolangciLint(module *model.App) Generator {
	return nop.New()
}

func Makefile(module *model.App) Generator {
	return nop.New()
}

func Docker(module *model.App) Generator {
	return nop.New()
}

func DockerCompose(module *model.App) Generator {
	return nop.New()
}

func Application(app *model.App) Generator {
	return Func(func() ([]model.File, error) {
		var pipe Pipeline

		for _, ep := range app.Entrypoints() {
			pipe = pipe.With(entrypoint.New(app, ep))
		}

		for _, c := range app.Components() {
			pipe = pipe.With(component.New(app, c))
		}

		return pipe.Generate()
	})
}
