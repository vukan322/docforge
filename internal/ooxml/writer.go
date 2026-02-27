package ooxml

import (
	"archive/zip"
	"fmt"
	"os"
)

func (d *Document) Save(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = out.Close() }()

	w := zip.NewWriter(out)

	d.files[documentXMLPath] = d.xmlDoc
	for name, data := range d.files {
		f, err := w.Create(name)
		if err != nil {
			_ = w.Close()
			return fmt.Errorf("failed to create zip entry %s: %w", name, err)
		}
		_, err = f.Write(data)
		if err != nil {
			_ = w.Close()
			return fmt.Errorf("failed to write zip entry %s: %w", name, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to finalize zip: %w", err)
	}
	return nil
}
