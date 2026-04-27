package generate

import (
	"context"

	"github.com/vladimirkorzhenevskiy/gogen/internal/generator"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/dto"
)

type FS interface {
	Write(ctx context.Context, files ...model.File) error
}

type UseCase struct {
	fs FS
}

func New(fs FS) *UseCase {
	return &UseCase{
		fs: fs,
	}
}

func (uc *UseCase) Execute(ctx context.Context, appDTO *dto.App) error {
	app, err := appDTO.ToModel()
	if err != nil {
		return err
	}

	return uc.generate(ctx, app)
}

func (uc *UseCase) generate(ctx context.Context, app *model.App) error {
	pipe := generator.Pipeline{
		generator.Gitignore(app),
		generator.GoMod(app),

		generator.Docker(app),
		generator.DockerCompose(app),

		generator.Application(app),
	}

	files, err := pipe.Generate()
	if err != nil {
		return err
	}

	return uc.fs.Write(ctx, files...)
}
