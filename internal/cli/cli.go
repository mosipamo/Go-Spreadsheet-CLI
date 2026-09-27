package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"myexcel/internal/parser"
	"myexcel/internal/spreadsheet"
	"myexcel/internal/storage"
)

type App struct {
	Grid         *spreadsheet.Grid
	Store        storage.Storage
	DefaultPath  string
	Clipboard    spreadsheet.Clipboard
	HasClipboard bool
	Out          io.Writer
}

type handlerFunc func(*App, []string) error

var commands = map[string]handlerFunc{
	"set":     cmdSet,
	"get":     cmdGet,
	"delete":  cmdDelete,
	"del":     cmdDelete,
	"sum":     cmdSum,
	"average": cmdAverage,
	"avg":     cmdAverage,
	"count":   cmdCount,
	"min":     cmdMin,
	"max":     cmdMax,
	"print":   cmdPrint,
	"copy":    cmdCopy,
	"paste":   cmdPaste,
	"find":    cmdFind,
	"clear":   cmdClear,
	"fill":    cmdFill,
	"sort":    cmdSort,
	"save":    cmdSave,
	"load":    cmdLoad,
	"help":    cmdHelp,
	"exit":    cmdExit,
	"quit":    cmdExit,
}

func (app *App) Execute(line string) (exit bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			exit = false
			err = fmt.Errorf("internal error recovered: %v", r)
		}
	}()

	cmd, perr := parser.Parse(line)
	if perr != nil {
		return false, nil
	}

	handler, ok := commands[cmd.Name]
	if !ok {
		return false, fmt.Errorf("unknown command %q (try 'help')", cmd.Name)
	}

	if runErr := handler(app, cmd.Args); runErr != nil {
		if errors.Is(runErr, ErrExit) {
			return true, nil
		}
		return false, runErr
	}
	return false, nil
}

func (app *App) Run(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for {
		fmt.Fprint(app.Out, "\n> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		exit, err := app.Execute(line)
		if err != nil {
			fmt.Fprintf(app.Out, "Error: %v\n", err)
			continue
		}
		if exit {
			break
		}
	}
	return scanner.Err()
}
