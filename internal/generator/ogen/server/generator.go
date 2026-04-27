package server

import (
	"os"
	"path"

	"github.com/ogen-go/ogen"
	"github.com/ogen-go/ogen/gen"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Generator struct {
	spec string
}

func New(spec string) Generator {
	return Generator{
		spec: spec,
	}
}

func (g Generator) Generate() ([]model.File, error) {
	data, err := os.ReadFile(g.spec)
	if err != nil {
		return nil, err
	}

	spec, err := ogen.Parse(data)
	if err != nil {
		return nil, err
	}

	generator, err := gen.NewGenerator(spec, gen.Options{
		Generator: gen.GenerateOptions{
			Features: &gen.FeatureOptions{
				Enable: gen.FeatureSet{
					gen.PathsServer.Name: struct{}{},
				},
				DisableAll: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	fs := &fileSystem{}

	if err = generator.WriteSource(fs, "server"); err != nil {
		return nil, err
	}

	return fs.files, nil
}

func Generate(spec string) ([]model.File, error) {
	return New(spec).Generate()
}

type fileSystem struct {
	files []model.File
}

func (fs *fileSystem) WriteFile(baseName string, content []byte) error {
	file := model.NewFile(
		path.Join("api/server", baseName),
		content,
		true,
	)

	fs.files = append(fs.files, file)

	return nil
}
