# CLAUDE.md

## Project Overview

Auth Gateway is a production-ready, multi-provider authentication system forked from New API. It provides a comprehensive authentication solution that can be easily integrated into any fullstack application.

**Tech Stack:**
- Backend: Go 1.25+ with Gin web framework
- Frontend: React 18 with Vite, Semi UI, Tailwind CSS
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
go build -o auth-gateway main.go
```

**Frontend (React + Vite):**
```bash
cd web
bun install              # Install dependencies
bun run dev              # Start dev server (port 5173, proxies to :3000)
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
docker run --name auth-gateway -d --restart always -p 3000:3000 \
  -e TZ=Asia/Shanghai -v ./data:/data \
  karl8080/auth-gateway:latest

# MySQL
docker run --name auth-gateway -d --restart always -p 3000:3000 \
  -e SQL_DSN="root:password@tcp(localhost:3306)/auth_gateway" \
  -e TZ=Asia/Shanghai -v ./data:/data \
  karl8080/auth-gateway:latest
```

## Authentication Features

Auth Gateway supports multiple authentication methods:

| Method | Description | Key Files |
|--------|-------------|-----------|
| Password | Username/email + bcrypt hashing | `controller/user.go` |
| GitHub OAuth | OAuth 2.0 | `controller/github.go` |
| Discord OAuth | OAuth 2.0 | `controller/discord.go` |
| OIDC | OpenID Connect (enterprise SSO) | `controller/oidc.go` |
| WeChat | QR code authentication | `controller/wechat.go` |
| Telegram | Bot widget login | `controller/telegram.go` |
| LinuxDO | Community OAuth with trust levels | `controller/linuxdo.go` |
| Passkey/WebAuthn | Passwordless biometric/hardware key | `controller/passkey.go`, `service/passkey/` |
| 2FA (TOTP) | Time-based OTP with backup codes | `controller/twofa.go`, `model/twofa.go` |
| Turnstile | Cloudflare bot protection | `middleware/turnstile-check.go` |

## Architecture

### High-Level Structure

```
main.go                    # Entry point, initializes DB/Redis/cache, starts Gin server
├── router/                # Route definitions
│   └── api-router.go      # Auth API routes
├── controller/            # HTTP handlers
│   ├── user.go            # Login, Register, Logout, User CRUD
│   ├── passkey.go         # WebAuthn endpoints
│   ├── twofa.go           # 2FA setup/verify endpoints
│   ├── github.go          # GitHub OAuth callback
│   ├── discord.go         # Discord OAuth callback
│   ├── oidc.go            # OIDC callback
│   ├── telegram.go        # Telegram login
│   ├── wechat.go          # WeChat login
│   └── linuxdo.go         # LinuxDO OAuth callback
├── middleware/            # Auth middleware
│   ├── auth.go            # UserAuth, AdminAuth, RootAuth, TokenAuth
│   ├── turnstile-check.go # Cloudflare Turnstile verification
│   └── rate-limit.go      # Rate limiting
├── model/                 # Database models (GORM)
│   ├── user.go            # User model with OAuth IDs
│   ├── passkey.go         # PasskeyCredential model
│   └── twofa.go           # TwoFA, TwoFABackupCode models
├── service/               # Business logic
│   └── passkey/           # WebAuthn service
│       ├── service.go     # WebAuthn configuration
│       ├── user.go        # WebAuthnUser interface
│       └── session.go     # Session data helpers
├── setting/               # Configuration
│   └── system_setting/    # OAuth & Passkey settings
├── common/                # Shared utilities
│   ├── password.go        # bcrypt hashing
│   └── totp.go            # TOTP generation/validation
└── web/                   # React frontend
    └── src/
        ├── components/auth/   # Auth UI components
        │   ├── LoginForm.jsx
        │   ├── RegisterForm.jsx
        │   └── TwoFAVerification.jsx
        └── helpers/           # OAuth/Passkey helpers
