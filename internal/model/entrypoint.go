package model

import "path"

type Entrypoint struct {
	name        string
	controllers []string
	components  []Component
}

func NewEntrypoint(name string, controllers []string) (*Entrypoint, error) {
	return &Entrypoint{
		name:        name,
		controllers: controllers,
	}, nil
}

func (e Entrypoint) Name() string {
	return e.name
}

func (e Entrypoint) Path() string {
	return path.Join("cmd", e.Name())
}

func (e Entrypoint) Components() []Component {
	return e.components
}

type Entrypoints []*Entrypoint
