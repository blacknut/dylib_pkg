package main

import (
	"fmt"
	"os/exec"
)

func installNameChangeLib(lib *dylib) error {
	for _, dep := range lib.deps {
		if err := installNameChange(lib.path, dep); err != nil {
			return err
		}
	}
	return nil
}

func installNameChange(path string, lib *dylib) error {
	newPath := fmt.Sprintf("@executable_path/%s", lib.name)

	if flagVerbose {
		fmt.Printf("Fixing file %s for lib: %s\n", path, lib.name)
	}

	// fmt.Println("install_name_tool", "-change", lib.path, newPath, path)

	// install_name_tool -change $dylib @executable_path/`basename $dylib` ./ga-client
	if out, err := exec.Command("install_name_tool", "-change", lib.path, newPath, path).CombinedOutput(); err != nil {
		fmt.Println("Failed to fix lib:", path, string(out))
		return err
	}

	return nil
}
