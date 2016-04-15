package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"path"
	"strings"
)

var oToolIgnorePrefixes = []string{"/System", "/usr/lib"}

func oToolExec(filePath string) ([]byte, error) {
	return exec.Command("otool", "-L", filePath).CombinedOutput()
}

func oToolParse(out []byte) ([]*dylib, error) {
	result := []*dylib{}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		parts := strings.SplitAfter(line, " ")
		if len(parts) < 2 {
			// permits to ignore first line
			continue
		}

		libPath := strings.TrimSpace(parts[0])
		if oToolIgnoreLibPath(libPath) {
			continue
		}

		lib := newDylib(path.Base(libPath), libPath)

		result = append(result, lib)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func oToolIgnoreLibPath(libPath string) bool {
	for _, prefix := range oToolIgnorePrefixes {
		if strings.HasPrefix(libPath, prefix) {
			return true
		}
	}
	return false
}
