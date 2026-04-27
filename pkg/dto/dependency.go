package dto

import (
	"maps"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Dependency struct {
	Name string `validate:"required"`
	FQDN string `validate:"required"`
}

func (d *Dependency) ToModel() (*model.Dependency, error) {
	return model.NewDependency(d.Name, d.FQDN)
}

type Dependencies []*Dependency

func (d *Dependencies) UnmarshalYAML(unmarshal func(any) error) error {
	var items map[string]string

	if err := unmarshal(&items); err != nil {
		return err
	}

	res := make(Dependencies, 0, len(items))

	for _, name := range slices.Sorted(maps.Keys(items)) {
		res = append(res, &Dependency{
			Name: name,
			FQDN: items[name],
		})
	}

	*d = res

	return nil
}

func (d Dependencies) ToModel() ([]*model.Dependency, error) {
	return slicesutil.TryMap(d, (*Dependency).ToModel)
}
