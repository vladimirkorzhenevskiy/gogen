package model

type Hook string

const (
	HookOnBeforeStart Hook = "OnBeforeStart"
	HookOnAfterStart  Hook = "OnAfterStart"
	HookOnBeforeStop  Hook = "OnBeforeStop"
	HookOnAfterStop   Hook = "OnAfterStop"
)

func (h Hook) String() string {
	return string(h)
}

func (h Hook) IsValid() bool {
	switch h {
	case HookOnBeforeStart, HookOnAfterStart, HookOnBeforeStop, HookOnAfterStop:
		return true
	default:
		return false
	}
}

func (h Hook) Values() []Hook {
	return []Hook{
		HookOnBeforeStart,
		HookOnAfterStart,
		HookOnBeforeStop,
		HookOnAfterStop,
	}
}
