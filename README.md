# 🔗 URL Shortener

A production-inspired URL Shortener built with **Go** and **PostgreSQL**, following Clean Architecture principles.

The project is being developed as part of a backend engineering journey focused on learning scalable software architecture rather than simply building CRUD applications.

---

## ✨ Features

### ✅ URL Management

- Create short URLs
- Retrieve URL details by Short ID
- Redirect users to the original URL

### ✅ Backend Architecture

- Clean Architecture
- Repository Pattern
- Service Layer
- Dependency Injection
- Application-level Error Handling

### ✅ Database

- PostgreSQL
- pgx Connection Pool
- SQL Migrations

---

## 🏗️ Architecture

```
                 HTTP Request
                       │
                       ▼
                HTTP Handler
                       │
                       ▼
                Service Layer
                       │
                       ▼
              Repository Layer
                       │
                       ▼
                 PostgreSQL
```

Each layer has a single responsibility.

| Layer | Responsibility |
|--------|----------------|
| Handler | HTTP request/response |
| Service | Business logic |
| Repository | Database operations |
| Database | Data persistence |

---

## 📁 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── apperrors/
│   ├── config/
│   ├── database/
│   ├── dto/
│   ├── generator/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── service/
│
├── migrations/
│
├── .env
├── go.mod
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 18+
- Git

---

### Clone the repository

```bash
git clone https://github.com/<your-username>/URL-Shortener.git

cd URL-Shortener
```

---

### Configure Environment Variables

Create a `.env` file.

```env
DATABASE_URL=postgres://username:password@localhost:5432/url_shortener

PORT=8080
```

---

### Create the Database

```sql
CREATE DATABASE url_shortener;
```

---

### Enable UUID Extension

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

---

### Run Migration

Execute:

```
migrations/001_create_urls_table.sql
```

using pgAdmin or psql.

---

### Install Dependencies

```bash
go mod tidy
```

---

### Start the Server

```bash
go run ./cmd/server
```

Server starts on

```
http://localhost:8080
```

---

## 📌 API

---

### Health Check

```
GET /health
```

Response

```json
{
    "status": "ok"
}
```

---

### Create Short URL

```
POST /api/urls
```

Request

```json
{
    "url":"https://www.google.com"
}
```

Response

```json
{
    "id":"6z96Gc",
    "short_url":"http://localhost:8080/6z96Gc"
}
```

---

### Retrieve URL

```
GET /api/urls/{shortID}
```

Example

```
GET /api/urls/6z96Gc
```

Response

```json
{
    "id":"...",
    "short_id":"6z96Gc",
    "original_url":"https://www.google.com",
    "created_at":"2026-08-06T19:31:36Z"
}
```

---

### Redirect

```
GET /{shortID}
```

Example

```
GET /6z96Gc
```

Response

```
302 Found
Location: https://www.google.com
```

---

## 🛣️ Roadmap

### ✅ v0.1.0

- HTTP Server
- PostgreSQL Integration
- Connection Pool
- URL Creation
- URL Retrieval
- URL Redirect

### 🚧 Upcoming

- Duplicate URL Detection
- Click Analytics
- Custom Short URLs
- URL Expiration
- Redis Cache
- Authentication
- Docker
- CI/CD Pipeline
- API Documentation
- Deployment

---

## 🧪 Tech Stack

- Go
- PostgreSQL
- pgx
- Standard Library HTTP Router

---

## 🎯 Learning Goals

This project emphasizes:

- Clean Architecture
- Layered Backend Design
- Dependency Injection
- Repository Pattern
- Database Design
- Production-Oriented Backend Development

---

## 📄 License

MIT License.