```

### Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND (React)                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ LoginForm   │ │RegisterForm │ │ 2FAVerify   │ │ OAuth Buttons/Callbacks ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MIDDLEWARE LAYER                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ TurnstileChk│ │ RateLimit   │ │ UserAuth    │ │ AdminAuth / RootAuth    ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            CONTROLLER LAYER                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ user.go     │ │ passkey.go  │ │ twofa.go    │ │ OAuth Controllers       ││
│  │ Login/Reg   │ │ WebAuthn    │ │ TOTP/Backup │ │ github/discord/oidc/etc ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MODEL LAYER                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ User        │ │ Passkey     │ │ TwoFA       │ │ TwoFABackupCode         ││
│  │ (user.go)   │ │ Credential  │ │ (twofa.go)  │ │ (twofa.go)              ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
```

## Key Configuration

### Environment Variables (.env)

**Critical for Auth:**
```bash
# Session (REQUIRED for multi-node)
SESSION_SECRET=your-session-secret-key
CRYPTO_SECRET=your-crypto-secret-key

# Database
SQL_DSN="root:password@tcp(localhost:3306)/auth_gateway"
REDIS_CONN_STRING=redis://localhost:6379/0

# GitHub OAuth
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx

# Discord OAuth (configured via admin settings UI)
# OIDC (configured via admin settings UI)

# Telegram OAuth
TELEGRAM_BOT_TOKEN=xxx
TELEGRAM_BOT_NAME=xxx

# WeChat OAuth
WECHAT_SERVER_ADDRESS=https://wechat-service.example.com
WECHAT_SERVER_TOKEN=xxx

# LinuxDO OAuth
LINUX_DO_CLIENT_ID=xxx
LINUX_DO_CLIENT_SECRET=xxx
LINUX_DO_MINIMUM_TRUST_LEVEL=0

# Turnstile (Cloudflare)
TURNSTILE_SITE_KEY=xxx
TURNSTILE_SECRET_KEY=xxx

# Rate Limiting
CRITICAL_RATE_LIMIT=20  # per minute
```

See `.env.example` for complete list.

### Multi-Instance Deployment

When running multiple instances:
1. Set `SESSION_SECRET` to same value on all nodes
2. Set `CRYPTO_SECRET` to same value if using shared Redis
3. Use Redis for caching (`REDIS_CONN_STRING`)
4. Use remote database (MySQL/PostgreSQL, not SQLite)

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./model -run TestUser
```

## Common Development Tasks

### Adding a New OAuth Provider

1. Create `controller/{provider}.go` with OAuth callback handler
2. Add route in `router/api-router.go`:
   ```go
   api.GET("/oauth/{provider}", controller.{Provider}OAuth)
   ```
3. Add settings in `setting/system_setting/{provider}.go`
4. Add OAuth ID field to `model/user.go`:
   ```go
   {Provider}ID string `gorm:"index"`
   ```
5. Add frontend button in `web/src/components/auth/LoginForm.jsx`
6. Add helper function in `web/src/helpers/index.js`

### Adding a New Auth API Endpoint

1. Define route in `router/api-router.go`
2. Create handler in `controller/`
3. Add middleware if auth required:
   ```go
   userRoute.POST("/endpoint", middleware.UserAuth(), controller.Handler)
   ```
4. Update model in `model/` if DB changes required

### Modifying Passkey/WebAuthn

Key files:
- `service/passkey/service.go` - WebAuthn configuration (RPID, origins)
- `setting/system_setting/passkey.go` - PasskeySettings struct
- `controller/passkey.go` - All passkey endpoints

### Modifying 2FA

Key files:
- `model/twofa.go` - TwoFA and TwoFABackupCode models
- `controller/twofa.go` - Setup, enable, disable, verify endpoints
- `common/totp.go` - TOTP generation/validation utilities

## API Endpoints Reference

### Public Endpoints (No Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/user/register` | Password registration |
| POST | `/api/user/login` | Password login |
| POST | `/api/user/login/2fa` | 2FA verification during login |
| POST | `/api/user/passkey/login/begin` | Initiate passkey login |
| POST | `/api/user/passkey/login/finish` | Complete passkey login |
| GET | `/api/oauth/state` | Generate OAuth state token |
| GET | `/api/oauth/github` | GitHub OAuth callback |
| GET | `/api/oauth/discord` | Discord OAuth callback |
| GET | `/api/oauth/oidc` | OIDC callback |
| GET | `/api/oauth/linuxdo` | LinuxDO OAuth callback |
| GET | `/api/oauth/wechat` | WeChat login |
| GET | `/api/oauth/telegram/login` | Telegram login |

