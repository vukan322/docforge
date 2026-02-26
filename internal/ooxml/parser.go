package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
)

const documentXMLPath = "word/document.xml"

type Document struct {
	files  map[string][]byte
	xmlDoc []byte
	funcs  map[string]any
}

func (d *Document) GetXML() []byte {
	return d.xmlDoc
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
	doc.xmlDoc = normalizeRuns(xmlDoc)

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

type XMLBody struct {
	XMLName  xml.Name `xml:"body"`
	InnerXML []byte   `xml:",innerxml"`
}

type XMLDocument struct {
	XMLName  xml.Name `xml:"document"`
	InnerXML []byte   `xml:",innerxml"`
}
