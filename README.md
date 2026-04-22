# 🚀 Micro-Elasticsearch (Go)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

A high-performance, in-memory full-text search engine written from scratch in Go. This project simulates the core mechanics of large-scale search engines like Elasticsearch, utilizing an **Inverted Index** for rapid lookups and the **TF-IDF algorithm** for intelligent document ranking.

## ✨ Features

- **Blazing Fast Search:** Custom in-memory Inverted Index for instantaneous full-text search queries `O(1)`.
- **Intelligent Ranking (TF-IDF):** Sorts search results based on Term Frequency & Inverse Document Frequency.
- **Text Analysis & Tokenization:** Built-in analyzer that converts text to lowercase, strips punctuation, and filters out common stop-words.
- **Persistent Storage (Source of Truth):** Uses PostgreSQL to durably store documents.
- **Automatic Index Rebuilding:** On startup, the engine queries PostgreSQL and seamlessly rebuilds the in-memory index to recover state.
- **Production-Ready Docker Setup:** Multi-stage Dockerfile resulting in a tiny `~15MB` application image, orchestrated with `docker-compose`.

## 🏗️ Architecture

1. **Write Path:** When a new document is sent via the API, it is first written to **PostgreSQL** (Source of Truth) for persistence, and then immediately indexed into the RAM (In-memory engine).
2. **Read Path:** Search queries hit the **Inverted Index** directly in memory, calculating the TF-IDF score on-the-fly and returning ranked results in milliseconds.
3. **Recovery:** If the container or server restarts, the engine automatically fetches all data from PostgreSQL and reconstructs the search index before accepting requests.

## 🚀 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) and Docker Compose

### Running the Engine

Spin up the entire stack (Search Engine + PostgreSQL) with a single command:

```bash
docker-compose up -d --build
```

The engine will be available at `http://localhost:8080`.

### Populating Dummy Data (Mini-Logstash)
To test the engine, you can run the provided data loader script which pushes sample documents into the engine:
```bash
go run cmd/loader/main.go
```

## 🔌 API Usage

### 1. Add a Document
```bash
curl -X POST http://localhost:8080/api/documents \
-H "Content-Type: application/json" \
-d '{"id":"test1", "text":"Go is amazing for building search engines"}'
```

### 2. Search Documents
```bash
curl "http://localhost:8080/api/search?q=search+engines"
```
**Example Response:**
```json
[
  {
    "Document": {
      "ID": "test1",
      "Text": "Go is amazing for building search engines"
    },
    "Score": 0.3521
  }
]
```

## 📂 Project Structure
- `cmd/server/`: The HTTP REST API entry point.
- `cmd/loader/`: Dummy data populator script.
- `internal/analyzer/`: Text tokenization and stop-word filtering.
- `internal/index/`: Core Inverted Index implementation with concurrency safety (`sync.RWMutex`).
- `internal/query/`: Search execution and TF-IDF scoring math.
- `internal/db/`: PostgreSQL integration and index rebuilding logic.

## 📝 License
This project is created for educational and portfolio purposes, open-sourced under the MIT License.
