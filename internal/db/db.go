package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"micro-es/internal/index"
)

// InitDB, veritabanı bağlantısını sağlar ve tabloları oluşturur.
func InitDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// Docker-compose ayağa kalkarken DB'nin hazır olması birkaç saniye sürebilir, bekliyoruz.
	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		return nil, pingErr
	}

	// Tabloyu oluştur (eğer yoksa)
	query := `CREATE TABLE IF NOT EXISTS documents (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL
	);`
	
	if _, err = db.Exec(query); err != nil {
		return nil, err
	}

	return db, nil
}

// SaveDocument, veriyi kalıcı olarak PostgreSQL'e yazar.
func SaveDocument(db *sql.DB, doc index.Document) error {
	// On Conflict (aynı ID ile gelirse) UPDATE et, yani üzerine yaz.
	query := "INSERT INTO documents (id, text) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET text = $2"
	_, err := db.Exec(query, doc.ID, doc.Text)
	return err
}

// LoadAllDocuments, sistem ilk açıldığında tüm verileri çekip RAM'e yüklemek için kullanılır.
func LoadAllDocuments(db *sql.DB) ([]index.Document, error) {
	rows, err := db.Query("SELECT id, text FROM documents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []index.Document
	for rows.Next() {
		var doc index.Document
		if err := rows.Scan(&doc.ID, &doc.Text); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
