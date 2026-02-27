package ooxml

import (
	"bytes"
	"regexp"
)

var (
	tableRowRegex  = regexp.MustCompile(`(?s)<w:tr\b[^>]*>.*?</w:tr>`)
	cellTextRegex  = regexp.MustCompile(`(?s)<w:t[^>]*>(.*?)</w:t>`)
	rangeMarkerReg = regexp.MustCompile(`{{range[^}]*}}`)
)

func preprocessTables(xmlDoc []byte) []byte {
	tableRegex := regexp.MustCompile(`(?s)<w:tbl\b[^>]*>.*?</w:tbl>`)
	return tableRegex.ReplaceAllFunc(xmlDoc, preprocessTable)
}

func preprocessTable(table []byte) []byte {
	rows := tableRowRegex.FindAllIndex(table, -1)
	if len(rows) == 0 {
		return table
	}

	rangeRowIdx := -1
	endRowIdx := -1
	dataRowIdx := -1
	depth := 0

	openRe := regexp.MustCompile(`{{(range|if|with)\b`)

	for i, loc := range rows {
		row := table[loc[0]:loc[1]]
		text := extractAllText(row)

		opens := len(openRe.FindAll(text, -1))
		ends := bytes.Count(text, []byte("{{end}}"))
		net := opens - ends

		if rangeRowIdx == -1 && bytes.Contains(text, []byte("{{range")) {
			rangeRowIdx = i
			depth = 1
			continue
		}

		if rangeRowIdx != -1 {
			depth += net
			if depth <= 0 {
				endRowIdx = i
				break
			}
			if dataRowIdx == -1 {
				dataRowIdx = i
			}
		}
	}

	if rangeRowIdx == -1 || endRowIdx == -1 || dataRowIdx == -1 {
		return table
	}

	rangeText := rangeMarkerReg.Find(table[rows[rangeRowIdx][0]:rows[rangeRowIdx][1]])
	if rangeText == nil {
		return table
	}

	dataRow := make([]byte, rows[dataRowIdx][1]-rows[dataRowIdx][0])
	copy(dataRow, table[rows[dataRowIdx][0]:rows[dataRowIdx][1]])

	var result []byte
	result = append(result, table[:rows[rangeRowIdx][0]]...)
	result = append(result, rangeText...)
	result = append(result, dataRow...)
	result = append(result, []byte("{{end}}")...)
	result = append(result, table[rows[endRowIdx][1]:]...)
	return result
}

func extractAllText(data []byte) []byte {
	matches := cellTextRegex.FindAllSubmatch(data, -1)
	var result []byte
	for _, m := range matches {
		result = append(result, m[1]...)
	}
	return result
}

func PreprocessTablesDebug(xmlDoc []byte) []byte {
	return preprocessTables(xmlDoc)
}
