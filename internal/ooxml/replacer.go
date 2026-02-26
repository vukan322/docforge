package ooxml

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func (d *Document) Replace(data any) error {
	fields, err := extractFields(data)
	if err != nil {
		return err
	}

	for placeholder, value := range fields {
		old := []byte("{{." + placeholder + "}}")
		new := []byte(sanitizeXML(fmt.Sprintf("%v", value)))
		d.xmlDoc = bytes.ReplaceAll(d.xmlDoc, old, new)
	}

	return nil
}

func extractFields(data any) (map[string]any, error) {
	result := make(map[string]any)

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			key := field.Name
			if tag := field.Tag.Get("docforge"); tag != "" {
				key = tag
			}
			result[key] = value.Interface()
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			result[key.String()] = v.MapIndex(key).Interface()
		}
	default:
		return nil, fmt.Errorf("data must be a struct or map, got %s", v.Kind())
	}

	return result, nil
}

func sanitizeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
