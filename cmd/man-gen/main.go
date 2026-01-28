package main

import (
	"log"

	"github.com/jmpargana/gq/internal/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	var header = &doc.GenManHeader{
		Title:   "gq",
		Section: "1",
		Manual:  "User Commands",
	}
	err := doc.GenManTree(cmd.RootCmd, header, ".")
	if err != nil {
		log.Fatalln(err)
	}
}
