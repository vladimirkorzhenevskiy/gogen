package controller

import (
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/controller/ogen"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Generator struct {
	controller *model.Controller
}

func New(controller *model.Controller) Generator {
	return Generator{
		controller: controller,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	switch g.controller.Protocol() {
	case model.ProtocolOgen:
		return ogen.New(g.controller).Generate()
	default:
		return []model.File{}, nil
	}
}
