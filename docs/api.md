# Chirpy API reference

Base URL for local development: `http://localhost:8080`.

JSON error responses have the shape:

```json
{"error":"message"}
```

Unless stated otherwise, endpoints do not require a trailing slash. The server uses Go's method-aware `http.ServeMux` routes.

## Health and admin

### `GET /api/healthz`

Checks that the HTTP server is responding.

**Success — `200 OK`**

Content type is `text/html; charset=utf-8`; body is the plain text `OK`.

### `GET /admin/metrics`

Returns an HTML page showing the number of requests served through `/app/` since process start or the last successful reset. Requests to other routes do not increment the counter.

**Success — `200 OK`** (`text/html; charset=utf-8`)

### `POST /admin/reset`

Deletes all users from the database, which also deletes associated chirps and refresh tokens through foreign-key cascades, and resets the `/app/` hit counter.

No request body is read. Only enabled when the configured `PLATFORM` value is exactly `dev`.

**Success — `200 OK`** (`text/plain; charset=utf-8`): `Reset complete`

**Errors:** `403 Forbidden` when platform is not `dev`; database failures produce a `500` JSON error. On a database failure, the current handler writes the error and then continues into its success response code, so the response body may contain both the error JSON and `Reset complete`.

## Users and authentication

### `POST /api/users`

Creates a user.

**Request — `application/json`**

```json
{
  "email": "reader@example.com",
  "password": "correct horse battery staple"
}
```

**Success — `201 Created`**

Returns a user object. `token` and `refresh_token` are empty strings on this response.

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z",
  "email": "reader@example.com",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false
}
```

**Errors:** `400 Bad Request` for malformed JSON; `401 Unauthorized` if password hashing fails; `500 Internal Server Error` if the user cannot be inserted (for example, a duplicate email).

### `POST /api/login`

Authenticates a user and issues an access JWT and refresh token.

**Request — `application/json`**

```json
{
  "email": "reader@example.com",
  "password": "correct horse battery staple"
}
```

**Success — `200 OK`**

Returns the user object with `token` and `refresh_token` populated. The access token is valid for one hour; the refresh token is valid for 60 days.

**Errors:** `400 Bad Request` for malformed JSON; `401 Unauthorized` for an unknown email or incorrect password; `500 Internal Server Error` if token creation fails.

### `POST /api/refresh`

Issues a new one-hour access JWT using a valid, unexpired, non-revoked refresh token.

**Request headers**

```http
Authorization: Bearer <refresh-token>
```

No request body is used.

**Success — `200 OK`**

```json
{"token":"<new-access-jwt>"}
```

**Errors:** `400 Bad Request` if the Authorization header is missing; `401 Unauthorized` if the refresh token cannot be used; `500 Internal Server Error` if a new JWT cannot be created.

### `POST /api/revoke`

Revokes a refresh token.

**Request headers**

```http
Authorization: Bearer <refresh-token>
```

No request body is used.

**Success — `204 No Content`**

**Errors:** `400 Bad Request` if the Authorization header is missing; `500 Internal Server Error` if the database update fails. The handler returns success even if the supplied token matched no database row.

### `PUT /api/users`

Updates the authenticated user's email and password. The access JWT identifies the user; no user ID is sent in the URL or body.

**Request headers**

```http
Authorization: Bearer <access-token>
Content-Type: application/json
```

**Request body**

```json
{
  "email": "new-address@example.com",
  "password": "new password"
}
```

**Success — `200 OK`**

Returns the updated user object. `token` and `refresh_token` are empty strings.

**Errors:** `401 Unauthorized` for missing or invalid access JWT; `400 Bad Request` for malformed JSON or a failed database update.

## Chirps

The chirp object returned by create and read endpoints has this shape:

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z",
  "body": "Hello, Chirpy!",
  "user_id": "00000000-0000-0000-0000-000000000000"
}
```

### `POST /api/chirps`

Creates a chirp for the authenticated user.

**Request headers**

```http
Authorization: Bearer <access-token>
Content-Type: application/json
```

**Request body**

```json
{"body":"Hello, Chirpy!"}
```

The body is limited to 140 characters. The words `kerfuffle`, `sharbert` and `fornax` are replaced with `****` when they appear as whole words separated by spaces; matching is case-insensitive.

**Success — `201 Created`**

Returns the created chirp object with the cleaned body.

**Errors:** `400 Bad Request` for malformed JSON or a body longer than 140 characters; `401 Unauthorized` if the access JWT is missing or invalid; `500 Internal Server Error` if insertion fails.

### `GET /api/chirps`

Lists chirps in ascending creation-time order. No authentication is required.

Optional query parameter:

| Parameter | Meaning |
| --- | --- |
| `author_id` | UUID; when present, only chirps belonging to that user are returned |

Example: `GET /api/chirps?author_id=00000000-0000-0000-0000-000000000000`

**Success — `200 OK`**

Returns a JSON array of chirp objects. An empty result is `[]`.

**Errors:** `400 Bad Request` if `author_id` is not a UUID; `500 Internal Server Error` if the query fails.

### `GET /api/chirps/{chirpID}`

Returns one chirp by UUID. No authentication is required.

**Success — `200 OK`** with a chirp object.

**Errors:** `400 Bad Request` if `chirpID` is not a UUID; `404 Not Found` if the chirp cannot be retrieved.

### `DELETE /api/chirps/{chirpID}`

Deletes a chirp. The authenticated user must own the chirp.

**Request headers**

```http
Authorization: Bearer <access-token>
```

**Success — `204 No Content`**

**Errors:** `400 Bad Request` if `chirpID` is not a UUID; `401 Unauthorized` if the access JWT is missing or invalid; `403 Forbidden` if the chirp belongs to a different user; `404 Not Found` if the chirp cannot be retrieved; `500 Internal Server Error` if deletion fails.

## Polka webhook

### `POST /api/polka/webhooks`

Accepts Polka user upgrade events. The request must carry the configured API key.

**Request headers**

```http
Authorization: ApiKey <POLKA_KEY>
Content-Type: application/json
```

**Request body**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "00000000-0000-0000-0000-000000000000"
  }
}
```

For `user.upgraded`, sets `is_chirpy_red` to `true` for the user. Other event names are acknowledged without changing a user.

**Success — `204 No Content`**

**Errors:** `401 Unauthorized` if the key is missing or invalid; `400 Bad Request` for malformed JSON; `404 Not Found` if an upgraded user ID does not exist; `500 Internal Server Error` if the update fails.

Implementation note: with a present but incorrect API key, the current handler writes a `401` response but does not return before decoding and processing the event. Do not rely on that request being side-effect free until this behavior is corrected.

## Static application

### `GET /app/` and other `/app/*` paths

Serves files from the process working directory with `/app` stripped from the requested path. For example, `/app/` serves `index.html`. These requests increment the counter displayed at `/admin/metrics`.
