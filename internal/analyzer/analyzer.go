package analyzer

import (
	"strings"
	"unicode"
)

// stopWords, arama motorunun yoksayacağı (anlam ifade etmeyen) bağlaç ve kelimelerdir.
// Şimdilik Türkçe ve İngilizce yaygın kelimeleri karma olarak koyalım.
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "the": true, "is": true, "are": true, "was": true,
	"in": true, "on": true, "at": true, "for": true, "to": true, "of": true, "with": true,
	"bir": true, "ve": true, "ile": true, "için": true, "bu": true, "şu": true, "o": true,
	"çok": true, "da": true, "de": true, "mı": true, "mi": true,
}

// Analyze, sisteme giren metni temizler ve aranabilir kelimeler (tokenlar) dizisi olarak döndürür.
func Analyze(text string) []string {
	// 1. Tüm metni küçük harfe çevir
	text = strings.ToLower(text)

	// 2. Noktalama işaretlerini kaldır ve boşluklara göre böl
	// FieldsFunc, harf veya rakam olmayan her karakterden kelimeyi böler.
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// 3. Stop-word'leri (gereksiz kelimeleri) filtrele
	var filtered []string
	for _, token := range tokens {
		if !stopWords[token] {
			filtered = append(filtered, token)
		}
	}

	return filtered
}
