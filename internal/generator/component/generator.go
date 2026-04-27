package component

import (
	"errors"
	"fmt"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/adapter"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/controller/ogen"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/middleware"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/nop"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/repository"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/usecase"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

var errUnknownComponentType = errors.New("unknown component type")

type Generator struct {
	component model.Component
}

func New(component model.Component) Generator {
	return Generator{
		component: component,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	switch component := g.component.(type) {
	case *model.Driver:
		return nop.New().Generate()
	case *model.Controller:
		return ogen.New(component).Generate()
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
