/*
Copyright © 2026 Joao Pargana jmpargana@gmail.com
*/
package main

import (
	"os"

	"github.com/jmpargana/gq/internal/cmd"
)

func Execute() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
