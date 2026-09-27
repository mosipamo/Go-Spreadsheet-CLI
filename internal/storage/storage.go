package storage

import "myexcel/internal/spreadsheet"

type CellRecord struct {
	Row   int     `json:"row"`
	Col   int     `json:"col"`
	Value float64 `json:"value"`
}

type Document struct {
	Cells []CellRecord `json:"cells"`
}

type Storage interface {
	Save(path string, doc Document) error
	Load(path string) (Document, error)
}

func ToDocument(snapshot map[spreadsheet.Position]spreadsheet.Cell) Document {
	doc := Document{Cells: make([]CellRecord, 0, len(snapshot))}
	for pos, cell := range snapshot {
		doc.Cells = append(doc.Cells, CellRecord{Row: pos.Row, Col: pos.Col, Value: cell.Value})
	}
	return doc
}

func ToGridCells(doc Document) map[spreadsheet.Position]spreadsheet.Cell {
	cells := make(map[spreadsheet.Position]spreadsheet.Cell, len(doc.Cells))
	for _, rec := range doc.Cells {
		cells[spreadsheet.Position{Row: rec.Row, Col: rec.Col}] = spreadsheet.Cell{Value: rec.Value}
	}
	return cells
}
