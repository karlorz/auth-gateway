# AGENTS.md

## Project Overview

New API is a next-generation AI model gateway and asset management system, forked from One API. It provides a unified API interface for multiple LLM providers (OpenAI, Claude, Gemini, etc.) with features like user management, token distribution, quota tracking, billing, and rate limiting.

**Tech Stack:**
- Backend: Go 1.25.1 with Gin web framework
- Frontend: React 18 with Vite, Semi UI components
- Database: SQLite (default), MySQL 5.7.8+, or PostgreSQL 9.6+
- Cache: Redis (optional but recommended)
- Package Manager: Bun (frontend)

## Building and Running

### Development Commands

**Backend (Go):**
```bash
# Run backend server
go run main.go

# Build binary
go build -o new-api main.go
```

**Frontend (React + Vite):**
```bash
cd web
bun install              # Install dependencies
bun run dev              # Start dev server (port 5173 by proxy to :3000)
bun run build            # Production build
bun run lint             # Check code formatting with Prettier
bun run lint:fix         # Fix formatting issues
bun run eslint           # Run ESLint
bun run eslint:fix       # Fix ESLint issues
```

**Full Stack Build:**
```bash
make all                 # Build frontend and start backend
```

### Running with Docker

**Docker Compose (Recommended):**
```bash
docker-compose up -d     # Start all services (app + postgres + redis)
docker-compose logs -f   # View logs
docker-compose down      # Stop all services
```

**Direct Docker:**
```bash
# SQLite
docker run --name new-api -d --restart always -p 3000:3000 \
  -e TZ=Asia/Shanghai -v /home/ubuntu/data/new-api:/data \
  calciumion/new-api:latest

# MySQL
docker run --name new-api -d --restart always -p 3000:3000 \
  -e SQL_DSN="root:123456@tcp(localhost:3306)/oneapi" \
  -e TZ=Asia/Shanghai -v /home/ubuntu/data/new-api:/data \
  calciumion/new-api:latest
```

## Architecture

### High-Level Structure

```
main.go                    # Entry point, initializes DB/Redis/cache, starts Gin server
├── router/                # Route definitions (API, relay, dashboard, web)
├── controller/            # HTTP handlers for API endpoints
├── middleware/            # Auth, rate limiting, CORS, distributor logic
├── relay/                 # Core proxy logic for forwarding requests to LLM providers
│   ├── channel/           # Provider-specific adapters (openai, claude, gemini, etc.)
│   └── helper/            # Request/response transformation utilities
├── model/                 # Database models (GORM) and caching logic
├── service/               # Business logic (quota, billing, notifications, token counting)
├── dto/                   # Data transfer objects for API requests/responses
├── common/                # Shared utilities and constants
├── constant/              # Application constants
└── web/                   # React frontend
    └── src/
        ├── pages/         # Page components
        ├── components/    # Reusable UI components
        └── helpers/       # Frontend utilities
```

### Request Flow for LLM Proxying

1. **Client Request** → `relay-router.go` (router)
2. **Authentication** → `middleware/auth.go` validates token
3. **Rate Limiting** → `middleware/rate-limit.go` or `model-rate-limit.go`
4. **Channel Selection** → `middleware/distributor.go` selects best provider channel
5. **Adapter Pattern** → `relay/relay_adaptor.go` calls provider-specific adapter
6. **Provider Channels** → `relay/channel/{provider}/` handles provider-specific logic
7. **Request Transformation** → Convert OpenAI format to provider format if needed
8. **Forward to Provider** → HTTP client in `service/http_client.go`
9. **Response Transformation** → Convert provider response back to OpenAI format
10. **Billing & Logging** → `service/quota.go` updates usage, logs to DB

### Channel Adapter System

Each LLM provider has its own adapter in `relay/channel/{provider}/`:
- `adaptor.go`: Implements the `channel.Adaptor` interface
- `main.go`: Provider-specific request/response handling
- Additional files for embeddings, images, audio, etc.

