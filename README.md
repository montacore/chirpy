# Chirpy

Chirpy is a small Go HTTP service for user accounts and short posts (“chirps”). It uses PostgreSQL for persistent data, JWT access tokens and refresh tokens for authentication, and serves a small static page under `/app/`.

## Requirements

- Go version declared by the module: `1.26.5`
- PostgreSQL with the Chirpy schema applied
- A `.env` file in the working directory

The server reads these settings from the process environment (the `.env` file is loaded at startup):

| Variable | Purpose |
| --- | --- |
| `DB_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret used to sign access tokens; required |
| `PLATFORM` | Set to `dev` to enable `POST /admin/reset` |
| `POLKA_KEY` | API key expected by the Polka webhook endpoint |

The application listens on port `8080`. It opens a PostgreSQL connection but does not run schema migrations itself. Apply the SQL files in `sql/schema/` using your migration workflow before starting the service.

## Run

Create `.env` with the settings above, then run:

```sh
go run .
```

The server starts on `http://localhost:8080`. Startup exits if `.env` cannot be loaded or `JWT_SECRET` is empty.

## Project structure

| Path | Description |
| --- | --- |
| `main.go` | Loads configuration, constructs the database and registers HTTP routes |
| `handler_*.go`, `handle_*.go` | HTTP endpoint handlers |
| `middleware.go` | Shared API configuration, file hit counter and admin handlers |
| `handlerValidate.go` | JSON/error response helpers and chirp profanity masking |
| `chirp.go`, `user.go` | JSON-facing Chirp and User response models |
| `internal/auth/` | Password hashing, JWT validation, bearer/API key parsing and refresh token generation |
| `internal/database/` | sqlc-generated Go database access code |
| `sql/queries/` | SQL query source used to generate database access code |
| `sql/schema/` | Ordered database schema migrations |
| `assets/`, `index.html` | Static site content served from the application directory |
| `docs/api.md` | Endpoint reference, payloads, authentication and status codes |

## API overview

The service exposes health, admin, user/authentication, chirp and Polka webhook routes. Most API responses use JSON; successful delete, revoke and webhook acknowledgement responses have no body. See [the API reference](docs/api.md) for request and response details.

Access protected endpoints with `Authorization: Bearer <access-token>`. The refresh and revoke routes expect the refresh token in the same Bearer header format. The Polka webhook uses `Authorization: ApiKey <key>`.

## Database and authentication notes

- New passwords are stored as Argon2id hashes.
- Access JWTs expire after one hour. Login-issued refresh tokens expire after 60 days and can be revoked.
- Chirp text is limited to 140 characters and the words `kerfuffle`, `sharbert` and `fornax` are masked case-insensitively when they match a space-delimited word exactly.
- Deleting a user cascades to that user's chirps and refresh tokens according to the schema.
- `POST /admin/reset` deletes all users (and cascading rows) and resets the in-memory file-server hit count. It is enabled only when `PLATFORM=dev`.
