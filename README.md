<div align="center">

![auth-gateway](/web/public/logo.png)

# Auth Gateway

🔐 **Production-Ready Multi-Provider Authentication System**

<p align="center">
  <a href="https://raw.githubusercontent.com/karlorz/auth-gateway/main/LICENSE">
    <img src="https://img.shields.io/github/license/karlorz/auth-gateway?color=brightgreen" alt="license">
  </a>
  <a href="https://github.com/karlorz/auth-gateway/releases/latest">
    <img src="https://img.shields.io/github/v/release/karlorz/auth-gateway?color=brightgreen&include_prereleases" alt="release">
  </a>
  <a href="https://hub.docker.com/r/karl8080/auth-gateway">
    <img src="https://img.shields.io/badge/docker-dockerHub-blue" alt="docker">
  </a>
</p>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-features">Features</a> •
  <a href="#-deployment">Deployment</a> •
  <a href="#-documentation">Documentation</a>
</p>

</div>

## 📝 Project Description

**Auth Gateway** provides a comprehensive authentication solution that can be easily integrated into any fullstack application:

- **Multiple Authentication Methods** - Password, OAuth, Passkey/WebAuthn, 2FA
- **Production Ready** - Battle-tested CI/CD, Docker support, multi-database
- **Easy to Customize** - Trim non-auth components for your specific needs

---

## ✨ Features

### 🔐 Authentication Methods

| Method | Description |
|--------|-------------|
| 🔑 Password | Username/email + password with bcrypt hashing |
| 🐙 GitHub OAuth | OAuth 2.0 integration |
| 💬 Discord OAuth | OAuth 2.0 integration |
| 🔗 OIDC | OpenID Connect for enterprise SSO |
| 📱 WeChat | QR code-based authentication |
| ✈️ Telegram | Bot widget authentication |
| 🐧 LinuxDO | Community OAuth with trust levels |
| 🔒 Passkey/WebAuthn | Passwordless biometric/hardware key auth |
| 📲 2FA (TOTP) | Time-based one-time passwords with backup codes |

### 🛡️ Security Features

- ✅ Cloudflare Turnstile bot protection
- ✅ Rate limiting on auth endpoints
- ✅ Session management with secure cookies
- ✅ CSRF protection for OAuth flows
- ✅ Account lockout after failed attempts
- ✅ Signature counter for passkey clone detection

### 🔧 Technical

- ⚡ Go 1.25+ with Gin framework
- ⚛️ React 18 with Vite, Semi UI, Tailwind CSS
- 🗄️ Multi-database: SQLite, MySQL, PostgreSQL
- 🚀 Redis caching (optional)
- 🐳 Docker & Docker Compose ready
- 🔄 GitHub Actions CI/CD

---

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the project
git clone https://github.com/karlorz/auth-gateway.git
cd auth-gateway

# Edit configuration
cp .env.example .env
nano .env

# Start services
docker-compose up -d
```

### Using Docker

```bash
# SQLite (default)
docker run --name auth-gateway -d --restart always \
  -p 3000:3000 \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  karl8080/auth-gateway:latest

# MySQL
docker run --name auth-gateway -d --restart always \
  -p 3000:3000 \
  -e SQL_DSN="root:123456@tcp(localhost:3306)/auth_gateway" \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  karl8080/auth-gateway:latest
```

### Development

```bash
# Backend
go run main.go

# Frontend
cd web
bun install
bun run dev
```

---

## 📚 Documentation

### 📖 [Authentication System Guide](docs/AUTH.md)

Comprehensive documentation covering:

- Architecture overview
- All authentication methods in detail
- Database schema
- API endpoints reference
- Backend & Frontend implementation
- Package extraction guide
- Integration examples

### Quick Links

| Topic | Link |
|-------|------|
| Auth Architecture | [docs/AUTH.md#architecture](docs/AUTH.md#architecture) |
| API Endpoints | [docs/AUTH.md#api-endpoints](docs/AUTH.md#api-endpoints) |
| Database Schema | [docs/AUTH.md#database-schema](docs/AUTH.md#database-schema) |
| Configuration | [docs/AUTH.md#configuration](docs/AUTH.md#configuration) |
| Integration Examples | [docs/AUTH.md#integration-examples](docs/AUTH.md#integration-examples) |

---

## ⚙️ Configuration

### Environment Variables

```bash
# Database
SQL_DSN="root:password@tcp(localhost:3306)/auth_gateway"

# Session (REQUIRED for multi-node)
SESSION_SECRET=your-session-secret
CRYPTO_SECRET=your-crypto-secret

# Redis (optional, recommended for production)
REDIS_CONN_STRING=redis://localhost:6379/0

# OAuth Providers
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
TELEGRAM_BOT_TOKEN=xxx
TELEGRAM_BOT_NAME=xxx

# Turnstile (Cloudflare)
TURNSTILE_SITE_KEY=xxx
TURNSTILE_SECRET_KEY=xxx
```

See `.env.example` for complete configuration options.

---

## 🗂️ Project Structure

```
auth-gateway/
├── controller/           # HTTP handlers
│   ├── user.go          # Password auth (login, register)
│   ├── passkey.go       # WebAuthn endpoints
│   ├── twofa.go         # 2FA endpoints
│   ├── github.go        # GitHub OAuth
│   ├── discord.go       # Discord OAuth
│   ├── oidc.go          # OpenID Connect
│   ├── telegram.go      # Telegram OAuth
│   ├── wechat.go        # WeChat OAuth
│   └── linuxdo.go       # LinuxDO OAuth
├── middleware/           # Auth middleware
│   ├── auth.go          # UserAuth, AdminAuth, TokenAuth
│   └── turnstile-check.go
├── model/               # Database models
│   ├── user.go          # User model
│   ├── passkey.go       # PasskeyCredential
│   └── twofa.go         # TwoFA, TwoFABackupCode
├── service/
│   └── passkey/         # WebAuthn service
├── setting/             # Configuration
│   └── system_setting/  # OAuth & Passkey settings
├── router/              # Route definitions
├── web/                 # React frontend
│   └── src/
│       └── components/auth/  # Auth UI components
├── docs/
│   └── AUTH.md          # Auth system documentation
├── .github/workflows/   # CI/CD
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

---

## 🗺️ Roadmap

- [ ] Trim non-auth components for lighter package
- [ ] Create standalone Go auth SDK
- [ ] Add more OAuth providers (Google, Apple, Microsoft)
- [ ] WebAuthn attestation format support
- [ ] Magic link authentication
- [ ] OAuth token refresh handling

---

## 🤝 Contributing

Contributions are welcome!

- 🐛 Report bugs
- 💡 Suggest features
- 📝 Improve documentation
- 🔧 Submit PRs

---

## 📄 License

MIT - See [LICENSE](LICENSE) for details.

---

## 🙏 Credits

- WebAuthn support via [go-webauthn](https://github.com/go-webauthn/webauthn)
- TOTP support via [pquerna/otp](https://github.com/pquerna/otp)

---

<div align="center">

### 💖 Thanks for using Auth Gateway!

If this project helps you, please give us a ⭐️ Star!

<sub>Built with ❤️</sub>

</div>
