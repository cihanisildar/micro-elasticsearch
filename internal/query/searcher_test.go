package query

import (
	"testing"

	"micro-es/internal/index"
)

func TestSearch(t *testing.T) {
	idx := index.NewIndex()

	idx.Add(index.Document{ID: "1", Text: "Kırmızı araba hızlı"})
	idx.Add(index.Document{ID: "2", Text: "Mavi araba yavaş"})
	idx.Add(index.Document{ID: "3", Text: "Kırmızı bisiklet"})

	// "kırmızı araba" araması
	results := Search(idx, "kırmızı araba")

	if len(results) != 3 {
		t.Fatalf("Beklenen sonuc sayisi 3, alinan %d", len(results))
	}

	// 1. Sonuç: Doc 1 olmalı (kırmızı ve araba var, tam eşleşme ve yüksek skor)
	if results[0].Document.ID != "1" {
		t.Errorf("Beklenen en iyi eslesme '1', alinan '%s'", results[0].Document.ID)
	}

	// 2. Sonuç: Doc 3 olmalı (sadece 'kırmızı' var ama 'araba'dan daha nadir olduğu için skoru yüksek)
	if results[1].Document.ID != "3" {
		t.Errorf("Beklenen ikinci eslesme '3', alinan '%s'", results[1].Document.ID)
	}

	// 3. Sonuç: Doc 2 olmalı (sadece 'araba' var, 'araba' kelimesi daha yaygın olduğundan skoru düşük)
	if results[2].Document.ID != "2" {
		t.Errorf("Beklenen ucuncu eslesme '2', alinan '%s'", results[2].Document.ID)
	}
}
