package root

import (
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/root/gitignore"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/root/golangci"
	"github.com/vladimirkorzhenevskiy/gogen/internal/generator/root/gomod"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

func Gitignore(app *model.App) gitignore.Generator {
	return gitignore.New(app)
}

func GolangCI(app *model.App) golangci.Generator {
	return golangci.New(app)
}

func GoMod(app *model.App) gomod.Generator {
	return gomod.New(app)
}
