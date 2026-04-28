package model

import (
	"errors"
	"fmt"

	"github.com/iancoleman/strcase"
)

var (
	ErrUnknownDriverKind   = errors.New("unknown driver kind")
	ErrUndefinedDriverName = errors.New("undefined driver name")
)

type Driver struct {
	name   string
	kind   DriverKind
	fqdn   string
	config Config
	path   string
}

func NewDriver(name string, kind DriverKind) (*Driver, error) {
	switch kind {
	case DriverKindPostgres:
		return NewPostgres(name)
	case DriverKindRedis:
		return NewRedis(name)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownDriverKind, kind)
	}
}

func NewPostgres(name string) (*Driver, error) {
	if name == "" {
		return nil, ErrUndefinedDriverName
	}

	driver := &Driver{
		fqdn: FQDN(DriverLayer.Group(), name),
		name: name,
		kind: DriverKindPostgres,
		path: "github.com/vladimirkorzhenevskiy/gogen/pkg/postgres",
		config: Config{
			Must(NewSetting("HOST", SettingTypeString, false, "localhost")),
			Must(NewSetting("PORT", SettingTypeUint16, false, "5432")),
			Must(NewSetting("USERNAME", SettingTypeString, true, "")),
			Must(NewSetting("PASSWORD", SettingTypeString, true, "")),
			Must(NewSetting("DB", SettingTypeString, true, "")),
			Must(NewSetting("SSL_MODE", SettingTypeString, false, "disable")),
			Must(NewSetting("MAX_CONNS", SettingTypeInt32, false, "10")),
		},
	}

	return driver, nil
}

func NewRedis(name string) (*Driver, error) {
	if name == "" {
		return nil, ErrUndefinedDriverName
	}

	name = strcase.ToLowerCamel(name)

	driver := &Driver{
		fqdn: FQDN(DriverLayer.Group(), name),
		name: name,
		kind: DriverKindRedis,
		path: "github.com/vladimirkorzhenevskiy/gogen/pkg/redis",
		config: Config{
			Must(NewSetting("HOST", SettingTypeString, false, "localhost")),
			Must(NewSetting("PORT", SettingTypeUint16, false, "6379")),
			Must(NewSetting("USERNAME", SettingTypeString, false, "")),
			Must(NewSetting("PASSWORD", SettingTypeString, false, "")),
			Must(NewSetting("DB", SettingTypeInt8, false, "0")),
		},
	}

	return driver, nil
}

func (d Driver) FQDN() string {
	return d.fqdn
}

func (d Driver) Name() string {
	return strcase.ToCamel(d.name)
}

func (d Driver) FullName() string {
	return strcase.ToCamel(d.name)
}

func (d Driver) Namespace() string {
	return ""
}

func (d Driver) Layer() Layer {
	return DriverLayer
}

func (d Driver) Alias() string {
	return Concat(strcase.ToCamel(d.name), strcase.ToCamel(DriverLayer.String()))
}

func (d Driver) Path() string {
	return d.path
}

func (d Driver) Config() Config {
	return d.config
}

func (d Driver) Dependencies() Dependencies { return nil }

func (d Driver) Hooks() []Hook { return nil }

type DriverKind string

const (
	DriverKindPostgres DriverKind = "postgres"
	DriverKindRedis    DriverKind = "redis"
)
