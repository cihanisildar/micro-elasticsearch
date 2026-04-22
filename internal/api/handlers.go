package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"micro-es/internal/db"
	"micro-es/internal/index"
	"micro-es/internal/query"
)

// GlobalIndex, sunucu çalıştığı sürece verilerimizi RAM'de tutacak olan ortak nesnemizdir.
var GlobalIndex = index.NewIndex()

// DB, veritabanı bağlantımızı tutan global değişkendir.
var DB *sql.DB

// AddDocumentRequest, kullanıcının bize göndereceği JSON verisinin yapısıdır.
type AddDocumentRequest struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// AddDocumentHandler, POST /api/documents isteğini karşılar ve sisteme yeni veri ekler.
func AddDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// Sadece POST destekliyoruz
	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST istegi kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	// Gelen JSON verisini oku
	var req AddDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Gecersiz JSON formati", http.StatusBadRequest)
		return
	}

	// Boş veri kontrolü
	if req.ID == "" || req.Text == "" {
		http.Error(w, "ID ve Text alanlari zorunludur", http.StatusBadRequest)
		return
	}

	// Index'e ekle
	doc := index.Document{
		ID:   req.ID,
		Text: req.Text,
	}

	// 1. Veriyi Kalıcı Olarak Veritabanına Yaz (Source of Truth)
	if DB != nil {
		if err := db.SaveDocument(DB, doc); err != nil {
			http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
			return
		}
	}

	// 2. Veriyi Hızlı Arama İçin RAM'e (Ters Dizin) Ekle (Search Engine)
	GlobalIndex.Add(doc)

	// Başarı cevabı dön
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"mesaj": "Doküman başarıyla eklendi",
		"id":    req.ID,
	})
}

// SearchHandler, GET /api/search?q=kelime isteğini karşılar ve arama yapar.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Sadece GET destekliyoruz
	if r.Method != http.MethodGet {
		http.Error(w, "Sadece GET istegi kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	// URL'den 'q' parametresini al (örnek: ?q=kirmizi+araba)
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "Arama terimi (q parametresi) bos olamaz", http.StatusBadRequest)
		return
	}

	// TF-IDF Arama fonksiyonumuzu çağır
	results := query.Search(GlobalIndex, q)

	// Eğer sonuç yoksa null dönmemesi için boş dizi oluştur
	if results == nil {
		results = []query.SearchResult{}
	}

	// Sonuçları JSON olarak dön
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
