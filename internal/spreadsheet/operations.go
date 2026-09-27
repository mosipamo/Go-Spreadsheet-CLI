package spreadsheet

import (
	"errors"
	"fmt"
	"sort"
)

func normalizeRange(r1, c1, r2, c2 int) (int, int, int, int, error) {
	if err := validatePosition(r1, c1); err != nil {
		return 0, 0, 0, 0, err
	}
	if err := validatePosition(r2, c2); err != nil {
		return 0, 0, 0, 0, err
	}
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	return r1, c1, r2, c2, nil
}

func (g *Grid) valuesInRange(r1, c1, r2, c2 int) ([]float64, error) {
	r1, c1, r2, c2, err := normalizeRange(r1, c1, r2, c2)
	if err != nil {
		return nil, err
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	var values []float64
	for p, c := range g.cells {
		if p.Row >= r1 && p.Row <= r2 && p.Col >= c1 && p.Col <= c2 {
			values = append(values, c.Value)
		}
	}
	return values, nil
}

func (g *Grid) allValues() []float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	values := make([]float64, 0, len(g.cells))
	for _, c := range g.cells {
		values = append(values, c.Value)
	}
	return values
}

func (g *Grid) Sum(r1, c1, r2, c2 int) (float64, error) {
	values, err := g.valuesInRange(r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, v := range values {
		total += v
	}
	return total, nil
}

func (g *Grid) Average(r1, c1, r2, c2 int) (float64, error) {
	values, err := g.valuesInRange(r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("average of range (%d,%d)-(%d,%d): %w", r1, c1, r2, c2, ErrEmptyGrid)
	}
	var total float64
	for _, v := range values {
		total += v
	}
	return total / float64(len(values)), nil
}

func (g *Grid) CountRange(r1, c1, r2, c2 int) (int, error) {
	values, err := g.valuesInRange(r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	return len(values), nil
}

func (g *Grid) CountAll() int {
	return g.Len()
}

func (g *Grid) MinRange(r1, c1, r2, c2 int) (float64, error) {
	values, err := g.valuesInRange(r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("min of range (%d,%d)-(%d,%d): %w", r1, c1, r2, c2, ErrEmptyGrid)
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m, nil
}

func (g *Grid) MaxRange(r1, c1, r2, c2 int) (float64, error) {
	values, err := g.valuesInRange(r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("max of range (%d,%d)-(%d,%d): %w", r1, c1, r2, c2, ErrEmptyGrid)
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m, nil
}

func (g *Grid) MinAll() (float64, error) {
	values := g.allValues()
	if len(values) == 0 {
		return 0, ErrEmptyGrid
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m, nil
}

func (g *Grid) MaxAll() (float64, error) {
	values := g.allValues()
	if len(values) == 0 {
		return 0, ErrEmptyGrid
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m, nil
}

type Clipboard struct {
	Rows, Cols int
	Values     map[Position]Cell
}

func (g *Grid) Copy(r1, c1, r2, c2 int) (Clipboard, error) {
	r1, c1, r2, c2, err := normalizeRange(r1, c1, r2, c2)
	if err != nil {
		return Clipboard{}, err
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	block := Clipboard{
		Rows:   r2 - r1 + 1,
		Cols:   c2 - c1 + 1,
		Values: make(map[Position]Cell),
	}
	for p, c := range g.cells {
		if p.Row >= r1 && p.Row <= r2 && p.Col >= c1 && p.Col <= c2 {
			block.Values[Position{p.Row - r1, p.Col - c1}] = c
		}
	}
	return block, nil
}

func (g *Grid) Paste(destRow, destCol int, block Clipboard) error {
	if err := validatePosition(destRow, destCol); err != nil {
		return err
	}
	if block.Values == nil {
		return errors.New("clipboard is empty: use 'copy' before 'paste'")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for offset, c := range block.Values {
		g.cells[Position{destRow + offset.Row, destCol + offset.Col}] = c
	}
	return nil
}

func (g *Grid) Find(value float64) []Position {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []Position
	for p, c := range g.cells {
		if c.Value == value {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Row != out[j].Row {
			return out[i].Row < out[j].Row
		}
		return out[i].Col < out[j].Col
	})
	return out
}

func (g *Grid) Clear(r1, c1, r2, c2 int) error {
	r1, c1, r2, c2, err := normalizeRange(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for row := r1; row <= r2; row++ {
		for col := c1; col <= c2; col++ {
			delete(g.cells, Position{row, col})
		}
	}
	return nil
}

func (g *Grid) Fill(r1, c1, r2, c2 int, value float64) error {
	r1, c1, r2, c2, err := normalizeRange(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for row := r1; row <= r2; row++ {
		for col := c1; col <= c2; col++ {
			g.cells[Position{row, col}] = Cell{Value: value}
		}
	}
	return nil
}

func (g *Grid) Sort(r1, c1, r2, c2 int, ascending bool) error {
	r1, c1, r2, c2, err := normalizeRange(r1, c1, r2, c2)
	if err != nil {
		return err
	}
	isSingleRow := r1 == r2
	isSingleCol := c1 == c2
	if !isSingleRow && !isSingleCol {
		return fmt.Errorf("sort only supports a single row or single column range, got (%d,%d)-(%d,%d)", r1, c1, r2, c2)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	var positions []Position
	if isSingleRow {
		for col := c1; col <= c2; col++ {
			p := Position{r1, col}
			if _, ok := g.cells[p]; ok {
				positions = append(positions, p)
			}
		}
	} else {
		for row := r1; row <= r2; row++ {
			p := Position{row, c1}
			if _, ok := g.cells[p]; ok {
				positions = append(positions, p)
			}
		}
	}
	if len(positions) < 2 {
		return nil
	}

	values := make([]float64, len(positions))
	for i, p := range positions {
		values[i] = g.cells[p].Value
	}
	sort.Float64s(values)
	if !ascending {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}
	for i, p := range positions {
		g.cells[p] = Cell{Value: values[i]}
	}
	return nil
}
