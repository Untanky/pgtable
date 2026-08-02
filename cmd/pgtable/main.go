package main

import (
	"fmt"
	"os"

	"github.com/untanky/pgtable/internal/parser"
	"github.com/untanky/pgtable/render"
)

func main() {
	table, err := parser.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	renderer := render.NewModel(table)

	if _, err := os.Stdout.WriteString(renderer.Render()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
