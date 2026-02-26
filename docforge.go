package docforge

import "github.com/vukan322/docforge/internal/ooxml"

type Document struct {
	doc *ooxml.Document
}

func Open(path string) (*Document, error) {
	d, err := ooxml.Open(path)
	if err != nil {
		return nil, err
	}
	return &Document{doc: d}, nil
}

func (d *Document) Replace(data any) error {
	return d.doc.Replace(data)
}

func (d *Document) Save(path string) error {
	return d.doc.Save(path)
}