The adaptor interface abstracts away provider differences, allowing the gateway to treat all providers uniformly. The `GetAdaptor()` function in `relay/relay_adaptor.go` returns the appropriate adapter based on channel type.

### Key Middleware Components

- **distributor.go**: Core routing logic - selects channels based on weights, priorities, and model availability. Handles retries on failure.
- **auth.go**: Token validation and user authentication
- **model-rate-limit.go**: Per-user, per-model rate limiting
- **rate-limit.go**: Global API rate limiting

### Database & Caching

- **model/**: GORM models for channels, tokens, users, logs, quotas, etc.
- **channel_cache.go**: In-memory cache for channel configs (synced from DB)
- **option.go**: System-wide settings cached in memory
- Redis is optional but recommended for multi-instance deployments

## Key Configuration

### Environment Variables (.env)

Critical settings:
- `SQL_DSN`: Database connection string
- `REDIS_CONN_STRING`: Redis connection (e.g., `redis://localhost:6379/0`)
- `SESSION_SECRET`: **Required for multi-node deployments**
- `CRYPTO_SECRET`: **Required for multi-node deployments with shared Redis**
- `STREAMING_TIMEOUT`: Stream timeout in seconds (default: 300)
- `SYNC_FREQUENCY`: Cache sync interval in seconds (default: 60)
- `BATCH_UPDATE_ENABLED=true`: Enable batch DB updates for performance
- `UPDATE_TASK=true`: Enable async task updates (Midjourney, Suno)
- `ERROR_LOG_ENABLED=true`: Enable error logging

See `.env.example` for complete list.

### Multi-Instance Deployment

When running multiple instances:
1. Set `SESSION_SECRET` to same value on all nodes
2. Set `CRYPTO_SECRET` to same value if using shared Redis
3. Use Redis for caching (`REDIS_CONN_STRING`)
4. Use remote database (MySQL/PostgreSQL, not SQLite)

## Testing

Currently, the codebase has minimal test coverage. When adding tests:
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./model -run TestChannelCache
```

## Common Development Tasks

### Adding a New LLM Provider

1. Create new directory in `relay/channel/{provider}/`
2. Implement the `channel.Adaptor` interface
3. Add provider constant to `constant/` (channel types)
4. Register adapter in `relay/relay_adaptor.go` `GetAdaptor()` function
5. Add provider-specific request/response structs to `dto/`
6. Add pricing info to `model/pricing_default.go` if needed

### Adding a New API Endpoint

1. Define route in `router/api-router.go`
2. Create handler in `controller/`
3. Add middleware if auth/rate limiting needed
4. Update corresponding model in `model/` if DB changes required
5. Add DTO in `dto/` for request/response validation

### Modifying Billing Logic

Billing is handled in `service/quota.go`. Key functions:
- `ConsumeQuota()`: Deduct quota from user
- `PostConsumeQuota()`: Record usage after streaming completes
- Token counting uses tiktoken-go in `service/token_counter.go`

### Format Conversion

Request format conversion (e.g., OpenAI → Claude) is in `service/convert.go`. Enable per-channel in the channel settings.

## Documentation

Official documentation: https://docs.newapi.pro/

Key docs:
- Feature introduction: https://docs.newapi.pro/wiki/features-introduction
- Installation guide: https://docs.newapi.pro/installation
- API reference: https://docs.newapi.pro/api
- Environment variables: https://docs.newapi.pro/installation/environment-variables

## Important Notes

- Default admin username: `root`, default password: `123456` (change immediately!)
- Frontend runs on port 3000 (proxied through backend)
- API endpoints are under `/api/`
- Relay endpoints are under `/v1/` (OpenAI-compatible)
- When modifying channel adapters, ensure proper error handling and quota refunds on failure
- Always test with `BATCH_UPDATE_ENABLED=false` first to avoid DB race conditions during development
