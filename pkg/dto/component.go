package dto

import (
	"errors"
	"maps"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Component struct {
	Name         string       `yaml:"name"         validate:"required"`
	Namespace    string       `yaml:"namespace"`
	Config       Config       `yaml:"config"       validate:"unique,dive"`
	Dependencies Dependencies `yaml:"dependencies" validate:"unique,dive"`
	Hooks        Hooks        `yaml:"hooks"        validate:"unique"`
}

func (c *Component) SetName(name string) {
	c.Name = name
}

func (c *Component) SetNamespace(namespace string) {
	c.Namespace = namespace
}

func (c *Component) ToModel() (model.Component, error) {
	return nil, errors.New("unsupported component type")
}

type component interface {
	SetName(name string)
	SetNamespace(namespace string)
	ToModel() (model.Component, error)
}

type Components[T component] []T

func (c *Components[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]map[string]T

	if err := unmarshal(&raw); err != nil {
		return err
	}

	res := make(Components[T], 0, len(raw))

	for _, namespace := range slices.Sorted(maps.Keys(raw)) {
		items := raw[namespace]

		for _, name := range slices.Sorted(maps.Keys(items)) {
			item := items[name]
			item.SetName(name)
			item.SetNamespace(namespace)

			res = append(res, item)
		}
	}

	*c = res

	return nil
}

func (c Components[T]) ToModel() (model.Components, error) {
	return slicesutil.TryMap(c, T.ToModel)
}

type Controller struct {
	Component `yaml:",inline"`

	Protocol string `yaml:"protocol"`
	Spec     string `yaml:"spec"`
}

func (c *Controller) ToModel() (model.Component, error) {
	config, err := c.Config.ToModel()
	if err != nil {
		return nil, err
	}

	dependencies, err := c.Dependencies.ToModel()
	if err != nil {
		return nil, err
	}

	hooks, err := c.Hooks.ToModel()
	if err != nil {
		return nil, err
	}

	switch model.Protocol(c.Protocol) {
	case model.ProtocolOgen:
		return model.NewOgenController(
			c.Name,
			c.Namespace,
			c.Spec,
			config,
			dependencies,
			hooks,
		)
	default:
		return model.NewController(
			c.Name,
			c.Namespace,
			config,
			dependencies,
			hooks,
		)
	}
}

type Controllers = Components[*Controller]

type Middleware struct {
	Component `yaml:",inline"`
}

func (m *Middleware) ToModel() (model.Component, error) {
	config, err := m.Config.ToModel()
	if err != nil {
		return nil, err
	}

	dependencies, err := m.Dependencies.ToModel()
	if err != nil {
		return nil, err
	}

	hooks, err := m.Hooks.ToModel()
	if err != nil {
		return nil, err
	}

	return model.NewMiddleware(
		m.Name,
		m.Namespace,
		config,
		dependencies,
		hooks,
	)
}

type Middlewares = Components[*Middleware]

type UseCase struct {
	Component `yaml:",inline"`
}

func (u *UseCase) ToModel() (model.Component, error) {
	config, err := u.Config.ToModel()
	if err != nil {
		return nil, err
	}

	dependencies, err := u.Dependencies.ToModel()
	if err != nil {
		return nil, err
	}

	hooks, err := u.Hooks.ToModel()
	if err != nil {
		return nil, err
	}

	return model.NewUseCase(
		u.Name,
		u.Namespace,
		config,
		dependencies,
		hooks,
	)
}

type UseCases = Components[*UseCase]

type Adapter struct {
	Component `yaml:",inline"`
}

func (a *Adapter) ToModel() (model.Component, error) {
	config, err := a.Config.ToModel()
	if err != nil {
		return nil, err
	}

	dependencies, err := a.Dependencies.ToModel()
	if err != nil {
		return nil, err
	}

	hooks, err := a.Hooks.ToModel()
	if err != nil {
		return nil, err
	}

	return model.NewAdapter(
		a.Name,
		a.Namespace,
		config,
		dependencies,
		hooks,
	)
}

type Adapters = Components[*Adapter]

type Repository struct {
	Component `yaml:",inline"`
}

func (r *Repository) ToModel() (model.Component, error) {
	config, err := r.Config.ToModel()
	if err != nil {
		return nil, err
	}

	dependencies, err := r.Dependencies.ToModel()
	if err != nil {
		return nil, err
	}

	hooks, err := r.Hooks.ToModel()
	if err != nil {
		return nil, err
	}

	return model.NewRepository(
		r.Name,
		r.Namespace,
		config,
		dependencies,
		hooks,
	)
}

type Repositories = Components[*Repository]
