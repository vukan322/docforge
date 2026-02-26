package docforge

import "github.com/vukan322/docforge/internal/ooxml"

func Open(path string) (*ooxml.Document, error) {
	return ooxml.Open(path)
}
