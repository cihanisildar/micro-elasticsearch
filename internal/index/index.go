package index

import (
	"sync"
	"micro-es/internal/analyzer"
)

// Index, arama motorumuzun kalbidir. Tüm veriler burada RAM'de tutulur.
type Index struct {
	// mu, eşzamanlı veri yazma/okuma işlemlerinde çakışmayı (race condition) engeller.
	mu sync.RWMutex

	// docs, orijinal metinleri ID'leri ile saklar (örn: "doc1" -> {ID: "doc1", Text: "..."})
	docs map[string]Document

	// invertedIndex (Ters Dizin), kelimelerden doküman ID'lerine giden haritadır.
	// Örn: "araba" -> ["doc1", "doc2"]
	invertedIndex map[string][]string
}

// NewIndex, boş bir indeks oluşturur ve döndürür.
func NewIndex() *Index {
	return &Index{
		docs:          make(map[string]Document),
		invertedIndex: make(map[string][]string),
	}
}

// Add, sisteme yeni bir doküman ekler.
func (idx *Index) Add(doc Document) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// 1. Orijinal dokümanı kaydet
	idx.docs[doc.ID] = doc

	// 2. Metni kelimelere (tokenlara) ayır
	tokens := analyzer.Analyze(doc.Text)

	// 3. Bir kelime dokümanda 5 kere geçse bile, inverted index'e "bu dokümanda var" 
	// bilgisini sadece 1 kez eklemeliyiz. Bunun için tokenları tekilleştiriyoruz.
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token] = true
	}

	// 4. Ters Dizini (Inverted Index) güncelle
	for token := range uniqueTokens {
		idx.invertedIndex[token] = append(idx.invertedIndex[token], doc.ID)
	}
}

// TotalDocs, sistemdeki toplam doküman sayısını döndürür.
func (idx *Index) TotalDocs() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.docs)
}

// GetDocIDs, verilen bir kelimenin (token) geçtiği doküman ID'lerini döndürür.
func (idx *Index) GetDocIDs(token string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.invertedIndex[token]
}

// GetDocument, verilen bir ID'ye ait dokümanı döndürür.
func (idx *Index) GetDocument(id string) (Document, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	doc, exists := idx.docs[id]
	return doc, exists
}
