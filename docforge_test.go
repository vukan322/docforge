package docforge_test

import (
	"os"
	"testing"

	"github.com/vukan322/docforge"
)

const testTemplate = "testdata/templates/test-template.docx"
const testOutput = "testdata/templates/test-output.docx"

func TestOpen_ValidFile(t *testing.T) {
	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if doc == nil {
		t.Fatal("expected document, got nil")
	}
}

func TestOpen_InvalidFile(t *testing.T) {
	_, err := docforge.Open("nonexistent.docx")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestReplace_Placeholder(t *testing.T) {
	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}

	err = doc.Replace(map[string]any{
		"document_name": "Test Document",
	})
	if err != nil {
		t.Fatalf("failed to replace: %v", err)
	}

	err = doc.Save(testOutput)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	defer func() { _ = os.Remove(testOutput) }()
}

func TestReplace_StructTag(t *testing.T) {
	type Invoice struct {
		DocumentName string `docforge:"document_name"`
	}

	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}

	err = doc.Replace(Invoice{DocumentName: "Tagged Invoice"})
	if err != nil {
		t.Fatalf("failed to replace: %v", err)
	}

	err = doc.Save(testOutput)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	defer func() { _ = os.Remove(testOutput) }()
}

func TestRender_WithItems(t *testing.T) {
	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}

	err = doc.Render(map[string]any{
		"document_name": "Test Invoice",
		"Items": []map[string]any{
			{"Product": "Valve", "Qty": 3, "Price": 150.00},
			{"Product": "Pump", "Qty": 1, "Price": 890.00},
		},
	}, testOutput)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	defer func() { _ = os.Remove(testOutput) }()
}

func TestRender_WithAddFunc(t *testing.T) {
	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}

	doc.AddFunc("upper", func(s string) string {
		result := make([]byte, len(s))
		for i := range s {
			if s[i] >= 'a' && s[i] <= 'z' {
				result[i] = s[i] - 32
			} else {
				result[i] = s[i]
			}
		}
		return string(result)
	})

	err = doc.Render(map[string]any{
		"document_name": "Test Invoice",
		"Items": []map[string]any{
			{"Product": "Valve", "Qty": 3, "Price": 150.00},
		},
	}, testOutput)
	if err != nil {
		t.Fatalf("failed to render with func: %v", err)
	}

	defer func() { _ = os.Remove(testOutput) }()
}

func TestSave_RoundTrip(t *testing.T) {
	doc, err := docforge.Open(testTemplate)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}

	err = doc.Save(testOutput)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	_, err = docforge.Open(testOutput)
	if err != nil {
		t.Fatalf("saved file is not a valid docx: %v", err)
	}

	defer func() { _ = os.Remove(testOutput) }()
}
