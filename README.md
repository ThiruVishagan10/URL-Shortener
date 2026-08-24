# URL Shortener

A Go and PostgreSQL URL shortener built with a clean, layered architecture.

This project includes:

- Short URL creation and redirect handling
- Google OAuth login
- Session-based authentication
- A PostgreSQL-backed persistence layer
- Migration-based schema management

## Features

- Create short URLs for authenticated users
- Reuse an existing short URL when the same original URL is submitted again
- Resolve short IDs to the original URL
- Redirect short IDs to the original destination
- Sign in with Google and maintain a session cookie
- Fetch the currently authenticated user with `/api/me`

## Architecture

The application follows a simple request flow:

`HTTP Handler -> Service Layer -> Repository Layer -> PostgreSQL`

Each layer has a single responsibility:

- Handler: request parsing, validation, and HTTP responses
- Service: business rules and orchestration
- Repository: database access
- Database: persistence

## Project Structure

```text
.
./cmd/server/main.go
./internal/apperrors/
./internal/auth/
./internal/config/
./internal/database/
./internal/generator/
./internal/handler/
./internal/middleware/
./internal/model/
./internal/repository/
./internal/service/
./docs/
./migrations/
./go.mod
./README.md
```

## Requirements

- Go 1.22 or newer
- PostgreSQL 18 or newer
- Google OAuth credentials

## Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://username:password@localhost:5432/url_shortener
PORT=8080
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

The application requires all of these values at startup.

## Database Setup

1. Create the database.

```sql
CREATE DATABASE url_shortener;
```

2. Make sure the UUID functions used by the migrations are available in your PostgreSQL version.

3. Run the SQL files in `migrations/` in order:

- `001_create_urls_table.sql`
- `002_add_unique_constraints.sql`
- `003_create_users_table.sql`
- `004_create_sessions_table.sql`
- `005_add_user_id_to_urls.sql`
- `006_add_visiblity_to_urls.sql`

## Running the App

Install dependencies:

```bash
go mod tidy
```

Start the server:

```bash
go run ./cmd/server
```

The server listens on the port defined by `PORT`.

## API

### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "healthy"
}
```

### Google Login

```http
GET /auth/google
```

Starts the OAuth flow and redirects to Google.

### Google Callback

```http
GET /auth/google/callback
```

Handles the OAuth callback, creates or finds the user, and sets a `session_id` cookie.

### Current User

```http
GET /api/me
```

Requires a valid session cookie.

### Create Short URL

```http
POST /api/urls
```

Requires authentication.

Request:

```json
{
  "url": "https://www.google.com"
}
```

Response:

```json
{
  "id": "6z96Gc",
  "short_url": "http://localhost:8080/6z96Gc"
}
```

Notes:

- The request body must contain a valid absolute URL.
- If the same original URL is submitted again, the existing short URL is returned.

### Retrieve URL Details

```http
GET /api/urls/{shortID}
```

Example:

```http
GET /api/urls/6z96Gc
```

Response:

```json
{
  "id": "uuid-value",
  "short_id": "6z96Gc",
  "original_url": "https://www.google.com",
  "created_at": "2026-08-06T19:31:36Z"
}
```

### Redirect

```http
GET /{shortID}
```

Example:

```http
GET /6z96Gc
```

Response:

- `302 Found`
- `Location: https://www.google.com`

## Data Model

The current URL record stores:

- `id`
- `short_id`
- `original_url`
- `user_id`
- `visibility`
- `created_at`

Users and sessions are also persisted for Google login and authenticated access.

## Roadmap

- Custom short URLs
- URL expiration
- Click analytics
- Redis caching
- Docker support
- CI/CD
- API documentation
- Deployment setup

## Tech Stack

- Go
- PostgreSQL
- pgx
- Standard library HTTP server
- Google OAuth

## License

MIT
