package main

import (
	"fmt"
	"os"

	"myexcel/internal/cli"
	"myexcel/internal/spreadsheet"
	"myexcel/internal/storage"
)

const defaultDataPath = "data/spreadsheet.json"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "fatal: recovered from panic: %v\n", r)
			os.Exit(1)
		}
	}()

	grid := spreadsheet.NewGrid()
	store := storage.NewJSONStorage()

	app := &cli.App{
		Grid:        grid,
		Store:       store,
		DefaultPath: defaultDataPath,
		Out:         os.Stdout,
	}

	if _, statErr := os.Stat(defaultDataPath); statErr == nil {
		doc, err := store.Load(defaultDataPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load checkpoint %q: %v\n", defaultDataPath, err)
		} else {
			grid.Restore(storage.ToGridCells(doc))
			fmt.Printf("Loaded checkpoint from %s (%d cell(s)).\n", defaultDataPath, grid.Len())
		}
	}

	fmt.Println("Go Spreadsheet CLI — type 'help' for commands, 'exit' to quit.")

	if err := app.Run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: reading input: %v\n", err)
		os.Exit(1)
	}
}
