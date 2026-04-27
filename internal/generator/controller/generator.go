package controller

import (
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/controller/ogen"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Generator struct {
	app        *model.App
	controller *model.Controller
}

func New(app *model.App, controller *model.Controller) Generator {
	return Generator{
		app:        app,
		controller: controller,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	switch g.controller.Protocol() {
	case model.ProtocolOgen:
		return ogen.New(g.app, g.controller).Generate()
	default:
		return []model.File{}, nil
	}
}
