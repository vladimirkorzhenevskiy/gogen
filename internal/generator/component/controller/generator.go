package controller

import (
	"fmt"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/controller/common"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/component/controller/ogen"
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
	fmt.Println("Generating controller")
	switch g.controller.Protocol() {
	case model.ProtocolOgen:
		fmt.Println("Generating ogen")
		return ogen.New(g.app, g.controller).Generate()
	default:
		fmt.Println("Generating default")
		return common.New(g.controller).Generate()
	}
}
