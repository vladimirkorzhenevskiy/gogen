package entrypoint

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func TestGenerator_Generate(t *testing.T) {
	type args struct {
		name        string
		module      string
		entrypoints model.Entrypoints
		components  model.Components
	}

	tests := []struct {
		name    string
		args    args
		want    []model.File
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should_successful",
			args: args{
				name:   "example",
				module: "example",
				components: model.Components{
					Must(model.NewOgenController(
						"api",
						"leaderboard",
						"openapi.yaml",
						model.Config{
							Must(model.NewSetting("host", "string", true, "localhost")),
							Must(model.NewSetting("port", "uint8", true, "8080")),
						},
						model.Dependencies{
							Must(model.NewDependency("getter", "usecases.leaderboard.get")),
						},
					)),
					Must(model.NewUseCase(
						"get",
						"leaderboard",
						model.Config{},
						model.Dependencies{
							Must(model.NewDependency("repository", "repositories.leaderboard.postgres")),
						},
					)),
					Must(model.NewRepository(
						"postgres",
						"leaderboard",
						model.Config{},
						model.Dependencies{
							Must(model.NewDependency("db", "drivers.postgres")),
						},
					)),
					Must(model.NewDriver(
						"postgres",
						"postgres",
					)),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := model.NewApp(
				tt.args.name,
				tt.args.module,
				tt.args.entrypoints,
				tt.args.components,
			)
			if err != nil {
				t.Fatal(err)
			}

			got, err := New(app, app.Entrypoints()[0]).Generate()
			if !tt.wantErr(t, err) {
				return
			}
			//assert.Equal(t, tt.want, got)

			for _, file := range got {
				f, err := os.Create("./example/" + file.Name()) // Creates or truncates the file
				if err != nil {
					log.Fatal(err)
				}

				f.Write(file.Data())
			}
		})
	}
}
