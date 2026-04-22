package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"micro-es/internal/api"
	"micro-es/internal/db"
)

func main() {
	// 1. Veritabanına Bağlan ve Verileri RAM'e Yükle (Index Rebuilding)
	dbUrl := os.Getenv("DB_URL")
	if dbUrl != "" {
		database, err := db.InitDB(dbUrl)
		if err != nil {
			log.Fatalf("❌ DB bağlantı hatası: %v", err)
		}
		api.DB = database
		
		docs, err := db.LoadAllDocuments(database)
		if err != nil {
			log.Fatalf("❌ DB okuma hatası: %v", err)
		}
		
		for _, doc := range docs {
			api.GlobalIndex.Add(doc)
		}
		fmt.Printf("📦 Veritabanından %d doküman RAM'e (Ters Dizin'e) yüklendi.\n", len(docs))
	} else {
		fmt.Println("⚠️ DB_URL tanımlanmamış, sadece RAM üzerinde (geçici) çalışıyor.")
	}

	// 2. API Endpoint'lerini Tanımla
	http.HandleFunc("/api/documents", api.AddDocumentHandler)
	http.HandleFunc("/api/search", api.SearchHandler)

	// Sunucuyu Başlat
	port := ":8080"
	fmt.Printf("🚀 Mikro-Elasticsearch sunucusu baslatiliyor...\n")
	fmt.Printf("📡 Dinlenen Port: %s\n", port)
	fmt.Printf("   - Yeni Doküman Ekle: POST http://localhost%s/api/documents\n", port)
	fmt.Printf("   - Arama Yap:         GET  http://localhost%s/api/search?q=kelime\n", port)
	fmt.Println("------------------------------------------------------")

	// Sunucuyu ayağa kaldırıp sürekli dinlemeye başla
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Sunucu baslatilamadi: %v", err)
	}
}
