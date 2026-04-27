package model

type Builder struct {
	name        string
	module      string
	entrypoints Entrypoints
	components  Components
}

func NewBuilder(name, module string) *Builder {
	return &Builder{
		name:   name,
		module: module,
	}
}

func (b *Builder) WithEntrypoints(entrypoints Entrypoints) *Builder {
	b.entrypoints = append(b.entrypoints, entrypoints...)

	return b
}

func (b *Builder) WithComponents(components Components) *Builder {
	b.components = append(b.components, components...)

	return b
}

func (b *Builder) Build() (*App, error) {
	return NewApp(b.name, b.module, b.entrypoints, b.components)
}
