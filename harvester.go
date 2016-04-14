package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

type harvester struct {
	// executable path
	root string

	// executable dependencies names
	deps map[string]bool

	// all libs
	libs map[string]*dylib

	// current lib stack
	libStack []*dylib
}

func newHarvester(root string) *harvester {
	return &harvester{
		root:     root,
		deps:     make(map[string]bool),
		libs:     make(map[string]*dylib),
		libStack: []*dylib{},
	}
}

func (h *harvester) collect() error {
	return h.collectFile(h.root, true)
}

func (h *harvester) collectFile(path string, isRoot bool) error {
	if flagVerbose {
		fmt.Println("Collecting:", path)
	}

	out, err := oToolExec(path)
	if err != nil {
		fmt.Println("Failed to exec otool:", err, string(out))
		return err
	}

	libs, err := oToolParse(out)
	if err != nil {
		return err
	}

	for _, lib := range libs {
		if curLib := h.curLib(); curLib != nil {
			curLib.addDep(lib)
		}

		if lib.path == path {
			continue
		}

		// root dependencies
		h.deps[lib.name] = true

		if h.libs[lib.name] == nil {
			// all libs
			h.libs[lib.name] = lib

			if err := h.collectLib(lib); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *harvester) collectLib(lib *dylib) error {
	h.pushLib(lib)

	if err := h.collectFile(lib.path, false); err != nil {
		return err
	}

	// pop
	h.popLib()

	return nil
}

func (h *harvester) pushLib(lib *dylib) {
	h.libStack = append(h.libStack, lib)
}

func (h *harvester) popLib() *dylib {
	var result *dylib

	result, h.libStack = h.libStack[len(h.libStack)-1], h.libStack[:len(h.libStack)-1]

	return result
}

func (h *harvester) curLib() *dylib {
	if len(h.libStack) == 0 {
		return nil
	}

	return h.libStack[len(h.libStack)-1]
}

func (h *harvester) destDir() string {
	if flagDest != "" {
		return flagDest
	}

	return path.Dir(h.root)
}

func (h *harvester) copy() error {
	for _, lib := range h.libs {
		dest := h.destLibPath(lib)

		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			if !flagForce {
				return fmt.Errorf("-force flag not set and destination file already exists: %s", dest)
			}
		}

		if err := copyFile(lib.path, dest); err != nil {
			return err
		}
	}

	return nil
}

func (h *harvester) destLibPath(lib *dylib) string {
	return path.Join(h.destDir(), lib.name)
}

func (h *harvester) fixReferences() error {
	// fixes executable
	for name := range h.deps {
		if err := installNameChange(h.root, h.libs[name]); err != nil {
			return err
		}
	}

	// fixes libs
	for _, lib := range h.libs {
		destLib := h.destLibPath(lib)

		for _, dep := range lib.deps {
			if err := installNameChange(destLib, dep); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(from string, to string) error {
	if flagVerbose {
		fmt.Printf("Copying: %s => %s\n", from, to)
	}

	// open source file
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	// open destination file
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
