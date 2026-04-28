package nop

import (
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

func (Generator) Generate() ([]model.File, error) {
	return []model.File{}, nil
}
