package main

import (
	"path"
	"strings"
)

const (
	executablePathPrefix = "@executable_path"
)

type dylib struct {
	name string
	path string

	deps []*dylib
}

func newDylib(name, path string) *dylib {
	return &dylib{
		name: name,
		path: path,
		deps: []*dylib{},
	}
}

func (d *dylib) addDep(lib *dylib) {
	d.deps = append(d.deps, lib)
}

func (d *dylib) isExecutablePath() bool {
	return strings.HasPrefix(d.path, executablePathPrefix)
}

func (d *dylib) absolutePath(execDir string) string {
	if d.isExecutablePath() {
		return path.Join(execDir, d.name)
	}
	return d.path
}
