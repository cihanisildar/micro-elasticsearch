package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

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

	if DB != nil {
		if err := db.SaveDocument(DB, doc); err != nil {
			http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
			return
		}
	}

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

// sampleDocuments, UI üzerinden yüklenebilecek örnek veri setidir.
var sampleDocuments = []index.Document{
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

// SeedHandler, POST /api/seed isteğini karşılar ve örnek 10 dokümanı sisteme yükler.
func SeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST istegi kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	loaded := 0
	for _, doc := range sampleDocuments {
		if DB != nil {
			if err := db.SaveDocument(DB, doc); err != nil {
				continue
			}
		}
		GlobalIndex.Add(doc)
		loaded++
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<span class="seed-success">✅ ` + fmt.Sprintf("%d", loaded) + ` doküman yüklendi!</span>`))
}

// SearchHTMLHandler, HTMX istekleri için JSON yerine Go Templates kullanarak doğrudan HTML döndürür.
func SearchHTMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Sadece GET istegi kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	
	var results []query.SearchResult
	
	if q == "" {
		// Arama boşsa tüm dokümanları (varsayılan 0 skoruyla) göster
		allDocs := GlobalIndex.GetAllDocuments()
		for _, doc := range allDocs {
			results = append(results, query.SearchResult{
				Document: doc,
				Score:    0.0,
			})
		}
	} else {
		results = query.Search(GlobalIndex, q)
	}

	if len(results) == 0 {
		w.Write([]byte(`<div class="no-results">Aradığınız kelimeye uygun doküman bulunamadı. 🔍</div>`))
		return
	}

	// Kartların HTML Şablonu
	const tmplHTML = `
		{{range $index, $item := .}}
		<div class="result-card" style="animation-delay: {{mul $index 0.08}}s;">
			<div class="result-text">{{highlight $item.Document.Text}}</div>
			<div class="result-meta">
				<div class="score-badge">Score: {{printf "%.4f" $item.Score}}</div>
				<div class="doc-id">ID: {{$item.Document.ID}}</div>
			</div>
		</div>
		{{end}}
	`

	// Animasyon gecikmesi ve vurgulama için yardımcı fonksiyonlar
	funcMap := template.FuncMap{
		"mul": func(a int, b float64) float64 {
			return float64(a) * b
		},
		"highlight": func(text string) template.HTML {
			if q == "" {
				return template.HTML(text)
			}
			terms := strings.Fields(q)
			highlighted := text
			for _, term := range terms {
				re := regexp.MustCompile(`(?i)(` + regexp.QuoteMeta(term) + `)`)
				highlighted = re.ReplaceAllString(highlighted, `<span class="highlight">$1</span>`)
			}
			return template.HTML(highlighted)
		},
	}

	tmpl, err := template.New("results").Funcs(funcMap).Parse(tmplHTML)
	if err != nil {
		http.Error(w, "Şablon hatası", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, results)
}
