package ogen

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
		name         string
		namespace    string
		config       model.Config
		dependencies model.Dependencies
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
				name:      "Example",
				namespace: "Example",
				config:    model.Config{},
				dependencies: model.Dependencies{
					Must(model.NewDependency("cache", "drivers.redis")),
					Must(model.NewDependency("client", "drivers.http.client")),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := model.NewOgenController(
				tt.args.name,
				tt.args.namespace,
				"openapi.yaml",
				tt.args.config,
				tt.args.dependencies,
			)
			if err != nil {
				t.Fatal(err)
			}

			got, err := New(controller).Generate()
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
