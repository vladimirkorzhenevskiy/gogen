package dto

import (
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type App struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec"`
}

type Metadata struct {
	Name   string `yaml:"name"   validate:"required"`
	Module string `yaml:"module" validate:"required"`
}

type Spec struct {
	Drivers      Drivers      `yaml:"drivers"`
	Entrypoints  Entrypoints  `yaml:"entrypoints"  validate:"dive"`
	Controllers  Controllers  `yaml:"controllers"  validate:"dive"`
	Middlewares  Middlewares  `yaml:"middlewares"  validate:"dive"`
	UseCases     UseCases     `yaml:"usecases"     validate:"dive"`
	Adapters     Adapters     `yaml:"adapters"     validate:"dive"`
	Repositories Repositories `yaml:"repositories" validate:"dive"`
}

func (a App) ToModel() (*model.App, error) {
	builder := model.NewBuilder(a.Metadata.Name, a.Metadata.Module)

	entrypoints, err := a.Spec.Entrypoints.ToModel()
	if err != nil {
		return nil, err
	}
	builder.WithEntrypoints(entrypoints)

	providers := []func() (model.Components, error){
		a.Spec.Drivers.ToModel,
		a.Spec.Controllers.ToModel,
		a.Spec.Middlewares.ToModel,
		a.Spec.UseCases.ToModel,
		a.Spec.Adapters.ToModel,
		a.Spec.Repositories.ToModel,
	}

	for _, provider := range providers {
		components, err := provider()
		if err != nil {
			return nil, err
		}
		builder.WithComponents(components)
	}

	return builder.Build()
}
