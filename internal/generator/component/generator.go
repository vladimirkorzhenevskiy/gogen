package component

import (
	"errors"
	"fmt"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/adapter"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/controller"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/middleware"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/nop"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/repository"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/usecase"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

var errUnknownComponentType = errors.New("unknown component type")

type Generator struct {
	app       *model.App
	component model.Component
}

func New(app *model.App, component model.Component) Generator {
	return Generator{
		app:       app,
		component: component,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	switch component := g.component.(type) {
	case *model.Driver:
		return nop.New().Generate()
	case *model.Controller:
		return controller.New(g.app, component).Generate()
	case *model.Middleware:
		return middleware.New(component).Generate()
	case *model.UseCase:
		return usecase.New(component).Generate()
	case *model.Adapter:
		return adapter.New(component).Generate()
	case *model.Repository:
		return repository.New(component).Generate()
	default:
		return nil, fmt.Errorf("%w: %T", errUnknownComponentType, component)
	}
}
