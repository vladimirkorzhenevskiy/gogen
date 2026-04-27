package model

type Layer string

const (
	DriverLayer     Layer = "driver"
	ControllerLayer Layer = "controller"
	MiddlewareLayer Layer = "middleware"
	UseCaseLayer    Layer = "useCase"
	AdapterLayer    Layer = "adapter"
	RepositoryLayer Layer = "repository"
)

func (l Layer) String() string {
	return string(l)
}

func (l Layer) Group() string {
	switch l {
	case ControllerLayer:
		return "controllers"
	case UseCaseLayer:
		return "useCases"
	case RepositoryLayer:
		return "repositories"
	case DriverLayer:
		return "drivers"
	case MiddlewareLayer:
		return "middlewares"
	case AdapterLayer:
		return "adapters"
	default:
		return string(l) + "s"
	}
}

func (l Layer) IsValid() bool {
	switch l {
	case DriverLayer, ControllerLayer, MiddlewareLayer, UseCaseLayer, AdapterLayer, RepositoryLayer:
		return true
	default:
		return false
	}
}
