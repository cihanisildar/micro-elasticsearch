package analyzer

import (
	"strings"
	"unicode"
)

var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "the": true, "is": true, "are": true, "was": true,
	"in": true, "on": true, "at": true, "for": true, "to": true, "of": true, "with": true,
	"bir": true, "ve": true, "ile": true, "için": true, "bu": true, "şu": true, "o": true,
	"çok": true, "da": true, "de": true, "mı": true, "mi": true,
}

// Analyze, sisteme giren metni temizler ve aranabilir kelimeler (tokenlar) dizisi olarak döndürür.
func Analyze(text string) []string {
	text = strings.ToLower(text)

	// FieldsFunc, harf veya rakam olmayan her karakterden kelimeyi böler.
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var filtered []string
	for _, token := range tokens {
		if !stopWords[token] {
			filtered = append(filtered, token)
		}
	}

	return filtered
}
