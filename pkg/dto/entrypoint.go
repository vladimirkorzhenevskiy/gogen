package dto

import (
	"maps"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Entrypoint struct {
	Name        string   `yaml:"name"        validate:"required"`
	Controllers []string `yaml:"controllers" validate:"required,min=1"`
}

func (e *Entrypoint) ToModel() (*model.Entrypoint, error) {
	return model.NewEntrypoint(
		e.Name,
		e.Controllers,
	)
}

type Entrypoints []*Entrypoint

func (e *Entrypoints) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]*Entrypoint

	if err := unmarshal(&raw); err != nil {
		return err
	}

	res := make(Entrypoints, 0, len(raw))

	for _, name := range slices.Sorted(maps.Keys(raw)) {
		entrypoint := raw[name]
		if entrypoint == nil {
			entrypoint = &Entrypoint{}
		}

		entrypoint.Name = name

		res = append(res, entrypoint)
	}

	*e = res

	return nil
}

func (e Entrypoints) ToModel() (model.Entrypoints, error) {
	return slicesutil.TryMap(e, (*Entrypoint).ToModel)
}
