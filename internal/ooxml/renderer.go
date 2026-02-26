package ooxml

import (
	"bytes"
	"fmt"
	"text/template"
)

func (d *Document) Render(data any, outputPath string) error {
	fields, err := extractFields(data)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{}
	for name, fn := range d.funcs {
		funcMap[name] = fn
	}

	tpl, err := template.New("docforge").Funcs(funcMap).Parse(string(d.xmlDoc))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, fields); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	d.xmlDoc = buf.Bytes()
	return d.Save(outputPath)
}

func (d *Document) AddFunc(name string, fn any) {
	if d.funcs == nil {
		d.funcs = make(map[string]any)
	}
	d.funcs[name] = fn
}
