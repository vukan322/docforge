package ooxml

import (
	"bytes"
	"testing"
)

func TestPreprocessTable_NoRangeMarker(t *testing.T) {
	input := []byte(`<w:tbl><w:tblPr></w:tblPr><w:tblGrid></w:tblGrid><w:tr><w:tc><w:p><w:r><w:t>Hello</w:t></w:r></w:p></w:tc></w:tr></w:tbl>`)
	result := preprocessTables(input)
	if string(result) != string(input) {
		t.Errorf("expected unchanged input when no range marker, got %s", result)
	}
}

func TestPreprocessTable_MovesRangeOutsideRow(t *testing.T) {
	input := []byte(`<w:tbl><w:tblPr></w:tblPr><w:tblGrid></w:tblGrid>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{range .Items}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{.Name}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{end}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`</w:tbl>`)

	result := preprocessTables(input)

	rangeIdx := bytes.Index(result, []byte("{{range .Items}}"))
	trIdx := bytes.Index(result, []byte("<w:tr>"))

	if rangeIdx == -1 {
		t.Fatal("expected {{range .Items}} in output")
	}
	if trIdx == -1 {
		t.Fatal("expected <w:tr> in output")
	}
	if rangeIdx > trIdx {
		t.Errorf("expected {{range .Items}} to appear before <w:tr>, but it appears after")
	}
}

func TestPreprocessTable_EndOutsideRow(t *testing.T) {
	input := []byte(`<w:tbl><w:tblPr></w:tblPr><w:tblGrid></w:tblGrid>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{range .Items}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{.Name}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{end}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`</w:tbl>`)

	result := preprocessTables(input)

	endIdx := bytes.Index(result, []byte("{{end}}"))
	lastTrIdx := bytes.LastIndex(result, []byte("</w:tr>"))

	if endIdx == -1 {
		t.Fatal("expected {{end}} in output")
	}
	if endIdx < lastTrIdx {
		t.Errorf("expected {{end}} to appear after last </w:tr>")
	}
}

func TestPreprocessTable_DataRowPreserved(t *testing.T) {
	input := []byte(`<w:tbl><w:tblPr></w:tblPr><w:tblGrid></w:tblGrid>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{range .Items}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{.Name}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr><w:tc><w:p><w:r><w:t>{{end}}</w:t></w:r></w:p></w:tc></w:tr>` +
		`</w:tbl>`)

	result := preprocessTables(input)

	if !bytes.Contains(result, []byte("{{.Name}}")) {
		t.Errorf("expected data row placeholder {{.Name}} to be preserved")
	}
}
