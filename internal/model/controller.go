package model

import "errors"

var (
	ErrUndefinedControllerSpec = errors.New("undefined controller spec")
)

type Controller struct {
	*BaseComponent

	protocol Protocol
	spec     string
}

func NewController(
	name,
	namespace string,
	config Config,
	dependencies Dependencies,
	hooks []Hook,
) (*Controller, error) {
	base, err := baseComponent(name, namespace, ControllerLayer, config, dependencies, hooks)
	if err != nil {
		return nil, err
	}

	return &Controller{
		BaseComponent: base,
	}, nil
}

func NewOgenController(
	name,
	namespace,
	spec string,
	config Config,
	dependencies Dependencies,
	hooks []Hook,
) (*Controller, error) {
	if spec == "" {
		return nil, ErrUndefinedControllerSpec
	}

	ogenConfig := Config{
		Must(NewSetting("HOST", SettingTypeString, false, "localhost")),
		Must(NewSetting("PORT", SettingTypeUint16, false, "8080")),
	}

	base, err := baseComponent(name, namespace, ControllerLayer, ogenConfig.Merge(config), dependencies, hooks)
	if err != nil {
		return nil, err
	}

	return &Controller{
		BaseComponent: base,

		protocol: ProtocolOgen,
		spec:     spec,
	}, nil
}

func (c *Controller) Protocol() Protocol {
	return c.protocol
}

func (c *Controller) Spec() string {
	return c.spec
}

type Protocol string

const (
	ProtocolOgen  Protocol = "ogen"
	ProtocolHTTP  Protocol = "http"
	ProtocolKafka Protocol = "kafka"
	ProtocolGRPC  Protocol = "grpc"
	ProtocolRedis Protocol = "redis"
	ProtocolCron  Protocol = "cron"
)

func (ct Protocol) IsValid() bool {
	switch ct {
	case ProtocolOgen, ProtocolHTTP, ProtocolKafka, ProtocolRedis, ProtocolGRPC, ProtocolCron:
		return true
	default:
		return false
	}
}
