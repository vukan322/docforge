# docforge

A Go library for creating, modifying, and rendering DOCX documents.
Handles everything from generating new files to filling templates with data,
without mutating your document styles.

> docforge is under active development. See the changelog for current status.

---

## Installation

'go get github.com/vukan322/docforge'

---

## Simple replacement

For basic placeholder replacement without a full template engine:

```go
doc, err := docforge.Open("template.docx")
if err != nil {
    log.Fatal(err)
}

err = doc.Replace(map[string]any{
    "Name":    "Luka",
    "Company": "Kantech",
})
if err != nil {
    log.Fatal(err)
}

err = doc.Save("output.docx")
if err != nil {
    log.Fatal(err)
}
```

In your Word template use '{{.Name}}' and '{{.Company}}' as placeholders.

---

## Full template rendering

For range, conditionals, and custom functions:

```go
doc, err := docforge.Open("template.docx")
if err != nil {
    log.Fatal(err)
}

doc.AddFunc("formatDate", func(t time.Time) string {
    return t.Format("02.01.2006")
})

err = doc.Render(map[string]any{
    "Name": "Luka",
    "Date": time.Now(),
    "Items": []map[string]any{
        {"Product": "Valve", "Qty": 3},
        {"Product": "Pump",  "Qty": 1},
    },
}, "output.docx")
if err != nil {
    log.Fatal(err)
}
```

---

## Struct tag support

Use the 'docforge' struct tag to control placeholder names:

```go
type Invoice struct {
    FirstName string  `docforge:"first_name"`
    Amount    float64 `docforge:"amount"`
}
```

'{{.first_name}}' and '{{.amount}}' in your template will map to the struct fields.
Without a tag, the field name is used as-is.

---

## Template syntax

| Syntax | Purpose |
|---|---|
| '{{.FieldName}}' | Simple value replacement |
| '{{range .Items}} ... {{end}}' | Table row iteration |
| '{{if .Flag}} ... {{end}}' | Conditional content |
| '{{.Date \| formatDate}}' | Custom function |

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
if err != nil {
    log.Fatal(err)
}
```

---

## License

This project is licensed under the MIT License - see the 'LICENSE' file for details.
