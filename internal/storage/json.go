package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type JSONStorage struct{}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{}
}

func (JSONStorage) Save(path string, doc Document) error {
	if path == "" {
		return fmt.Errorf("save: path is empty")
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("save: creating directory %q: %w", dir, err)
		}
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("save: encoding JSON: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("save: writing file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("save: finalizing file: %w", err)
	}
	return nil
}

func (JSONStorage) Load(path string) (Document, error) {
	if path == "" {
		return Document{}, fmt.Errorf("load: path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, fmt.Errorf("load: reading file: %w", err)
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return Document{}, fmt.Errorf("load: decoding JSON: %w", err)
	}
	return doc, nil
}
