package spreadsheet

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrInvalidPosition = errors.New("row and column must be positive integers")
	ErrCellEmpty       = errors.New("cell is empty")
	ErrEmptyGrid       = errors.New("no populated cells to operate on")
)

type Position struct {
	Row int
	Col int
}

type Grid struct {
	mu    sync.RWMutex
	cells map[Position]Cell
}

func NewGrid() *Grid {
	return &Grid{cells: make(map[Position]Cell)}
}

func validatePosition(row, col int) error {
	if row <= 0 || col <= 0 {
		return fmt.Errorf("%w: got row=%d col=%d", ErrInvalidPosition, row, col)
	}
	return nil
}

func (g *Grid) Set(row, col int, value float64) error {
	if err := validatePosition(row, col); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cells[Position{row, col}] = Cell{Value: value}
	return nil
}

func (g *Grid) Get(row, col int) (float64, error) {
	if err := validatePosition(row, col); err != nil {
		return 0, err
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	c, ok := g.cells[Position{row, col}]
	if !ok {
		return 0, fmt.Errorf("cell (%d,%d): %w", row, col, ErrCellEmpty)
	}
	return c.Value, nil
}

func (g *Grid) Has(row, col int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.cells[Position{row, col}]
	return ok
}

func (g *Grid) Delete(row, col int) error {
	if err := validatePosition(row, col); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.cells[Position{row, col}]; !ok {
		return fmt.Errorf("cell (%d,%d): %w", row, col, ErrCellEmpty)
	}
	delete(g.cells, Position{row, col})
	return nil
}

func (g *Grid) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.cells)
}

func (g *Grid) Bounds() (minRow, minCol, maxRow, maxCol int, ok bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	first := true
	for p := range g.cells {
		if first {
			minRow, maxRow = p.Row, p.Row
			minCol, maxCol = p.Col, p.Col
			first = false
			continue
		}
		if p.Row < minRow {
			minRow = p.Row
		}
		if p.Row > maxRow {
			maxRow = p.Row
		}
		if p.Col < minCol {
			minCol = p.Col
		}
		if p.Col > maxCol {
			maxCol = p.Col
		}
	}
	return minRow, minCol, maxRow, maxCol, !first
}

func (g *Grid) Snapshot() map[Position]Cell {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make(map[Position]Cell, len(g.cells))
	for p, c := range g.cells {
		out[p] = c
	}
	return out
}

func (g *Grid) Restore(cells map[Position]Cell) {
	g.mu.Lock()
	defer g.mu.Unlock()
	newCells := make(map[Position]Cell, len(cells))
	for p, c := range cells {
		if p.Row <= 0 || p.Col <= 0 {
			continue
		}
		newCells[p] = c
	}
	g.cells = newCells
}
