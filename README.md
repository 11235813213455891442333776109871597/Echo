# Echo

A lightweight HTTP echo server written in Go. It reflects every incoming request back to the caller as JSON, including the method, path, query parameters, headers, and body.

## Usage

```bash
# Build
go build -o echo .

# Run (default: :8080)
./echo

# Run on a custom address
ECHO_ADDR=:9090 ./echo
```

### Example

```bash
curl -s -X POST 'http://localhost:8080/hello?foo=bar' \
  -H 'Content-Type: application/json' \
  -d '{"greeting":"world"}'
```

```json
{
  "method": "POST",
  "path": "/hello",
  "query": {"foo": ["bar"]},
  "headers": {"Content-Type": ["application/json"]},
  "body": "{\"greeting\":\"world\"}"
}
```

## Development

```bash
# Run tests
go test ./...

# Run with race detector
go test -race ./...
```