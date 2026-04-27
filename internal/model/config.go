package model

import (
	"errors"
)

var (
	ErrUndefinedSettingName = errors.New("undefined setting name")
	ErrInvalidSettingType   = errors.New("invalid setting type")
)

type Config []*Setting

func (c Config) Merge(config Config) Config {
	return append(c, config...)
}

type Setting struct {
	name         string
	settingType  SettingType
	required     bool
	defaultValue string
}

func NewSetting(name string, settingType SettingType, required bool, defaultValue string) (*Setting, error) {
	if name == "" {
		return nil, ErrUndefinedSettingName
	}

	if !settingType.IsValid() {
		return nil, ErrInvalidSettingType
	}

	return &Setting{
		name:         name,
		settingType:  settingType,
		required:     required,
		defaultValue: defaultValue,
	}, nil
}

func (s *Setting) Name() string {
	return s.name
}

func (s *Setting) Type() SettingType {
	return s.settingType
}

func (s *Setting) Required() bool {
	return s.required
}

func (s *Setting) Default() string {
	return s.defaultValue
}

const (
	SettingTypeString   SettingType = "string"
	SettingTypeInt      SettingType = "int"
	SettingTypeInt8     SettingType = "int8"
	SettingTypeInt16    SettingType = "int16"
	SettingTypeInt32    SettingType = "int32"
	SettingTypeInt64    SettingType = "int64"
	SettingTypeUint     SettingType = "uint"
	SettingTypeUint8    SettingType = "uint8"
	SettingTypeUint16   SettingType = "uint16"
	SettingTypeUint32   SettingType = "uint32"
	SettingTypeUint64   SettingType = "uint64"
	SettingTypeFloat32  SettingType = "float32"
	SettingTypeFloat64  SettingType = "float64"
	SettingTypeBool     SettingType = "bool"
	SettingTypeDuration SettingType = "duration"
)

type SettingType string

func (st SettingType) String() string {
	return string(st)
}

func (st SettingType) IsValid() bool {
	switch st {
	case SettingTypeString,
		SettingTypeInt,
		SettingTypeInt8,
		SettingTypeInt16,
		SettingTypeInt32,
		SettingTypeInt64,
		SettingTypeUint,
		SettingTypeUint8,
		SettingTypeUint16,
		SettingTypeUint32,
		SettingTypeUint64,
		SettingTypeFloat32,
		SettingTypeFloat64,
		SettingTypeBool,
		SettingTypeDuration:
		return true
	}

	return false
}
