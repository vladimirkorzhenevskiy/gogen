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

func (uc *UseCase) Execute(ctx context.Context, inputDTO *dto.App) error {
	app, err := inputDTO.ToModel()
	if err != nil {
		return err
	}

	return uc.generate(ctx, app)
}

func (uc *UseCase) generate(ctx context.Context, app *model.App) error {
	files, err := generator.App(app).Generate()
	if err != nil {
		return err
	}

	return uc.fs.Write(ctx, files...)
}
