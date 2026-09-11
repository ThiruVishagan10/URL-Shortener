# URL Shortener

A Go HTTP API for creating, managing, and resolving shortened URLs. Data is stored in PostgreSQL; users sign in through Google OAuth 2.0 and are authenticated with server-side sessions held in an HTTP-only cookie.

The service is intended to be used with a separate frontend. Requests from that frontend are allowed only when their `Origin` exactly matches `FRONTEND_URL`.

## What it does

- Creates six-character, URL-safe short IDs for authenticated users.
- Redirects public short URLs to their original destinations.
- Supports `public` and `private` URL visibility.
- Lets each user list, change the visibility of, and delete their own URLs.
- Reuses a user's existing short URL when they submit the same original URL again.
- Signs users in with Google OAuth 2.0 using PKCE and CSRF state validation.
- Stores sessions in PostgreSQL with a 24-hour expiry.
- Rate-limits URL creation and the start of Google sign-in to 10 requests per second per client IP, with a burst of 20.

## Architecture

```
HTTP request
  -> middleware (CORS, rate limiting, authentication)
  -> handler (HTTP input and output)
  -> service (business rules and access control)
  -> repository (SQL)
  -> PostgreSQL
```

The repository uses the standard library `net/http` router introduced in Go 1.22 and `pgx` for PostgreSQL access.

## Project layout

```text
cmd/server/          Application entry point and route registration
internal/auth/       Google OAuth provider and identity validation
internal/config/     Environment-based configuration
internal/database/   PostgreSQL connection-pool setup
internal/handler/    HTTP handlers
internal/middleware/ CORS, session authentication, and IP rate limiting
internal/model/      URL, user, session, and visibility models
internal/repository/ PostgreSQL queries
internal/service/    URL, user, and session business logic
migrations/          PostgreSQL schema migrations
```

## Prerequisites

- Go 1.26.4 or later (the version declared in `go.mod`)
- PostgreSQL 18 or later. The schema uses `uuidv7()`.
- A Google OAuth client configured for the application callback URL.

## Configuration

Create a `.env` file in the project root. The server loads it when present; environment variables can be used instead.

```env
DATABASE_URL=postgres://username:password@localhost:5432/url_shortener
PORT=8080
FRONTEND_URL=http://localhost:3000
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
APP_BASE_URL=http://localhost:8080
```

All seven values are required. Do not commit `.env` or OAuth credentials.

`FRONTEND_URL` has two roles: the root route redirects there, and it is the single allowed CORS origin. Its value must exactly match the frontend browser origin, including scheme, host, and port.

`APP_BASE_URL` is the public address used in created short URLs. Do not include a trailing slash. For the deployed API, set it to `https://url-shortener-4ykd.onrender.com` in Render's environment settings.

In Google Cloud, add the exact value of `GOOGLE_REDIRECT_URL` to the OAuth client's authorized redirect URIs. Use either `localhost` or `127.0.0.1` consistently throughout the browser URL and configuration; they are different origins.

## Database setup

Create the database first:

```sql
CREATE DATABASE url_shortener;
```

Then apply the executable schema migrations in this order:

```text
001_create_urls_table.sql
002_add_unique_constraints.sql
003_create_users_table.sql
004_create_sessions_table.sql
005_add_user_id_to_urls.sql
006_add_visibility_to_urls.sql
008_fix_url_unique_constraint.sql
```

Migration `007_make_url_unique_per_user.sql` is currently a diagnostic query with a misspelled PostgreSQL catalog name, not a schema migration. Do not run it. Migration `008` replaces the global `original_url` uniqueness rule with a `(user_id, original_url)` constraint, which allows different users to shorten the same destination.

## Run locally

Install dependencies and start the server:

```bash
go mod download
go run ./cmd/server
```

The service listens on `http://localhost:<PORT>`. With the example configuration, that is `http://localhost:8080`.

Run the current test suite with:

```bash
go test ./...
```

## Authentication and sessions

