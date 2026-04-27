package fs

import (
	"context"
	"fmt"
	"os"

	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Adapter struct{}

func New() *Adapter {
	return &Adapter{}
}

func (a Adapter) Write(_ context.Context, files ...model.File) error {
	for _, file := range files {

		fmt.Println(file.Path())
		if !file.Overwrite() {
			if _, err := os.Stat(file.Path()); err == nil {
				continue
			}
		}

		if err := os.MkdirAll(file.Dir(), 0755); err != nil {
			return fmt.Errorf("failed to create dir %s: %w", file.Dir(), err)
		}

		if err := os.WriteFile(file.Path(), file.Data(), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", file.Path(), err)
		}
	}

	return nil
}
