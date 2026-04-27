package model

import "errors"

var (
	ErrUndefinedDependencyName = errors.New("undefined dependency name")
	ErrUndefinedDependencyFQDN = errors.New("undefined dependency fqdn")
)

type Dependency struct {
	name string
	fqdn string
}

func NewDependency(name, fqdn string) (*Dependency, error) {
	if name == "" {
		return nil, ErrUndefinedDependencyName
	}

	if fqdn == "" {
		return nil, ErrUndefinedDependencyFQDN
	}

	return &Dependency{
		name: name,
		fqdn: fqdn,
	}, nil
}

func (d *Dependency) Name() string {
	return d.name
}

func (d *Dependency) FQDN() string {
	return d.fqdn
}

type Dependencies []*Dependency
