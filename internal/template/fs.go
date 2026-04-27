package template

import (
	"io/fs"
)

func MustSub(fileSystem fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fileSystem, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
