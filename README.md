# Go Spreadsheet CLI

A lightweight terminal-based spreadsheet written in Go. Cells are addressed by 1-indexed `(row, col)` coordinates, stored in memory in a thread-safe sparse grid, and persisted as JSON checkpoints.

## Features

- **Cell CRUD:** `set`, `get`, `delete` (`del`)
- **Aggregations:** `sum`, `average` (`avg`), `count`, `min`, `max` — over a range or whole grid
- **View:** `print` a range or the whole grid (`.` = empty cell)
- **Edit blocks:** `copy` / `paste`, `clear`, `fill`, `find`, `sort` (single row/column, `asc`/`desc`)
- **Persistence:** `save` / `load` JSON checkpoints (atomic write via temp file + rename)
- **Robust CLI:** case-insensitive commands, friendly `help`, no panics (recovered and reported)

## Project Structure

```
.
├── cmd/sheet/main.go              # entrypoint, loads data/spreadsheet.json, starts REPL
├── internal/
│   ├── cli/
│   │   ├── cli.go                 # App, Execute(), Run() REPL loop
│   │   └── commands.go            # all command handlers + help text
│   ├── parser/command.go          # splits input line -> {Name, Args}
│   ├── spreadsheet/
│   │   ├── cell.go                # Cell{Value float64}
│   │   ├── grid.go                # thread-safe sparse Grid map[Position]Cell
│   │   └── operations.go          # sum/avg/count/min/max/copy/paste/find/clear/fill/sort
│   └── storage/
│       ├── storage.go             # Document/CellRecord types, Storage interface
│       └── json.go                # JSONStorage implementation
├── data/spreadsheet.json          # default checkpoint
└── go.mod                         # module myexcel, go 1.22
```

Key design points:

- `spreadsheet.Grid` is a sparse `map[Position]Cell` guarded by `sync.RWMutex`.
- `parser` does no validation — arg count/type checking lives in `cli` handlers.
- `storage` only deals with serializable `Document`; conversion via `ToDocument` / `ToGridCells`.
- Coordinates are 1-indexed and must be positive (`row >= 1, col >= 1`). Ranges are auto-normalized, so `(2,2)-(1,1)` works.

## Prerequisites

- Go 1.22+

## Installation & Run

```bash
# from project root
go run ./cmd/sheet

# or build a binary
go build -o sheet ./cmd/sheet
./sheet
```

On startup the app auto-loads `data/spreadsheet.json` if it exists:

```
Loaded checkpoint from data/spreadsheet.json (8 cell(s)).
Go Spreadsheet CLI — type 'help' for commands, 'exit' to quit.
```

## Usage

Type commands at the `> ` prompt. Commands are case-insensitive, extra whitespace is ignored.

### Command Reference

| Command           | Usage                                  | Description                                       |
| ----------------- | -------------------------------------- | ------------------------------------------------- |
| `set`             | `set <row> <col> <value>`              | Set cell value                                    |
| `get`             | `get <row> <col>`                      | Read cell value                                   |
| `delete` / `del`  | `delete <row> <col>`                   | Clear a single cell                               |
| `sum`             | `sum <r1> <c1> <r2> <c2>`              | Sum a range (empty cells = 0)                     |
| `average` / `avg` | `average <r1> <c1> <r2> <c2>`          | Average populated cells in range                  |
| `count`           | `count [<r1> <c1> <r2> <c2>]`          | Count populated cells (range or whole grid)       |
| `min`             | `min [<r1> <c1> <r2> <c2>]`            | Minimum (range or whole grid)                     |
| `max`             | `max [<r1> <c1> <r2> <c2>]`            | Maximum (range or whole grid)                     |
| `print`           | `print [<r1> <c1> <r2> <c2>]`          | Print range or whole grid (limit 100,000 cells)   |
| `copy`            | `copy <r1> <c1> <r2> <c2>`             | Copy range to clipboard                           |
| `paste`           | `paste <row> <col>`                    | Paste clipboard at position                       |
| `find`            | `find <value>`                         | List `(row, col)` positions equal to value        |
| `clear`           | `clear <r1> <c1> <r2> <c2>`            | Clear a range                                     |
| `fill`            | `fill <r1> <c1> <r2> <c2> <value>`     | Fill range with value                             |
| `sort`            | `sort <r1> <c1> <r2> <c2> [asc\|desc]` | Sort a single row/column (default `asc`)          |
| `save`            | `save [path]`                          | Save checkpoint (default `data/spreadsheet.json`) |
| `load`            | `load [path]`                          | Load checkpoint (replaces grid)                   |
| `help`            | `help`                                 | Show help                                         |
| `exit` / `quit`   | `exit`                                 | Quit                                              |

### Example Session

```
> set 1 1 10
Set Successfully.
> set 1 2 20
Set Successfully.
> set 2 1 30
Set Successfully.
> set 2 2 40
Set Successfully.
> get 1 1
10
> sum 1 1 2 2
100
> average 1 1 2 2
25
> count 1 1 2 2
4
> min 1 1 2 2
10
> max
40
> print 1 1 2 2
10        20
30        40
> copy 1 1 1 2
Copied Successfully.
> paste 3 1
Pasted Successfully.
> find 10
(1, 1)=10
(3, 1)=10
> fill 5 5 6 6 99
Filled Successfully.
> sort 1 1 1 2 desc
Sorted Successfully.
> save
Saved Successfully.
> clear 1 1 2 2
Clear Successfully.
> load
Loaded Successfully.
> exit
```

### Notes & Limits

- All cell values are `float64`. Whole numbers print without decimals.
- `average`/`min`/`max` on an empty range or grid return an error.
- `sort` only supports a single row or single column range.
- `print` without args prints the bounding box of populated cells, or `(empty grid)`.
- `save` creates parent directories as needed and writes atomically (`path.tmp` → rename).
- Unknown commands and bad arguments print `Error: ...` without exiting; panics are recovered.

## Data Format

`data/spreadsheet.json`:

```json
{
	"cells": [
		{ "row": 1, "col": 1, "value": 10 },
		{ "row": 1, "col": 2, "value": 20 }
	]
}
```

## Development

```bash
go vet ./...
go build ./...
```
