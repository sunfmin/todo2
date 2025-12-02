# API Contracts

This directory contains the API contract definitions for the Simple Todo App.

## Files

### todo.proto
Protocol Buffer definition for the Todo API. This is the source of truth for:
- Data structures (Todo, requests, responses)
- Field validation rules
- Type safety

**Usage**:
```bash
# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --validate_out="lang=go:." --validate_opt=paths=source_relative \
       todo.proto
```

### openapi.yaml
OpenAPI 3.0 specification for the REST API. Provides:
- HTTP endpoint definitions
- Request/response schemas
- Error responses
- API documentation

**Usage**:
```bash
# View in Swagger UI
docker run -p 8080:8080 -e SWAGGER_JSON=/contracts/openapi.yaml \
           -v $(pwd):/contracts swaggerapi/swagger-ui

# Generate client code (optional)
openapi-generator-cli generate -i openapi.yaml -g go -o client/
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/todos` | Create a new todo |
| GET | `/api/v1/todos` | List todos (paginated) |
| GET | `/api/v1/todos/{id}` | Get a single todo |
| PUT | `/api/v1/todos/{id}` | Update a todo |
| DELETE | `/api/v1/todos/{id}` | Delete a todo |
| GET | `/health` | Health check |

## Request/Response Examples

### Create Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"description": "Buy groceries"}'
```

Response (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Buy groceries",
  "completed": false,
  "created_at": "2025-12-02T03:00:00Z",
  "updated_at": "2025-12-02T03:00:00Z"
}
```

### List Todos
```bash
curl http://localhost:8080/api/v1/todos?limit=20&offset=0
```

Response (200 OK):
```json
{
  "todos": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "description": "Buy groceries",
      "completed": false,
      "created_at": "2025-12-02T03:00:00Z",
      "updated_at": "2025-12-02T03:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

### Update Todo
```bash
curl -X PUT http://localhost:8080/api/v1/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

Response (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Buy groceries",
  "completed": true,
  "created_at": "2025-12-02T03:00:00Z",
  "updated_at": "2025-12-02T03:05:00Z"
}
```

### Delete Todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/550e8400-e29b-41d4-a716-446655440000
```

Response (204 No Content)

## Validation Rules

### Description Field
- **Required**: Yes (for create)
- **Min Length**: 1 character
- **Max Length**: 500 characters
- **Pattern**: Cannot start with whitespace
- **Special Characters**: Allowed (including emojis)

### ID Field
- **Format**: UUID v4
- **Example**: `550e8400-e29b-41d4-a716-446655440000`

### Pagination
- **Limit**: 1-100 (default: 20)
- **Offset**: >= 0 (default: 0)

## Error Responses

All errors follow this format:
```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {}
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_INPUT` | 400 | Invalid request data |
| `TODO_NOT_FOUND` | 404 | Todo does not exist |
| `INTERNAL_ERROR` | 500 | Server error |

## Versioning

API version is included in the URL path: `/api/v1/`

Breaking changes will increment the major version (v2, v3, etc.).

## Testing

See `../data-model.md` for protobuf message examples and test fixtures.

## References

- [Protocol Buffers](https://protobuf.dev/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate)