package main

import (
	"fmt"
	"os"
	"path/filepath"
	"ws/cmd/root"
)

func main() {
	rootCmd := root.NewRootCommand(filepath.Base(os.Args[0]))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)

	}
}
