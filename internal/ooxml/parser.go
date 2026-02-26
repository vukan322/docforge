package ooxml

import (
	"archive/zip"
	"fmt"
	"io"
)

const documentXMLPath = "word/document.xml"

type Document struct {
	files  map[string][]byte
	xmlDoc []byte
}

func Open(path string) (*Document, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open docx: %w", err)
	}
	defer r.Close()

	doc := &Document{
		files: make(map[string][]byte),
	}

	for _, f := range r.File {
		data, err := readZipFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", f.Name, err)
		}
		doc.files[f.Name] = data
	}

	xmlDoc, ok := doc.files[documentXMLPath]
	if !ok {
		return nil, fmt.Errorf("invalid docx: missing %s", documentXMLPath)
	}
	doc.xmlDoc = xmlDoc

	return doc, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}
