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
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	d.files[documentXMLPath] = d.xmlDoc

	for name, data := range d.files {
		f, err := w.Create(name)
		if err != nil {
			return fmt.Errorf("failed to create zip entry %s: %w", name, err)
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write zip entry %s: %w", name, err)
		}
	}

	return nil
}
