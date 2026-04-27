package model

import (
	"errors"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

var (
	ErrUnknownComponentLayer       = errors.New("unknown component layer")
	ErrUndefinedComponentName      = errors.New("undefined component name")
	ErrUndefinedComponentNamespace = errors.New("undefined component namespace")
	ErrUndefinedComponentLayer     = errors.New("undefined component layer")
)

type Component interface {
	FQDN() string
	Name() string
	FullName() string
	Namespace() string
	Layer() Layer
	Path() string
	Alias() string
	Config() Config
	Dependencies() Dependencies
}

func Is[T any](value Component) bool {
	if _, ok := value.(T); ok {
		return true
	}

	return false
}

type Components []Component

func (c Components) Filter(predicate func(Component) bool) Components {
	return slicesutil.Filter(c, predicate)
}

func (c Components) WithLayer(layer Layer) Components {
	return slicesutil.Filter(c, func(component Component) bool {
		return component.Layer() == layer
	})
}

func (c Components) WithConfig() Components {
	return slicesutil.Filter(c, func(component Component) bool {
		return len(component.Config()) != 0
	})
}

func (c Components) Drivers() Components {
	return c.Filter(Is[*Driver])
}

func (c Components) Controllers() Components {
	return c.Filter(Is[*Controller])
}

func (c Components) Middlewares() Components {
	return c.Filter(Is[*Middleware])
}

func (c Components) UseCases() Components {
	return c.Filter(Is[*UseCase])
}

func (c Components) Adapters() Components {
	return c.Filter(Is[*Adapter])
}

func (c Components) Repositories() Components {
	return c.Filter(Is[*Repository])
}

func (c Components) ByNamespace() map[string]Components {
	res := make(map[string]Components, len(c))

	for _, component := range c {
		components := res[component.Namespace()]
		components = append(components, component)

		res[component.Namespace()] = components
	}

	return res
}

func (c Components) FQDN() []string {
	return slicesutil.Map(c, Component.FQDN)
}

func FQDN(parts ...string) string {
	var i int

	for j := range parts {
		if parts[j] == "" {
			continue
		}

		parts[i] = strings.ToLower(parts[j])

		i++
	}

	return strings.Join(parts[:i], ".")
}

func Concat(parts ...string) string {
	return strings.Join(parts, "")
}

type BaseComponent struct {
	name         string
	namespace    string
	layer        Layer
	config       Config
	dependencies Dependencies
	path         string
	alias        string
}

func baseComponent(name, namespace string, layer Layer, config Config, dependencies Dependencies) (*BaseComponent, error) {
	if name == "" {
		return nil, ErrUndefinedComponentName
	}

	if namespace == "" {
		return nil, ErrUndefinedComponentNamespace
	}

	if layer == "" {
		return nil, ErrUndefinedComponentLayer
	}

	if !layer.IsValid() {
		return nil, ErrUnknownComponentLayer
	}

	name = strcase.ToLowerCamel(name)
	namespace = strcase.ToLowerCamel(namespace)

	component := &BaseComponent{
		name:         name,
		namespace:    namespace,
		layer:        layer,
		config:       config,
		dependencies: dependencies,
	}

	return component, nil
}

func (b *BaseComponent) Name() string {
	return strcase.ToCamel(b.name)
}

func (b *BaseComponent) FullName() string {
	return strcase.ToCamel(b.namespace) + strcase.ToCamel(b.name)
}

func (b *BaseComponent) Namespace() string { return b.namespace }

func (b *BaseComponent) Layer() Layer { return b.layer }

func (b *BaseComponent) Config() Config { return b.config }

func (b *BaseComponent) FQDN() string {
	return FQDN(b.layer.Group(), b.namespace, b.name)
}

func (b *BaseComponent) Path() string {
	return path.Join(
		"internal",
		strings.ToLower(b.layer.String()),
		strings.ToLower(b.namespace),
		strings.ToLower(b.name),
	)
}

func (b *BaseComponent) Alias() string {
	return Concat(
		strcase.ToCamel(b.namespace),
		strcase.ToCamel(b.name),
		strcase.ToCamel(b.layer.String()),
	)
}

func (b *BaseComponent) Dependencies() Dependencies {
	return b.dependencies
}

type Adapter struct {
	*BaseComponent
}

func NewAdapter(name, namespace string, config Config, dependencies Dependencies) (*Adapter, error) {
	base, err := baseComponent(name, namespace, AdapterLayer, config, dependencies)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		BaseComponent: base,
	}, nil
}

type Middleware struct {
	*BaseComponent
}

func NewMiddleware(name, namespace string, config Config, dependencies Dependencies) (*Middleware, error) {
	base, err := baseComponent(name, namespace, MiddlewareLayer, config, dependencies)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		BaseComponent: base,
	}, nil
}

type Repository struct {
	*BaseComponent
}

func NewRepository(name, namespace string, config Config, dependencies Dependencies) (*Repository, error) {
	base, err := baseComponent(name, namespace, RepositoryLayer, config, dependencies)
	if err != nil {
		return nil, err
	}

	return &Repository{
		BaseComponent: base,
	}, nil
}

type UseCase struct {
	*BaseComponent
}

func NewUseCase(name, namespace string, config Config, dependencies Dependencies) (*UseCase, error) {
	base, err := baseComponent(name, namespace, UseCaseLayer, config, dependencies)
	if err != nil {
		return nil, err
	}

	return &UseCase{
		BaseComponent: base,
	}, nil
}
