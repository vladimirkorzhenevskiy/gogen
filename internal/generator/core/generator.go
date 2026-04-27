package core

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
