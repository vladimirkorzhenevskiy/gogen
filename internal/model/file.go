package model

import "path/filepath"

type File struct {
	path      string
	data      []byte
	overwrite bool
}

func NewFile(path string, data []byte, overwrite bool) File {
	return File{
		path:      filepath.Clean(path),
		data:      data,
		overwrite: overwrite,
	}
}

func (f File) Path() string {
	return f.path
}

func (f File) Dir() string {
	return filepath.Dir(f.path)
}

func (f File) Name() string {
	return filepath.Base(f.path)
}

func (f File) Data() []byte {
	return f.data
}

func (f File) Overwrite() bool {
	return f.overwrite
}
