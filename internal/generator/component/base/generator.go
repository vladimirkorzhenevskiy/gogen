package base

import (
	"embed"
	"io/fs"

	"github.com/vladimirkorzhenevskiy/gogen/internal/template"
)

//go:embed templates/*
var templates embed.FS

func RootFS() fs.FS {
	return template.MustSub(templates, "templates")
}

func Pipeline(fs embed.FS) *template.Pipeline {
	return template.NewPipeline(RootFS(), template.MustSub(fs, "templates"))
}
