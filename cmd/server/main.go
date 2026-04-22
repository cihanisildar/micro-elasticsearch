package main

import (
	"fmt"
	"log"
	"net/http"

	"micro-es/internal/api"
)

func main() {
	// API Endpoint'lerini (Yönlendirmelerini) Tanımla
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
