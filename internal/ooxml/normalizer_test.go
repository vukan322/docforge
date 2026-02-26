package ooxml

import (
	"bytes"
	"testing"
)

func TestNormalizeRuns_NoSplit(t *testing.T) {
	input := []byte(`<w:p><w:r><w:rPr></w:rPr><w:t>{{.Name}}</w:t></w:r></w:p>`)
	result := normalizeRuns(input)
	if string(result) != string(input) {
		t.Errorf("expected unchanged input, got %s", result)
	}
}

func TestNormalizeRuns_SplitPlaceholder(t *testing.T) {
	input := []byte(`<w:p><w:r><w:rPr></w:rPr><w:t>{{.Na</w:t></w:r><w:r><w:rPr></w:rPr><w:t>me}}</w:t></w:r></w:p>`)
	result := normalizeRuns(input)
	if !bytes.Contains(result, []byte("{{.Name}}")) {
		t.Errorf("expected merged placeholder {{.Name}}, got %s", result)
	}
}

func TestNormalizeRuns_PreservesNonPlaceholderRuns(t *testing.T) {
	input := []byte(`<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Hello</w:t></w:r><w:r><w:rPr><w:i/></w:rPr><w:t>World</w:t></w:r></w:p>`)
	result := normalizeRuns(input)
	if !bytes.Contains(result, []byte("Hello")) || !bytes.Contains(result, []byte("World")) {
		t.Errorf("expected both runs preserved, got %s", result)
	}
}
