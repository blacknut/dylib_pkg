package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	version = "0.0.1"
)

var (
	flagVersion bool
	flagVerbose bool
	flagDest    string
	flagForce   bool
	flagNoop    bool
)

func init() {
	flag.BoolVar(&flagVersion, "version", false, "Display version")
	flag.BoolVar(&flagVerbose, "verbose", false, "Display verbose infos")
	flag.StringVar(&flagDest, "dest", "", "Destination directory path (default: in the same dir as executable))")
	flag.BoolVar(&flagForce, "force", false, "Overwrite dylib files in destination directory")
	flag.BoolVar(&flagNoop, "noop", false, "Dry run")

	flag.Usage = func() {
		fmt.Println("Usage: $ dylib_pkg /path/to/executable")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	checkFlags()

	if flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	execPath := args[0]

	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", execPath)
		flag.Usage()
		os.Exit(1)
	}

	h := newHarvester(execPath)

	if err := h.collect(); err != nil {
		fmt.Println("Failed to collect dylibs:", err)
		os.Exit(1)
	}

	if flagVerbose {
		fmt.Printf("Found %d libs:\n", len(h.libs))
		h.print()
	}

	if err := h.copy(); err != nil {
		fmt.Println("Failed to copy dylibs:", err)
		os.Exit(1)
	}

	if err := h.fixReferences(); err != nil {
		fmt.Println("Failed to fix dylibs references:", err)
		os.Exit(1)
	}
}

func checkFlags() {
	if flagDest != "" {
		if _, err := os.Stat(flagDest); os.IsNotExist(err) {
			fmt.Println("Dest does not exist:", flagDest)
			flag.Usage()
			os.Exit(1)
		}
	}
}
