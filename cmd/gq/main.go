/*
Copyright © 2026 Joao Pargana jmpargana@gmail.com
*/
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jmpargana/gq/internal/ast"
	json "github.com/jmpargana/gq/internal/gqjson"
	"github.com/jmpargana/gq/internal/lexer"
	"github.com/jmpargana/gq/internal/parser"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "gq",
	Short: "jq-like command-line tool written in go",
	Long:  `A fast, simple, and expressive way to query and transform JSON data from the command line`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Printf("Version: %s\n", version)
			os.Exit(0)
		}

		if err := requireStdin(); err != nil {
			return err
		}

		if len(os.Args) < 2 {
			return fmt.Errorf("no program provided")
		}

		r := bufio.NewReader(os.Stdin)
		tokens := lexer.Lex(os.Args[1])
		p := parser.NewParser(tokens)
		t := p.ParseExpr()
		obj := json.ParseObject(r)
		result := ast.Transform(obj, t)
		json.Print(result)
		
		return nil
	},
}

func requireStdin() error {
	stat, err := os.Stdin.Stat()
    if err != nil {
        return fmt.Errorf("failed to stat stdin: %w", err)
    }

    // If stdin is a terminal, nothing was piped
    if stat.Mode()&os.ModeCharDevice != 0 {
        return fmt.Errorf("no input provided on stdin")
    }

    return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "Display CLI version")
}

func main() {
	Execute()
}
