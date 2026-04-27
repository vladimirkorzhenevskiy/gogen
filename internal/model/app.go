package model

import (
	"fmt"
	"maps"
	"path"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

var (
	ErrUndefinedAppName       = fmt.Errorf("undefined app name")
	ErrUndefinedAppModuleName = fmt.Errorf("undefined app module name")
	ErrCycleDetected          = fmt.Errorf("cycle detected")
	ErrDuplicatesDetected     = fmt.Errorf("duplicates detected")
)

type App struct {
	name        string
	module      string
	entrypoints Entrypoints
	components  map[string]Component
}

func NewApp(
	name string,
	module string,
	entrypoints Entrypoints,
	components Components,
) (*App, error) {
	if name == "" {
		return nil, ErrUndefinedAppName
	}

	if module == "" {
		return nil, ErrUndefinedAppModuleName
	}

	if len(entrypoints) == 0 {
		entrypoints = Entrypoints{
			{
				name:        name,
				controllers: components.Controllers().FQDN(),
			},
		}
	}

	app := &App{
		name:        name,
		module:      module,
		entrypoints: entrypoints,
		components:  make(map[string]Component),
	}

	if err := app.resolveAll(components); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Module() string {
	return a.module
}

func (a *App) GoVersion() string {
	return "1.25"
}

func (a *App) Layers() []Layer {
	return []Layer{
		DriverLayer,
		ControllerLayer,
		MiddlewareLayer,
		UseCaseLayer,
		AdapterLayer,
		RepositoryLayer,
	}
}

func (a *App) Entrypoints() Entrypoints {
	return a.entrypoints
}

func (a *App) Components() Components {
	return slices.Collect(maps.Values(a.components))
}

func (a *App) Import(fqdn string) string {
	component, ok := a.components[fqdn]
	if !ok {
		return path.Join(a.module, fqdn)
	}

	if Is[*Driver](component) {
		return component.Path()
	}

	return path.Join(a.module, component.Path())
}

func (a *App) Component(fqdn string) Component {
	return a.components[fqdn]
}

func (a *App) resolveAll(components Components) error {
	registry := make(map[string]Component, len(components))

	for _, item := range components {
		if _, ok := registry[item.FQDN()]; ok {
			return fmt.Errorf("%w: %v", ErrDuplicatesDetected, item.FQDN())
		}

		registry[item.FQDN()] = item
	}

	return slicesutil.TryEach(a.entrypoints, func(entrypoint *Entrypoint) error {
		return a.resolve(registry, entrypoint)
	})
}

func (a *App) resolve(registry map[string]Component, entrypoint *Entrypoint) error {
	components, err := resolve(registry, entrypoint.controllers...)
	if err != nil {
		return err
	}

	entrypoint.components = components

	for _, component := range components {
		a.components[component.FQDN()] = component
	}

	return nil
}

func resolve(registry map[string]Component, controllers ...string) ([]Component, error) {
	visited := make(map[string]bool)
	temporary := make(map[string]bool)
	sorted := make([]Component, 0)

	var visit func(fqdn string) error

	visit = func(fqdn string) error {
		comp, exists := registry[fqdn]
		if !exists {
			return fmt.Errorf("link error: component %s not found", fqdn)
		}

		if temporary[fqdn] {
			return fmt.Errorf("%w: %s", ErrCycleDetected, fqdn)
		}

		if !visited[fqdn] {
			temporary[fqdn] = true

			// Рекурсивно обходим все зависимости компонента
			for _, dep := range comp.Dependencies() {
				if err := visit(dep.fqdn); err != nil {
					return err
				}
			}

			visited[fqdn] = true
			temporary[fqdn] = false
			sorted = append(sorted, comp)
		}

		return nil
	}

	for _, controller := range controllers {
		if err := visit(controller); err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
