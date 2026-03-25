# Echo

A lightweight Go HTTP server that echoes back the request's method, headers, query parameters, and body. It also provides JWT-based user authentication.

## Features

- **Echo endpoint** – reflects method, headers, query params, and body as JSON.
- **Login endpoint** – issues a signed JWT on successful authentication.
- **Secure echo endpoint** – protected by Bearer token middleware.
- **Algorithm-confusion fix** – `ValidateToken` enforces HMAC signing and rejects tokens that use the `none` algorithm or any other unexpected signing method.

## Endpoints

| Method | Path           | Auth required | Description                        |
|--------|----------------|---------------|------------------------------------|
| ANY    | `/echo`        | No            | Echo request details               |
| POST   | `/login`       | No            | Obtain a JWT (see credentials)     |
| ANY    | `/echo/secure` | Yes (Bearer)  | Echo request details (protected)   |

## Quick start

```bash
go run .
# Server starts on :8080 (override with PORT env var)
```

### Login

```bash
curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"password123"}'
# {"token":"<jwt>"}
```

### Echo (public)

```bash
curl -s "http://localhost:8080/echo?foo=bar" \
  -H 'X-Custom: hello' \
  -d 'hello world'
```

### Echo (protected)

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"password123"}' | jq -r .token)

curl -s http://localhost:8080/echo/secure \
  -H "Authorization: Bearer $TOKEN"
```

## Running tests

```bash
go test ./...
```

## Security notes

- The JWT secret (`jwtSecret` in `auth/auth.go`) is a placeholder. In production, load it from a secrets manager or environment variable.
- Default credentials are hard-coded for demonstration only; replace with a proper credential store in production.
