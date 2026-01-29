package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jmpargana/gq/internal/ast"
	json "github.com/jmpargana/gq/internal/gqjson"
	"github.com/jmpargana/gq/internal/lexer"
	"github.com/jmpargana/gq/internal/parser"
	"github.com/jmpargana/gq/internal/stream"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gq",
	Short: "jq-like command-line tool written in go",
	Long: `A fast, simple, and expressive way to query 
and transform JSON data from the command line.

Most of the jqlang syntax (https://jqlang.org/manual/) should be available.
List of battle tested features:
	- indexing numbers, strings, identifiers
	- indexing iterators (.[])
	- cartesian product of iterators
	- chaing indexes
	- array creation
	- dictionary creation
	- nested piping
	
Additionally, you can also view the AST of your jqlang expression.
`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireStdin(); err != nil {
			return err
		}

		if len(os.Args) < 2 {
			return fmt.Errorf("no program provided")
		}

		r := bufio.NewReader(os.Stdin)

		// TODO: replace with cobra args
		tokens := lexer.Lex(os.Args[1])
		for _, tok := range tokens {
			if tok.Kind == lexer.ILLEGAL {
				return fmt.Errorf("illegal token found: %v", tok.Value)
			}
		}

		p := parser.NewParser(tokens)

		t := p.ParseExpr()

		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			fmt.Println("Generated AST:")
			fmt.Println()
			fmt.Println(ast.PrintAST(t, 0))
		}

		obj := json.ParseObject(r)
		stream := stream.NewS(obj)

		result := ast.TransformStream(stream, t)
		fmt.Printf("%s", result.String())

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

func init() {
	RootCmd.Flags().BoolP("debug", "d", false, "Displays AST from requested expression")
}
