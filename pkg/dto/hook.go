package dto

import (
	"fmt"
	"strings"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/slicesutil"
)

type Hook string

func (h Hook) ToModel() (model.Hook, error) {
	switch strings.ToLower(string(h)) {
	case strings.ToLower(model.HookOnBeforeStart.String()):
		return model.HookOnBeforeStart, nil
	case strings.ToLower(model.HookOnAfterStart.String()):
		return model.HookOnAfterStart, nil
	case strings.ToLower(model.HookOnBeforeStop.String()):
		return model.HookOnBeforeStop, nil
	case strings.ToLower(model.HookOnAfterStop.String()):
		return model.HookOnAfterStop, nil
	default:
		return "", fmt.Errorf("invalid hook: %v", h)
	}
}

type Hooks []Hook

func (h Hooks) ToModel() ([]model.Hook, error) {
	return slicesutil.TryMap(h, Hook.ToModel)
}
