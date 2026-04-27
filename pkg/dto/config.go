package dto

import (
	"maps"
	"slices"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Setting struct {
	Name     string `yaml:"name"      validate:"required"`
	Type     string `yaml:"type"      validate:"required"`
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
}

func (s *Setting) ToModel() (*model.Setting, error) {
	return model.NewSetting(
		s.Name,
		model.SettingType(s.Type),
		s.Required,
		s.Default,
	)
}

type Config []*Setting

func (c *Config) UnmarshalYAML(unmarshal func(any) error) error {
	var settings map[string]*Setting

	if err := unmarshal(&settings); err != nil {
		return err
	}

	res := make(Config, 0, len(settings))

	for _, name := range slices.Sorted(maps.Keys(settings)) {
		setting := settings[name]
		if setting == nil {
			setting = &Setting{}
		}

		setting.Name = name

		res = append(res, setting)
	}

	*c = res

	return nil
}

func (c Config) ToModel() (model.Config, error) {
	return slicesutil.TryMap(c, (*Setting).ToModel)
}
