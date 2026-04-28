package template

import (
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vladimirkorzhenevskiy/gogen/internal/model"
)

type Pipeline struct {
	filesystems []fs.FS
	functions   template.FuncMap
	tasks       []Task
}

func NewPipeline(filesystems ...fs.FS) *Pipeline {
	functions := template.FuncMap{
		"ToLower":          strings.ToLower,
		"ToUpper":          strings.ToUpper,
		"ToCamel":          strcase.ToCamel,
		"ToLowerCamel":     strcase.ToLowerCamel,
		"ToSnake":          strcase.ToSnake,
		"ToScreamingSnake": strcase.ToScreamingSnake,
	}

	return &Pipeline{
		filesystems: filesystems,
		functions:   functions,
	}
}

func (p *Pipeline) With(tasks ...Task) *Pipeline {
	p.tasks = append(p.tasks, tasks...)

	return p
}

func (p *Pipeline) WithFuncs(functions template.FuncMap) *Pipeline {
	for k, v := range functions {
		p.functions[k] = v
	}

	return p
}

func (p *Pipeline) WithFileSystems(filesystems ...fs.FS) *Pipeline {
	p.filesystems = append(p.filesystems, filesystems...)

	return p
}

func (p *Pipeline) Execute() ([]model.File, error) {
	tpl := template.New("").Funcs(p.functions)

	for i, filesystem := range p.filesystems {
		var err error

		tpl, err = tpl.ParseFS(filesystem, "*.tmpl")
		if err != nil {
			return nil, fmt.Errorf("failed to parse FS at index %d: %w", i, err)
		}
	}

	var files []model.File

	for _, task := range p.tasks {
		file, err := task.Execute(tpl)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (p *Pipeline) Generate() ([]model.File, error) {
	return p.Execute()
}
