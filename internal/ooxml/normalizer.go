package ooxml

import (
	"bytes"
	"regexp"
)

var (
	runRegex  = regexp.MustCompile(`(?s)<w:r\b[^>]*>.*?</w:r>`)
	textRegex = regexp.MustCompile(`(?s)<w:t\b[^>]*>(.*?)</w:t>`)
)

func normalizeRuns(xmlDoc []byte) []byte {
	paragraphRegex := regexp.MustCompile(`(?s)<w:p\b[^>]*>.*?</w:p>`)
	return paragraphRegex.ReplaceAllFunc(xmlDoc, normalizeParagraph)
}

func normalizeParagraph(para []byte) []byte {
	runs := runRegex.FindAllIndex(para, -1)
	if len(runs) == 0 {
		return para
	}

	var result []byte
	i := 0
	lastEnd := 0

	for i < len(runs) {
		start, end := runs[i][0], runs[i][1]
		runBytes := para[start:end]
		text := extractText(runBytes)

		if bytes.Contains(text, []byte("{{")) && !bytes.Contains(text, []byte("}}")) {
			merged := runBytes
			j := i + 1
			for j < len(runs) {
				nextRun := para[runs[j][0]:runs[j][1]]
				nextText := extractText(nextRun)
				merged = mergeRunText(merged, nextText)
				j++
				if bytes.Contains(extractText(merged), []byte("}}")) {
					break
				}
			}
			result = append(result, para[lastEnd:start]...)
			result = append(result, merged...)
			lastEnd = runs[j-1][1]
			i = j
		} else {
			result = append(result, para[lastEnd:end]...)
			lastEnd = end
			i++
		}
	}

	result = append(result, para[lastEnd:]...)
	return result
}

func extractText(r []byte) []byte {
	match := textRegex.FindSubmatch(r)
	if match == nil {
		return nil
	}
	return match[1]
}

func mergeRunText(base, additional []byte) []byte {
	return textRegex.ReplaceAllFunc(base, func(match []byte) []byte {
		current := textRegex.FindSubmatch(match)[1]
		merged := append(current, additional...)
		return bytes.Replace(match, current, merged, 1)
	})
}
