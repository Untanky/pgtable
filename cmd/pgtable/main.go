package main

import (
	"fmt"
	"os"

	"github.com/untanky/pgtable"
	"github.com/untanky/pgtable/internal/parser"
)

func main() {
	table, err := parser.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := pgtable.Render(os.Stdout, table); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
