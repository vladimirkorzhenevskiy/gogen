package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/dto"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

type GenerateUseCase interface {
	Execute(ctx context.Context, app *dto.App) error
}

type Controller struct {
	generateUseCase GenerateUseCase
}

func New(generateUseCase GenerateUseCase) *Controller {
	return &Controller{
		generateUseCase: generateUseCase,
	}
}

func (c *Controller) Generate(ctx context.Context, cmd *cli.Command) error {
	file, err := os.Open("gogen.yaml")
	if err != nil {
		return fmt.Errorf("failed to open gogen.yaml: %w", err)
	}
	defer func() { _ = file.Close() }()

	var app *dto.App

	if err := yaml.NewDecoder(file).Decode(&app); err != nil {
		return fmt.Errorf("failed to decode gogen.yaml: %w", err)
	}

	if err := validator.New().Struct(app); err != nil {
		return fmt.Errorf("failed to validate gogen.yaml: %w", err)
	}

	return c.generateUseCase.Execute(ctx, app)
}

func (c *Controller) Run(ctx context.Context) error {
	return c.Root().Run(ctx, os.Args)
}

func (c *Controller) Root() *cli.Command {
	return &cli.Command{
		Name:  "github.com/vladimirkorzhenevskiy/gogen",
		Usage: "github.com/vladimirkorzhenevskiy/gogen [command]",
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "generate [command options]",
				Action: c.Generate,
			},
		},
	}
}
