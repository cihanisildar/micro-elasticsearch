package query

import (
	"math"
	"sort"

	"micro-es/internal/analyzer"
	"micro-es/internal/index"
)

// SearchResult, arama sonucu dönen dokümanı ve skorunu (ilgi puanı) tutar.
type SearchResult struct {
	Document index.Document
	Score    float64
}

// Search, verilen arama terimine göre dokümanları puanlar ve sıralar.
func Search(idx *index.Index, query string) []SearchResult {
	// 1. Arama metnini kelimelere ayır (örn: "Kırmızı Araba" -> ["kırmızı", "araba"])
	tokens := analyzer.Analyze(query)
	if len(tokens) == 0 {
		return nil
	}

	docScores := make(map[string]float64)
	totalDocs := float64(idx.TotalDocs())

	// 2. Her bir kelime için TF-IDF hesapla
	for _, token := range tokens {
		matchedDocIDs := idx.GetDocIDs(token)
		if len(matchedDocIDs) == 0 {
			continue // Bu kelime hiçbir dokümanda yoksa atla
		}

		// IDF (Inverse Document Frequency) Hesaplaması
		// Formül: log10( Toplam Doküman Sayısı / Kelimenin Bulunduğu Doküman Sayısı )
		// Bir kelime ne kadar az dokümanda geçiyorsa, değeri o kadar yüksektir.
		idf := math.Log10(totalDocs / float64(len(matchedDocIDs)))

		// Bu kelimenin geçtiği her bir doküman için skor hesapla
		for _, docID := range matchedDocIDs {
			doc, exists := idx.GetDocument(docID)
			if !exists {
				continue
			}

			// Dokümanı kelimelere ayır
			docTokens := analyzer.Analyze(doc.Text)
			termCount := 0
			for _, dt := range docTokens {
				if dt == token {
					termCount++
				}
			}

			// TF (Term Frequency) Hesaplaması
			// Formül: Kelimenin o dokümanda geçme sayısı / Dokümandaki toplam kelime sayısı
			tf := float64(termCount) / float64(len(docTokens))

			// Dokümanın toplam puanına bu kelimenin TF-IDF değerini ekle
			docScores[docID] += tf * idf
		}
	}

	// 3. Puanlanmış dokümanları listeye çevir
	var results []SearchResult
	for docID, score := range docScores {
		if doc, exists := idx.GetDocument(docID); exists {
			results = append(results, SearchResult{
				Document: doc,
				Score:    score,
			})
		}
	}

	// 4. Skorlara göre büyükten küçüğe sırala (Ranking)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}
