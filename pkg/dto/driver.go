package dto

import (
	"maps"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Driver struct {
	Name string `yaml:"name" validate:"required"`
	Kind string `yaml:"kind" validate:"required"`
}

func (d *Driver) ToModel() (model.Component, error) {
	return model.NewDriver(d.Name, model.DriverKind(d.Kind))
}

type Drivers []*Driver

func (d *Drivers) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]*Driver

	if err := unmarshal(&raw); err != nil {
		return err
	}

	res := make(Drivers, 0, len(raw))

	for _, name := range slices.Sorted(maps.Keys(raw)) {
		driver := raw[name]
		if driver == nil {
			driver = &Driver{}
		}

		driver.Name = name

		res = append(res, driver)
	}

	*d = res

	return nil
}

func (d Drivers) ToModel() (model.Components, error) {
	return slicesutil.TryMap(d, (*Driver).ToModel)
}
