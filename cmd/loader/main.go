package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Document, API'mize göndereceğimiz JSON yapısıdır.
type Document struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// sampleData, sisteme test amacıyla yükleyeceğimiz 10 adet örnek makale/cümledir.
var sampleData = []Document{
	{ID: "doc1", Text: "Go dili ile yüksek performanslı arama motoru geliştirmek oldukça keyiflidir."},
	{ID: "doc2", Text: "Elasticsearch, Lucene tabanlı dağıtık bir arama ve analiz motorudur."},
	{ID: "doc3", Text: "Docker kullanarak projelerimizi izole bir şekilde saniyeler içinde ayağa kaldırabiliriz."},
	{ID: "doc4", Text: "TF-IDF algoritması kelimelerin doküman içindeki ve tüm veri setindeki önemini hesaplar."},
	{ID: "doc5", Text: "Go dili concurrency (eşzamanlılık) konusunda goroutine'ler sayesinde çok başarılıdır."},
	{ID: "doc6", Text: "Arama algoritmalarında ters dizin (inverted index) kullanmak performansı O(1) seviyesine yaklaştırır."},
	{ID: "doc7", Text: "Mikroservis mimarisinde Docker ve Kubernetes ikilisi modern yazılımın vazgeçilmezidir."},
	{ID: "doc8", Text: "Veri yapıları ve algoritmalar mülakatlarda en çok sorulan konuların başında gelir."},
	{ID: "doc9", Text: "Go dilinin standart kütüphanesi HTTP sunucusu kurmak için dış bağımlılık olmadan yeterince güçlüdür."},
	{ID: "doc10", Text: "Kırmızı arabalar çok hızlı gider ama mavi arabalar daha konforludur."},
}

func main() {
	url := "http://localhost:8080/api/documents"
	fmt.Println("🚀 Örnek veriler (Mini-Logstash) sunucuya yükleniyor...")

	for _, doc := range sampleData {
		jsonData, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("JSON dönüştürme hatası: %v", err)
		}

		// Sunucuya POST isteği gönderiyoruz
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("🚨 Sunucuya bağlanılamadı. Sunucunun 8080 portunda açık olduğundan emin olun: %v", err)
		}
		
		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("✅ Doküman eklendi: %s\n", doc.ID)
		} else {
			fmt.Printf("❌ Hata oluştu (HTTP Kodu: %d): %s\n", resp.StatusCode, doc.ID)
		}
		resp.Body.Close()
		
		// İstekleri saniyede yüzlerce atmamak için çok ufak bir mola veriyoruz
		time.Sleep(20 * time.Millisecond)
	}
	
	fmt.Println("🎉 Tüm veriler başarıyla yüklendi!")
	fmt.Println("👉 Tarayıcınızda test edebilirsiniz:")
	fmt.Println("   http://localhost:8080/api/search?q=go+dili")
	fmt.Println("   http://localhost:8080/api/search?q=docker")
}
