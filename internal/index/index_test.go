package index

import (
	"reflect"
	"testing"
)

func TestIndex(t *testing.T) {
	idx := NewIndex()

	doc1 := Document{ID: "1", Text: "Kırmızı hızlı araba"}
	doc2 := Document{ID: "2", Text: "Mavi araba yavaş"}

	idx.Add(doc1)
	idx.Add(doc2)

	if idx.TotalDocs() != 2 {
		t.Errorf("Beklenen doküman sayısı 2, alınan %d", idx.TotalDocs())
	}

	// "araba" kelimesi her iki dokümanda da var
	arabaIDs := idx.GetDocIDs("araba")
	expectedAraba := []string{"1", "2"}
	if !reflect.DeepEqual(arabaIDs, expectedAraba) {
		t.Errorf("'araba' için beklenen %v, alınan %v", expectedAraba, arabaIDs)
	}

	// "kırmızı" kelimesi sadece 1 numaralı dokümanda var
	kirmiziIDs := idx.GetDocIDs("kırmızı")
	expectedKirmizi := []string{"1"}
	if !reflect.DeepEqual(kirmiziIDs, expectedKirmizi) {
		t.Errorf("'kırmızı' için beklenen %v, alınan %v", expectedKirmizi, kirmiziIDs)
	}

	// "uçak" kelimesi hiçbir dokümanda yok
	ucakIDs := idx.GetDocIDs("uçak")
	if len(ucakIDs) != 0 {
		t.Errorf("'uçak' kelimesi bulunmamalıydı, alınan %v", ucakIDs)
	}
}
