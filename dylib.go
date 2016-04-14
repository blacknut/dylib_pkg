package main

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