1. Send the browser to `GET /auth/google`.
2. The server creates a random OAuth state and PKCE verifier, stores both in short-lived HTTP-only cookies, and redirects to Google.
3. Google returns to `GET /auth/google/callback`.
4. After validating the callback and identity token, the server creates or finds the user, persists a session, and sets the `session_id` cookie.
5. The callback redirects to `/`, which then redirects to `FRONTEND_URL`.

Sessions expire after 24 hours. The OAuth state and verifier cookies expire after five minutes and are removed after a successful callback. The current cookies use `SameSite=Lax` and `Secure=false`, which is appropriate only for local HTTP development; configure secure cookies before a production deployment.

Browser clients must send credentials for authenticated cross-origin API calls, for example with `fetch(..., { credentials: "include" })`.

## API

Unless otherwise stated, errors are plain-text HTTP error responses. Protected endpoints require a valid `session_id` cookie.

### `GET /`

Redirects (`302 Found`) to `FRONTEND_URL`.

### `GET /health`

Returns service health without authentication.

```json
{"status":"healthy"}
```

### `GET /auth/google`

Begins Google OAuth. This route is rate-limited and redirects the browser to Google.

### `GET /auth/google/callback`

Completes the OAuth flow. It requires the `code` and `state` query parameters and the short-lived OAuth cookies set by `/auth/google`.

### `POST /api/auth/logout`

Deletes the current server-side session when a `session_id` cookie is present and clears that cookie. Returns `204 No Content`. Calling it with no session cookie also returns `204 No Content`.

### `GET /api/me`

Requires authentication. Returns plain text in the current implementation:

```text
Authenticated user: <user-id>
```

### `POST /api/urls`

Requires authentication and is rate-limited. Create a URL with a valid absolute URL (a scheme and host are required) and an explicit visibility value.

```http
POST /api/urls
Content-Type: application/json

{"url":"https://example.com/articles/123","visibility":"public"}
```

Valid visibility values are `public` and `private`. The current handler requires `visibility` to be supplied. When the authenticated user has already shortened the same original URL, the existing record is returned instead of creating another one.

Successful response (`201 Created`):

```json
{"id":"aB3dE9","short_url":"<APP_BASE_URL>/aB3dE9"}
```

The `short_url` response uses `APP_BASE_URL`.

### `GET /api/urls`

Requires authentication. Lists URLs owned by the current user, newest first.

```json
{
  "urls": [
    {
      "id": "uuid-value",
      "short_id": "aB3dE9",
      "original_url": "https://example.com/articles/123",
      "visibility": "public",
      "created_at": "2026-09-11T10:00:00Z"
    }
  ]
}
```

### `GET /api/urls/{shortID}`

Returns details for a public URL to anyone. A private URL is returned only to its owner, who must send their session cookie. Unavailable, private-to-another-user, and nonexistent URLs return `404 Not Found`.

### `PATCH /api/urls/{shortID}`

Requires authentication and ownership. Changes a URL's visibility.

```http
PATCH /api/urls/aB3dE9
Content-Type: application/json

{"visibility":"private"}
```

Successful response:

```json
{"short_id":"aB3dE9","visibility":"private"}
```

### `DELETE /api/urls/{shortID}`

Requires authentication and ownership. Deletes the URL and returns `204 No Content`. A URL not owned by the caller returns `404 Not Found`.

### `GET /{shortID}`

Redirects a public short URL with `302 Found` to its original URL. Private URLs redirect only for their owner; others receive `404 Not Found`.

## Current operational notes

- The in-memory rate-limiter map does not expire inactive client IP entries. It resets whenever the process restarts and is not shared across multiple server instances.
- The redirect and URL-details routes use optional authentication. A missing or invalid session does not expose private URLs.
- URL ownership is enforced for listing, editing, and deletion.
- There is no container setup, CI pipeline, custom aliases, expiry for short URLs, click analytics, or configurable public base URL yet.

## Technology

- Go and the standard library HTTP server
- PostgreSQL and `pgx/v5`
- Google OAuth 2.0 and OpenID Connect
- `golang.org/x/time/rate` for token-bucket rate limiting

## License

MIT
