# docforge

A Go library for creating, modifying, and rendering DOCX documents.
Handles everything from generating new files to filling templates with data,
without mutating your document styles.

> docforge is under active development. See the changelog for current status.

---

## Installation

'go get github.com/vukan322/docforge'

---

## Create a document from scratch

```go
doc := docforge.New()

doc.AddParagraph("Hello from docforge")
doc.AddTable([][]string{
    {"Name", "Quantity"},
    {"Valve", "3"},
    {"Pump", "1"},
})

err := doc.Save("output.docx")
```

---

## Fill a template with data

```go
tpl, err := docforge.OpenTemplate("template.docx")
if err != nil {
    log.Fatal(err)
}

data := map[string]any{
    "Name":    "Luka",
    "Company": "Kantech",
}

err = tpl.Render(data, "output.docx")
```

In your Word template use '{{.Name}}' and '{{.Company}}' as placeholders.

---

## Modify an existing document

```go
doc, err := docforge.Open("existing.docx")
if err != nil {
    log.Fatal(err)
}

doc.ReplaceParagraph("old text", "new text")

err = doc.Save("modified.docx")
```

---

## Template Syntax

| Syntax | Purpose |
|---|---|
| '{{.FieldName}}' | Simple value replacement |
| '{{range .Slice}} ... {{end}}' | Table row iteration |
| '{{if .Flag}} ... {{end}}' | Conditional content |

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.
