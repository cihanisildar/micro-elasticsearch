package analyzer

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Temel cumle",
			input:    "Kırmızı arabalar çok hızlı gidiyorlar!",
			expected: []string{"kırmızı", "arabalar", "hızlı", "gidiyorlar"},
		},
		{
			name:     "Noktalama isaretleri",
			input:    "Merhaba, nasılsın? Umarım her şey yolundadır.",
			expected: []string{"merhaba", "nasılsın", "umarım", "her", "şey", "yolundadır"},
		},
		{
			name:     "Sadece stop words",
			input:    "Bir ve ile için",
			expected: nil,
		},
		{
			name:     "Ingilizce metin",
			input:    "The quick brown fox jumps over the lazy dog",
			expected: []string{"quick", "brown", "fox", "jumps", "over", "lazy", "dog"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Analyze(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Analyze() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
