package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"myexcel/internal/storage"
)

var ErrExit = errors.New("exit requested")

func parseInt(label, s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number, got %q", label, s)
	}
	return n, nil
}

func parseFloat(label, s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number, got %q", label, s)
	}
	return f, nil
}

func requireArgs(name string, args []string, want int, usage string) error {
	if len(args) != want {
		return fmt.Errorf("usage: %s %s", name, usage)
	}
	return nil
}

func parseRowCol(args []string) (int, int, error) {
	row, err := parseInt("row", args[0])
	if err != nil {
		return 0, 0, err
	}
	col, err := parseInt("col", args[1])
	if err != nil {
		return 0, 0, err
	}
	return row, col, nil
}

func parseRange(args []string) (int, int, int, int, error) {
	r1, err := parseInt("row1", args[0])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	c1, err := parseInt("col1", args[1])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	r2, err := parseInt("row2", args[2])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	c2, err := parseInt("col2", args[3])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return r1, c1, r2, c2, nil
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func cmdSet(app *App, args []string) error {
	if err := requireArgs("set", args, 3, "<row> <col> <value>"); err != nil {
		return err
	}
	row, col, err := parseRowCol(args[:2])
	if err != nil {
		return err
	}
	value, err := parseFloat("value", args[2])
	if err != nil {
		return err
	}
	if err := app.Grid.Set(row, col, value); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Set Successfully.")
	return nil
}

func cmdGet(app *App, args []string) error {
	if err := requireArgs("get", args, 2, "<row> <col>"); err != nil {
		return err
	}
	row, col, err := parseRowCol(args)
	if err != nil {
		return err
	}
	v, err := app.Grid.Get(row, col)
	if err != nil {
		return err
	}
	fmt.Fprintln(app.Out, formatNumber(v))
	return nil
}

func cmdDelete(app *App, args []string) error {
	if err := requireArgs("delete", args, 2, "<row> <col>"); err != nil {
		return err
	}
	row, col, err := parseRowCol(args)
	if err != nil {
		return err
	}
	if err := app.Grid.Delete(row, col); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Deleted Successfully.")
	return nil
}

func cmdSum(app *App, args []string) error {
	if err := requireArgs("sum", args, 4, "<row1> <col1> <row2> <col2>"); err != nil {
		return err
	}
	r1, c1, r2, c2, err := parseRange(args)
	if err != nil {
		return err
	}
	total, err := app.Grid.Sum(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	fmt.Fprintln(app.Out, formatNumber(total))
	return nil
}

func cmdAverage(app *App, args []string) error {
	if err := requireArgs("average", args, 4, "<row1> <col1> <row2> <col2>"); err != nil {
		return err
	}
	r1, c1, r2, c2, err := parseRange(args)
	if err != nil {
		return err
	}
	avg, err := app.Grid.Average(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	fmt.Fprintln(app.Out, formatNumber(avg))
	return nil
}

func cmdCount(app *App, args []string) error {
	switch len(args) {
	case 0:
		fmt.Fprintln(app.Out, app.Grid.CountAll())
		return nil
	case 4:
		r1, c1, r2, c2, err := parseRange(args)
		if err != nil {
			return err
		}
		n, err := app.Grid.CountRange(r1, c1, r2, c2)
		if err != nil {
			return err
		}
		fmt.Fprintln(app.Out, n)
		return nil
	default:
		return fmt.Errorf("usage: count [<row1> <col1> <row2> <col2>]")
	}
}

func cmdMin(app *App, args []string) error {
	switch len(args) {
	case 0:
		v, err := app.Grid.MinAll()
		if err != nil {
			return err
		}
		fmt.Fprintln(app.Out, formatNumber(v))
		return nil
	case 4:
		r1, c1, r2, c2, err := parseRange(args)
		if err != nil {
			return err
		}
		v, err := app.Grid.MinRange(r1, c1, r2, c2)
		if err != nil {
			return err
		}
		fmt.Fprintln(app.Out, formatNumber(v))
		return nil
	default:
		return fmt.Errorf("usage: min [<row1> <col1> <row2> <col2>]")
	}
}

func cmdMax(app *App, args []string) error {
	switch len(args) {
	case 0:
		v, err := app.Grid.MaxAll()
		if err != nil {
			return err
		}
		fmt.Fprintln(app.Out, formatNumber(v))
		return nil
	case 4:
		r1, c1, r2, c2, err := parseRange(args)
		if err != nil {
			return err
		}
		v, err := app.Grid.MaxRange(r1, c1, r2, c2)
		if err != nil {
			return err
		}
		fmt.Fprintln(app.Out, formatNumber(v))
		return nil
	default:
		return fmt.Errorf("usage: max [<row1> <col1> <row2> <col2>]")
	}
}

func cmdPrint(app *App, args []string) error {
	var r1, c1, r2, c2 int
	switch len(args) {
	case 0:
		minRow, minCol, maxRow, maxCol, ok := app.Grid.Bounds()
		if !ok {
			fmt.Fprintln(app.Out, "(empty grid)")
			return nil
		}
		r1, c1, r2, c2 = minRow, minCol, maxRow, maxCol
	case 4:
		var err error
		r1, c1, r2, c2, err = parseRange(args)
		if err != nil {
			return err
		}
		if r1 > r2 {
			r1, r2 = r2, r1
		}
		if c1 > c2 {
			c1, c2 = c2, c1
		}
	default:
		return fmt.Errorf("usage: print [<row1> <col1> <row2> <col2>]")
	}

	const maxCells = 100_000
	if (r2-r1+1)*(c2-c1+1) > maxCells {
		return fmt.Errorf("range too large to print (%d cells, limit %d)", (r2-r1+1)*(c2-c1+1), maxCells)
	}

	for row := r1; row <= r2; row++ {
		var line strings.Builder
		for col := c1; col <= c2; col++ {
			cell := "."
			if v, err := app.Grid.Get(row, col); err == nil {
				cell = formatNumber(v)
			}
			fmt.Fprintf(&line, "%-10s", cell)
		}
		fmt.Fprintln(app.Out, strings.TrimRight(line.String(), " "))
	}
	return nil
}

func cmdCopy(app *App, args []string) error {
	if err := requireArgs("copy", args, 4, "<row1> <col1> <row2> <col2>"); err != nil {
		return err
	}
	r1, c1, r2, c2, err := parseRange(args)
	if err != nil {
		return err
	}
	block, err := app.Grid.Copy(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	app.Clipboard = block
	app.HasClipboard = true
	fmt.Fprintln(app.Out, "Copied Successfully.")
	return nil
}

func cmdPaste(app *App, args []string) error {
	if err := requireArgs("paste", args, 2, "<row> <col>"); err != nil {
		return err
	}
	row, col, err := parseRowCol(args)
	if err != nil {
		return err
	}
	if !app.HasClipboard {
		return fmt.Errorf("clipboard is empty; use 'copy' first")
	}
	if err := app.Grid.Paste(row, col, app.Clipboard); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Pasted Successfully.")
	return nil
}

func cmdFind(app *App, args []string) error {
	if err := requireArgs("find", args, 1, "<value>"); err != nil {
		return err
	}
	value, err := parseFloat("value", args[0])
	if err != nil {
		return err
	}
	positions := app.Grid.Find(value)
	if len(positions) == 0 {
		fmt.Fprintln(app.Out, "(no matches)")
		return nil
	}
	for _, p := range positions {
		fmt.Fprintf(app.Out, "(%d, %d)=%g\n", p.Row, p.Col, value)
	}
	return nil
}

func cmdClear(app *App, args []string) error {
	if err := requireArgs("clear", args, 4, "<r1> <c1> <r2> <c2>"); err != nil {
		return err
	}
	r1, c1, r2, c2, err := parseRange(args)
	if err != nil {
		return err
	}
	if err := app.Grid.Clear(r1, c1, r2, c2); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Cleared Successfully.")
	return nil
}

func cmdFill(app *App, args []string) error {
	if err := requireArgs("fill", args, 5, "<r1> <c1> <r2> <c2> <value>"); err != nil {
		return err
	}
	r1, c1, r2, c2, err := parseRange(args[:4])
	if err != nil {
		return err
	}
	value, err := parseFloat("value", args[4])
	if err != nil {
		return err
	}
	if err := app.Grid.Fill(r1, c1, r2, c2, value); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Filled Successfully.")
	return nil
}

func cmdSort(app *App, args []string) error {
	if len(args) != 4 && len(args) != 5 {
		return fmt.Errorf("usage: sort <row1> <col1> <row2> <col2> [asc|desc]")
	}
	r1, c1, r2, c2, err := parseRange(args[:4])
	if err != nil {
		return err
	}
	ascending := true
	if len(args) == 5 {
		switch strings.ToLower(args[4]) {
		case "asc":
			ascending = true
		case "desc":
			ascending = false
		default:
			return fmt.Errorf("sort direction must be 'asc' or 'desc', got %q", args[4])
		}
	}
	if err := app.Grid.Sort(r1, c1, r2, c2, ascending); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Sorted Successfully.")
	return nil
}

func cmdSave(app *App, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: save [path]")
	}
	path := app.DefaultPath
	if len(args) == 1 {
		path = args[0]
	}
	doc := storage.ToDocument(app.Grid.Snapshot())
	if err := app.Store.Save(path, doc); err != nil {
		return err
	}
	fmt.Fprintln(app.Out, "Saved Successfully.")
	return nil
}

func cmdLoad(app *App, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: load [path]")
	}
	path := app.DefaultPath
	if len(args) == 1 {
		path = args[0]
	}
	doc, err := app.Store.Load(path)
	if err != nil {
		return err
	}
	app.Grid.Restore(storage.ToGridCells(doc))
	fmt.Fprintln(app.Out, "Loaded Successfully.")
	return nil
}

func cmdHelp(app *App, args []string) error {
	fmt.Fprint(app.Out, `Commands:
  set <row> <col> <value>              set a cell
  get <row> <col>                      read a cell
  delete <row> <col>                   clear a single cell
  sum <r1> <c1> <r2> <c2>              sum a range
  average <r1> <c1> <r2> <c2>          average a range
  count [<r1> <c1> <r2> <c2>]          count populated cells (range or whole grid)
  min [<r1> <c1> <r2> <c2>]            minimum value (range or whole grid)
  max [<r1> <c1> <r2> <c2>]            maximum value (range or whole grid)
  print [<r1> <c1> <r2> <c2>]          print a range or the whole grid
  copy <r1> <c1> <r2> <c2>             copy a range to the clipboard
  paste <row> <col>                    paste the clipboard at a position
  find <value>                         list positions of cells equal to a value
  clear <r1> <c1> <r2> <c2>            clear a range
  fill <r1> <c1> <r2> <c2> <value>     fill a range with a value
  sort <r1> <c1> <r2> <c2> [asc|desc]  sort a single row/column range
  save [path]                          save a checkpoint (default: data/spreadsheet.json)
  load [path]                          load a checkpoint
  help                                 show this message
  exit | quit                          quit the program
`)
	return nil
}

func cmdExit(app *App, args []string) error {
	return ErrExit
}
