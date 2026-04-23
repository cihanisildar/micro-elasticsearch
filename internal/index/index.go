package index

import (
	"sync"
	"micro-es/internal/analyzer"
)

// Index, arama motorumuzun kalbidir. Tüm veriler burada RAM'de tutulur.
type Index struct {
	mu sync.RWMutex

	docs map[string]Document

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

	idx.docs[doc.ID] = doc

	tokens := analyzer.Analyze(doc.Text)

	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token] = true
	}

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

// GetAllDocuments, sistemdeki tüm dokümanları liste halinde döndürür.
func (idx *Index) GetAllDocuments() []Document {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	
	var allDocs []Document
	for _, doc := range idx.docs {
		allDocs = append(allDocs, doc)
	}
	return allDocs
}