### Authenticated Endpoints (User Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/user/logout` | Logout |
| GET | `/api/user/self` | Get current user info |
| PUT | `/api/user/self` | Update current user |
| GET | `/api/user/passkey` | Get passkey status |
| POST | `/api/user/passkey/register/begin` | Begin passkey registration |
| POST | `/api/user/passkey/register/finish` | Complete passkey registration |
| DELETE | `/api/user/passkey` | Remove passkey |
| GET | `/api/user/2fa/status` | Get 2FA status |
| POST | `/api/user/2fa/setup` | Initialize 2FA setup |
| POST | `/api/user/2fa/enable` | Enable 2FA |
| POST | `/api/user/2fa/disable` | Disable 2FA |
| POST | `/api/user/2fa/backup_codes` | Regenerate backup codes |

### Admin Endpoints (Admin Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/user/` | List all users |
| GET | `/api/user/search` | Search users |
| POST | `/api/user/` | Create user |
| PUT | `/api/user/` | Update user |
| DELETE | `/api/user/:id` | Delete user |
| POST | `/api/user/manage` | Manage user (enable/disable/promote) |
| DELETE | `/api/user/:id/reset_passkey` | Admin reset user's passkey |
| DELETE | `/api/user/:id/2fa` | Admin disable user's 2FA |

## Documentation

- **Full Auth System Guide**: [docs/AUTH.md](docs/AUTH.md) - Comprehensive documentation covering architecture, all auth methods, database schema, and package extraction guide
- **Original New API Docs**: https://docs.newapi.pro/

## Important Notes

- Default admin username: `root`, default password: `123456` (change immediately!)
- Frontend dev server runs on port 5173, proxies API to backend on port 3000
- API endpoints are under `/api/`
- Session middleware must be configured before auth routes
- Always use HTTPS in production (required for Passkey/WebAuthn)

## Roadmap

- [x] Trim non-auth components (relay/, channel controllers, billing) - **IN PROGRESS**
  - Removed: relay/, LLM controllers, token/channel models, billing services
  - Removed: non-auth frontend pages (Channel, Token, TopUp, Log, Midjourney, etc.)
  - Remaining: Fix remaining import errors in misc.go, option.go controllers
- [ ] Complete backend build cleanup (fix remaining undefined references)
- [ ] Create standalone Go auth SDK
- [ ] Add more OAuth providers (Google, Apple, Microsoft)
- [ ] Magic link authentication
- [ ] OAuth token refresh handling

## Current Trimming Status

**Files Removed:**
- `relay/` - Entire LLM relay folder
- Controllers: relay, channel, token, log, pricing, redemption, midjourney, task, model, playground, topup, billing, etc.
- Models: channel, token, log, pricing, redemption, ability, task, midjourney, etc.
- Services: quota, token_counter, http_client, convert, midjourney, task, etc.
- DTOs: All LLM-related request/response types
- Frontend Pages: Channel, Token, TopUp, Log, Midjourney, Model, Playground, Pricing, etc.

**Files Modified:**
- `main.go` - Removed LLM initialization
- `router/main.go` - Removed relay/video routers
- `router/api-router.go` - Auth-only endpoints
- `model/main.go` - Auth-only migrations
- `model/option.go` - Auth-only settings
- `middleware/auth.go` - Removed TokenAuth for LLM

**Known Issues:**
- Controller misc.go and option.go still reference removed packages
- Need to create stub functions or further cleanup
