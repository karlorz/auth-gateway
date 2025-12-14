# Auth Gateway - Authentication System Guide

This document provides comprehensive documentation for the Auth Gateway authentication system, a reusable package that can be easily imported into other fullstack applications.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Authentication Methods](#authentication-methods)
4. [Database Schema](#database-schema)
5. [API Endpoints](#api-endpoints)
6. [Backend Implementation](#backend-implementation)
7. [Frontend Implementation](#frontend-implementation)
8. [Configuration](#configuration)
9. [Package Extraction Guide](#package-extraction-guide)
10. [Integration Examples](#integration-examples)

---

## Overview

The Auth Gateway authentication system provides a comprehensive, modular authentication solution supporting:

- **Password Authentication** - Traditional username/email + password login
- **OAuth 2.0 Providers** - GitHub, Discord, OIDC, WeChat, Telegram, LinuxDO
- **Passkey/WebAuthn** - Passwordless authentication using biometrics/hardware keys
- **Two-Factor Authentication (2FA)** - TOTP-based with backup codes
- **Bot Protection** - Cloudflare Turnstile integration

### Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.25+ with Gin framework |
| Frontend | React 18 with Semi UI |
| Database | GORM (SQLite/MySQL/PostgreSQL) |
| Session | gin-contrib/sessions |
| WebAuthn | go-webauthn/webauthn |
| TOTP | pquerna/otp |

---

## Architecture

### High-Level Flow

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
│                             SERVICE LAYER                                    │
│  ┌─────────────────────────────┐ ┌──────────────────────────────────────────┐│
│  │ service/passkey/            │ │ common/totp.go, common/password.go       ││
│  │ - service.go (WebAuthn)     │ │ - Password hashing (bcrypt)              ││
│  │ - user.go (WebAuthnUser)    │ │ - TOTP generation/validation             ││
│  │ - session.go (SessionData)  │ │ - Backup code generation                 ││
│  └─────────────────────────────┘ └──────────────────────────────────────────┘│
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

---

## Authentication Methods

### 1. Password Authentication

**Flow:**
1. User submits username/email + password
2. Backend validates credentials against bcrypt hash
3. If 2FA enabled → returns `require_2fa: true`, stores pending session
4. If 2FA not enabled → creates session, returns user data

**Key Files:**
- `controller/user.go:30-95` - Login handler
- `controller/user.go:145-273` - Register handler
- `model/user.go:501-517` - ValidateAndFill

**Security Features:**
- Password hashing with bcrypt
- Rate limiting on login endpoint
- Turnstile bot protection (optional)

### 2. OAuth 2.0 Providers

All OAuth providers follow the same pattern:

```
1. Frontend → GET /api/oauth/state → Generate CSRF state token
2. Frontend → Redirect to provider's authorization URL
3. Provider → Callback to /api/oauth/{provider}?code=XXX&state=YYY
4. Backend → Validate state, exchange code for token
5. Backend → Fetch user info from provider
6. Backend → Create/login user, setup session
```

#### GitHub OAuth
- **Controller:** `controller/github.go`
- **Settings:** `common.GitHubClientId`, `common.GitHubClientSecret`
- **Endpoints:** Token: `https://github.com/login/oauth/access_token`, User: `https://api.github.com/user`

#### Discord OAuth
- **Controller:** `controller/discord.go`
- **Settings:** `setting/system_setting/discord.go` → `DiscordSettings`
- **Endpoints:** Token: `https://discord.com/api/v10/oauth2/token`, User: `https://discord.com/api/v10/users/@me`

#### OIDC (OpenID Connect)
- **Controller:** `controller/oidc.go`
- **Settings:** `setting/system_setting/oidc.go` → `OIDCSettings`
- **Configurable endpoints:** `AuthorizationEndpoint`, `TokenEndpoint`, `UserInfoEndpoint`

#### Telegram
- **Controller:** `controller/telegram.go`
- **Settings:** `common.TelegramBotToken`, `common.TelegramBotName`
- **Security:** HMAC-SHA256 signature verification

#### WeChat
- **Controller:** `controller/wechat.go`
- **Settings:** `common.WeChatServerAddress`, `common.WeChatServerToken`
- **Flow:** QR code scan → Verification code → Backend validation

#### LinuxDO
- **Controller:** `controller/linuxdo.go`
- **Settings:** `common.LinuxDOClientId`, `common.LinuxDOClientSecret`
- **Feature:** Trust level validation (`LinuxDOMinimumTrustLevel`)

### 3. Passkey/WebAuthn

**Registration Flow:**
```
1. POST /api/user/passkey/register/begin → Get creation options
2. Browser prompts for biometric/PIN
3. POST /api/user/passkey/register/finish → Submit attestation
4. Server validates & stores credential
```

**Login Flow (Discoverable):**
```
1. POST /api/user/passkey/login/begin → Get assertion options
2. Browser performs resident key lookup
3. POST /api/user/passkey/login/finish → Submit assertion
4. Server validates signature, creates session
```

**Key Files:**
- `controller/passkey.go` - All passkey endpoints
- `service/passkey/service.go` - WebAuthn configuration
- `service/passkey/user.go` - WebAuthnUser interface implementation
- `service/passkey/session.go` - Session data storage
- `model/passkey.go` - PasskeyCredential model
- `setting/system_setting/passkey.go` - PasskeySettings

**Configuration Options:**
```go
type PasskeySettings struct {
    Enabled              bool   // Enable/disable passkey auth
    RPDisplayName        string // Relying Party display name
    RPID                 string // Domain (e.g., "example.com")
    Origins              string // Comma-separated allowed origins
    AllowInsecureOrigin  bool   // Allow HTTP (for localhost)
    UserVerification     string // "preferred", "required", "discouraged"
    AttachmentPreference string // Platform authenticator preference
}
```

### 4. Two-Factor Authentication (2FA/TOTP)

**Setup Flow:**
```
1. POST /api/user/2fa/setup → Generate secret + QR code + backup codes
2. User scans QR with authenticator app
3. POST /api/user/2fa/enable → Submit verification code to enable
```

**Verification Flow (during login):**
```
1. Password login returns require_2fa: true
2. POST /api/user/login/2fa → Submit 6-digit TOTP or backup code
3. Server validates, completes login
```

**Key Files:**
- `controller/twofa.go` - All 2FA endpoints
- `model/twofa.go` - TwoFA and TwoFABackupCode models
- `common/totp.go` - TOTP generation/validation utilities

**Security Features:**
- Account lockout after 5 failed attempts (configurable)
- Lockout duration: 15 minutes (configurable)
- 10 single-use backup codes
- Backup codes hashed with bcrypt

---

## Database Schema

### User Table

```sql
CREATE TABLE users (
    id              INT PRIMARY KEY AUTO_INCREMENT,
    username        VARCHAR(20) UNIQUE NOT NULL,
    password        VARCHAR(255) NOT NULL,
    display_name    VARCHAR(20),
    role            INT DEFAULT 1,          -- 1=common, 10=admin, 100=root
    status          INT DEFAULT 1,          -- 1=enabled, 2=disabled
    email           VARCHAR(50) INDEX,
    github_id       VARCHAR(255) INDEX,
    discord_id      VARCHAR(255) INDEX,
    oidc_id         VARCHAR(255) INDEX,
    wechat_id       VARCHAR(255) INDEX,
    telegram_id     VARCHAR(255) INDEX,
    linux_do_id     VARCHAR(255) INDEX,
    access_token    CHAR(32) UNIQUE INDEX,  -- API access token
    quota           INT DEFAULT 0,
    used_quota      INT DEFAULT 0,
    request_count   INT DEFAULT 0,
    `group`         VARCHAR(64) DEFAULT 'default',
    aff_code        VARCHAR(32) UNIQUE,
    aff_count       INT DEFAULT 0,
    aff_quota       INT DEFAULT 0,
    aff_history     INT DEFAULT 0,
    inviter_id      INT INDEX,
    setting         TEXT,                   -- JSON user settings
    remark          VARCHAR(255),
    stripe_customer VARCHAR(64) INDEX,
    deleted_at      TIMESTAMP INDEX
);
```

### PasskeyCredential Table

```sql
CREATE TABLE passkey_credentials (
    id               INT PRIMARY KEY AUTO_INCREMENT,
    user_id          INT UNIQUE NOT NULL,
    credential_id    VARCHAR(512) UNIQUE NOT NULL,  -- base64 encoded
    public_key       TEXT NOT NULL,                  -- base64 encoded
    attestation_type VARCHAR(255),
    aaguid           VARCHAR(512),                   -- base64 encoded
    sign_count       INT DEFAULT 0,
    clone_warning    BOOLEAN,
    user_present     BOOLEAN,
    user_verified    BOOLEAN,
    backup_eligible  BOOLEAN,
    backup_state     BOOLEAN,
    transports       TEXT,                           -- JSON array
    attachment       VARCHAR(32),
    last_used_at     TIMESTAMP,
    created_at       TIMESTAMP,
    updated_at       TIMESTAMP,
    deleted_at       TIMESTAMP INDEX
);
```

### TwoFA Table

```sql
CREATE TABLE two_fas (
    id              INT PRIMARY KEY AUTO_INCREMENT,
    user_id         INT UNIQUE NOT NULL,
    secret          VARCHAR(255) NOT NULL,  -- TOTP secret
    is_enabled      BOOLEAN,
    failed_attempts INT DEFAULT 0,
    locked_until    TIMESTAMP,
    last_used_at    TIMESTAMP,
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP INDEX
);
```

### TwoFABackupCode Table

```sql
CREATE TABLE two_fa_backup_codes (
    id         INT PRIMARY KEY AUTO_INCREMENT,
    user_id    INT NOT NULL INDEX,
    code_hash  VARCHAR(255) NOT NULL,  -- bcrypt hash
    is_used    BOOLEAN,
    used_at    TIMESTAMP,
    created_at TIMESTAMP,
    deleted_at TIMESTAMP INDEX
);
```

---

## API Endpoints

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
| GET | `/api/verification` | Send email verification |
| GET | `/api/reset_password` | Send password reset email |
| POST | `/api/user/reset` | Reset password |

### Authenticated Endpoints (User Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/user/logout` | Logout |
| GET | `/api/user/self` | Get current user info |
| PUT | `/api/user/self` | Update current user |
| DELETE | `/api/user/self` | Delete account |
| GET | `/api/user/token` | Generate access token |
| GET | `/api/user/passkey` | Get passkey status |
| POST | `/api/user/passkey/register/begin` | Begin passkey registration |
| POST | `/api/user/passkey/register/finish` | Complete passkey registration |
| POST | `/api/user/passkey/verify/begin` | Begin passkey verification |
| POST | `/api/user/passkey/verify/finish` | Complete passkey verification |
| DELETE | `/api/user/passkey` | Remove passkey |
| GET | `/api/user/2fa/status` | Get 2FA status |
| POST | `/api/user/2fa/setup` | Initialize 2FA setup |
| POST | `/api/user/2fa/enable` | Enable 2FA |
| POST | `/api/user/2fa/disable` | Disable 2FA |
| POST | `/api/user/2fa/backup_codes` | Regenerate backup codes |
| GET | `/api/oauth/email/bind` | Bind email to account |
| GET | `/api/oauth/wechat/bind` | Bind WeChat to account |
| GET | `/api/oauth/telegram/bind` | Bind Telegram to account |

### Admin Endpoints (Admin Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/user/` | List all users |
| GET | `/api/user/search` | Search users |
| GET | `/api/user/:id` | Get user by ID |
| POST | `/api/user/` | Create user |
| PUT | `/api/user/` | Update user |
| DELETE | `/api/user/:id` | Delete user |
| POST | `/api/user/manage` | Manage user (enable/disable/promote/demote) |
| DELETE | `/api/user/:id/reset_passkey` | Admin reset user's passkey |
| GET | `/api/user/2fa/stats` | Get 2FA statistics |
| DELETE | `/api/user/:id/2fa` | Admin disable user's 2FA |

---

## Backend Implementation

### Core Files to Extract

```
auth/
├── controller/
│   ├── user.go          # Login, Register, Logout, User management
│   ├── passkey.go       # WebAuthn endpoints
│   ├── twofa.go         # 2FA endpoints
│   ├── github.go        # GitHub OAuth
│   ├── discord.go       # Discord OAuth
│   ├── oidc.go          # OpenID Connect
│   ├── telegram.go      # Telegram OAuth
│   ├── wechat.go        # WeChat OAuth
│   └── linuxdo.go       # LinuxDO OAuth
├── middleware/
│   ├── auth.go          # UserAuth, AdminAuth, RootAuth, TokenAuth
│   └── turnstile-check.go # Cloudflare Turnstile
├── model/
│   ├── user.go          # User model
│   ├── passkey.go       # PasskeyCredential model
│   └── twofa.go         # TwoFA, TwoFABackupCode models
├── service/
│   └── passkey/
│       ├── service.go   # WebAuthn configuration
│       ├── user.go      # WebAuthnUser interface
│       └── session.go   # Session data helpers
├── setting/
│   └── system_setting/
│       ├── passkey.go   # PasskeySettings
│       ├── discord.go   # DiscordSettings
│       └── oidc.go      # OIDCSettings
└── common/
    ├── password.go      # Password hashing
    ├── totp.go          # TOTP utilities
    └── utils.go         # Random string generation, etc.
```

### Key Dependencies

```go
// go.mod
require (
    github.com/gin-gonic/gin v1.9+
    github.com/gin-contrib/sessions v0.0.5
    github.com/go-webauthn/webauthn v0.10+
    github.com/pquerna/otp v1.4+
    golang.org/x/crypto v0.18+ // bcrypt
    gorm.io/gorm v1.25+
)
```

### Session Management

The auth system uses `gin-contrib/sessions` for session management:

```go
import "github.com/gin-contrib/sessions"

// Setup session (after successful login)
func setupLogin(user *model.User, c *gin.Context) {
    session := sessions.Default(c)
    session.Set("id", user.Id)
    session.Set("username", user.Username)
    session.Set("role", user.Role)
    session.Set("status", user.Status)
    session.Set("group", user.Group)
    session.Save()
    // ... return user data
}

// Read session (in middleware)
func authHelper(c *gin.Context, minRole int) {
    session := sessions.Default(c)
    username := session.Get("username")
    role := session.Get("role")
    id := session.Get("id")
    status := session.Get("status")
    // ... validate and set context
}
```

### Password Hashing

```go
// common/password.go
import "golang.org/x/crypto/bcrypt"

func Password2Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func ValidatePasswordAndHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### TOTP Implementation

```go
// common/totp.go
import "github.com/pquerna/otp/totp"

func GenerateTOTPSecret(username string) (*otp.Key, error) {
    return totp.Generate(totp.GenerateOpts{
        Issuer:      "AuthGateway",
        AccountName: username,
    })
}

func ValidateTOTPCode(secret, code string) bool {
    return totp.Validate(code, secret)
}

func GenerateQRCodeData(secret, username string) string {
    return fmt.Sprintf("otpauth://totp/AuthGateway:%s?secret=%s&issuer=AuthGateway", username, secret)
}
```

---

## Frontend Implementation

### Key Components

```
web/src/components/auth/
├── LoginForm.jsx        # Main login page with OAuth buttons
├── RegisterForm.jsx     # Registration form
├── TwoFAVerification.jsx # 2FA code input modal
├── PasswordResetForm.jsx
└── OAuth2Callback.jsx   # OAuth callback handler
```

### Helper Functions

```javascript
// helpers/index.js

// OAuth helpers
export async function onGitHubOAuthClicked(clientId, options = {}) {
    const stateRes = await API.get('/api/oauth/state');
    const state = stateRes.data.data;
    window.location.href = `https://github.com/login/oauth/authorize?client_id=${clientId}&state=${state}&scope=user:email`;
}

export async function onDiscordOAuthClicked(clientId, options = {}) {
    const stateRes = await API.get('/api/oauth/state');
    const state = stateRes.data.data;
    const redirectUri = encodeURIComponent(`${window.location.origin}/oauth/discord`);
    window.location.href = `https://discord.com/api/oauth2/authorize?client_id=${clientId}&redirect_uri=${redirectUri}&response_type=code&scope=identify&state=${state}`;
}

// Passkey helpers
export async function isPasskeySupported() {
    if (!window.PublicKeyCredential) return false;
    try {
        return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
    } catch {
        return false;
    }
}

export function prepareCredentialRequestOptions(options) {
    // Convert base64url to ArrayBuffer
    const prepared = { ...options };
    if (prepared.challenge) {
        prepared.challenge = base64urlToArrayBuffer(prepared.challenge);
    }
    if (prepared.allowCredentials) {
        prepared.allowCredentials = prepared.allowCredentials.map(cred => ({
            ...cred,
            id: base64urlToArrayBuffer(cred.id)
        }));
    }
    return prepared;
}

export function buildAssertionResult(assertion) {
    return {
        id: assertion.id,
        rawId: arrayBufferToBase64url(assertion.rawId),
        response: {
            authenticatorData: arrayBufferToBase64url(assertion.response.authenticatorData),
            clientDataJSON: arrayBufferToBase64url(assertion.response.clientDataJSON),
            signature: arrayBufferToBase64url(assertion.response.signature),
            userHandle: assertion.response.userHandle
                ? arrayBufferToBase64url(assertion.response.userHandle)
                : null
        },
        type: assertion.type
    };
}
```

### Status Object

The frontend reads authentication configuration from `localStorage.status`:

```javascript
const status = {
    // Password auth
    password_login_enabled: true,
    password_register_enabled: true,
    register_enabled: true,
    email_verification: false,

    // OAuth providers
    github_oauth: true,
    github_client_id: "xxx",
    discord_oauth: true,
    discord_client_id: "xxx",
    oidc_enabled: true,
    oidc_client_id: "xxx",
    oidc_authorization_endpoint: "https://...",
    wechat_login: true,
    wechat_qrcode: "https://...",
    telegram_oauth: true,
    telegram_bot_name: "xxx",
    linuxdo_oauth: true,
    linuxdo_client_id: "xxx",

    // Passkey
    passkey_login: true,

    // Turnstile
    turnstile_check: true,
    turnstile_site_key: "xxx",

    // Legal
    user_agreement_enabled: true,
    privacy_policy_enabled: true
};
```

---

## Configuration

### Environment Variables

```bash
# Session (REQUIRED for multi-node)
SESSION_SECRET=your-session-secret-key
CRYPTO_SECRET=your-crypto-secret-key

# GitHub OAuth
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx

# Discord OAuth (via settings API)
# Stored in database via system_setting

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

# Passkey (via settings API)
# Stored in database via system_setting

# OIDC (via settings API)
# Stored in database via system_setting

# Rate Limiting
CRITICAL_RATE_LIMIT=20  # per minute
```

### Settings API Structure

Settings are stored in the database and managed via `setting/system_setting/`:

```go
// Passkey
type PasskeySettings struct {
    Enabled              bool   `json:"enabled"`
    RPDisplayName        string `json:"rp_display_name"`
    RPID                 string `json:"rp_id"`
    Origins              string `json:"origins"`
    AllowInsecureOrigin  bool   `json:"allow_insecure_origin"`
    UserVerification     string `json:"user_verification"`
    AttachmentPreference string `json:"attachment_preference"`
}

// Discord
type DiscordSettings struct {
    Enabled      bool   `json:"enabled"`
    ClientId     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

// OIDC
type OIDCSettings struct {
    Enabled               bool   `json:"enabled"`
    ClientId              string `json:"client_id"`
    ClientSecret          string `json:"client_secret"`
    WellKnown             string `json:"well_known"`
    AuthorizationEndpoint string `json:"authorization_endpoint"`
    TokenEndpoint         string `json:"token_endpoint"`
    UserInfoEndpoint      string `json:"user_info_endpoint"`
}
```

---

## Package Extraction Guide

### Step 1: Create Package Structure

```bash
mkdir -p auth-pkg/{controller,middleware,model,service/passkey,setting,common}
```

### Step 2: Extract Core Files

```go
// auth-pkg/auth.go - Main package entry point
package auth

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Config struct {
    // Session
    SessionSecret string

    // Database
    DB *gorm.DB

    // OAuth Providers
    GitHubClientID     string
    GitHubClientSecret string
    // ... other providers

    // Passkey
    PasskeyEnabled     bool
    PasskeyRPID        string
    PasskeyRPDisplayName string
    PasskeyOrigins     []string

    // 2FA
    TwoFAEnabled       bool
    MaxFailAttempts    int
    LockoutDuration    int // seconds

    // Turnstile
    TurnstileEnabled   bool
    TurnstileSiteKey   string
    TurnstileSecretKey string
}

func RegisterRoutes(router *gin.Engine, config Config) {
    // Initialize database models
    config.DB.AutoMigrate(&User{}, &PasskeyCredential{}, &TwoFA{}, &TwoFABackupCode{})

    // Setup routes
    api := router.Group("/api")
    {
        // Public routes
        api.POST("/user/register", Register(config))
        api.POST("/user/login", Login(config))
        api.POST("/user/login/2fa", Verify2FALogin(config))
        api.POST("/user/passkey/login/begin", PasskeyLoginBegin(config))
        api.POST("/user/passkey/login/finish", PasskeyLoginFinish(config))
        api.GET("/oauth/state", GenerateOAuthCode(config))
        api.GET("/oauth/github", GitHubOAuth(config))
        // ... more routes

        // Authenticated routes
        self := api.Group("/user")
        self.Use(UserAuth(config))
        {
            self.GET("/self", GetSelf(config))
            self.PUT("/self", UpdateSelf(config))
            // ... more routes
        }
    }
}
```

### Step 3: Adapt Models for Standalone Use

```go
// auth-pkg/model/user.go
package model

type User struct {
    ID           int    `gorm:"primaryKey"`
    Username     string `gorm:"unique;not null"`
    Password     string `gorm:"not null"`
    DisplayName  string
    Email        string `gorm:"index"`
    Role         int    `gorm:"default:1"`
    Status       int    `gorm:"default:1"`
    // OAuth IDs
    GitHubID     string `gorm:"index"`
    DiscordID    string `gorm:"index"`
    OidcID       string `gorm:"index"`
    WeChatID     string `gorm:"index"`
    TelegramID   string `gorm:"index"`
    LinuxDOID    string `gorm:"index"`
    // Timestamps
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### Step 4: Create Configuration Interface

```go
// auth-pkg/config.go
package auth

type AuthProvider interface {
    // Password
    PasswordLoginEnabled() bool
    PasswordRegisterEnabled() bool

    // OAuth
    GitHubOAuthEnabled() bool
    GitHubClientID() string
    GitHubClientSecret() string
    // ... other providers

    // Passkey
    PasskeyEnabled() bool
    PasskeySettings() *PasskeySettings

    // 2FA
    TwoFASettings() *TwoFASettings

    // Turnstile
    TurnstileEnabled() bool
    TurnstileSiteKey() string
    TurnstileSecretKey() string
}
```

---

## Integration Examples

### Basic Integration (Gin)

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "your-project/auth"
)

func main() {
    // Setup database
    db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

    // Setup Gin
    r := gin.Default()

    // Setup sessions
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("session", store))

    // Register auth routes
    auth.RegisterRoutes(r, auth.Config{
        DB:            db,
        SessionSecret: "your-secret",

        // GitHub OAuth
        GitHubClientID:     "xxx",
        GitHubClientSecret: "xxx",

        // Passkey
        PasskeyEnabled:       true,
        PasskeyRPID:          "yourdomain.com",
        PasskeyRPDisplayName: "Your App",
        PasskeyOrigins:       []string{"https://yourdomain.com"},

        // 2FA
        TwoFAEnabled:    true,
        MaxFailAttempts: 5,
        LockoutDuration: 900, // 15 minutes
    })

    r.Run(":8080")
}
```

### React Frontend Integration

```javascript
// src/auth/AuthProvider.jsx
import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check if user is logged in
        const savedUser = localStorage.getItem('user');
        if (savedUser) {
            setUser(JSON.parse(savedUser));
        }
        setLoading(false);
    }, []);

    const login = async (username, password) => {
        const res = await axios.post('/api/user/login', { username, password });
        if (res.data.success) {
            if (res.data.data.require_2fa) {
                return { require2FA: true };
            }
            setUser(res.data.data);
            localStorage.setItem('user', JSON.stringify(res.data.data));
            return { success: true };
        }
        return { error: res.data.message };
    };

    const verify2FA = async (code) => {
        const res = await axios.post('/api/user/login/2fa', { code });
        if (res.data.success) {
            setUser(res.data.data);
            localStorage.setItem('user', JSON.stringify(res.data.data));
            return { success: true };
        }
        return { error: res.data.message };
    };

    const logout = async () => {
        await axios.get('/api/user/logout');
        setUser(null);
        localStorage.removeItem('user');
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, verify2FA, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}
```

---

## Security Considerations

### Best Practices Implemented

1. **Password Security**
   - bcrypt hashing with default cost
   - No password stored in plaintext
   - Password validation before hashing

2. **Session Security**
   - Session secret required
   - Session data not exposed to client
   - Session cleared on logout

3. **OAuth Security**
   - CSRF protection via state tokens
   - State validated on callback
   - Tokens not stored client-side

4. **Passkey Security**
   - User verification required by default
   - Signature counter for clone detection
   - HTTPS required (configurable for localhost)

5. **2FA Security**
   - TOTP with 30-second window
   - Account lockout after failed attempts
   - Backup codes hashed (not plaintext)
   - Single-use backup codes

6. **Rate Limiting**
   - Critical rate limit on auth endpoints
   - Global API rate limit
   - Lockout mechanism for 2FA

### Recommendations

1. Always use HTTPS in production
2. Set strong session secrets
3. Enable Turnstile for public-facing apps
4. Require 2FA for admin accounts
5. Monitor failed login attempts
6. Implement password complexity requirements
7. Regular security audits of OAuth configurations

---

## Troubleshooting

### Common Issues

**Passkey not working:**
- Check RPID matches your domain
- Ensure HTTPS is enabled (or AllowInsecureOrigin for localhost)
- Verify Origins include your full URL with protocol

**OAuth callback fails:**
- Verify callback URL in provider settings matches exactly
- Check state token is being saved/retrieved from session
- Ensure session middleware is configured before auth routes

**2FA always failing:**
- Check server time is synchronized (NTP)
- Verify secret is stored correctly (base32 encoded)
- Test with known TOTP generator

**Session not persisting:**
- Ensure session middleware is applied before routes
- Check cookie settings (domain, path, secure)
- Verify SESSION_SECRET is set for multi-node deployments

---

*This documentation is for Auth Gateway. For the latest updates, refer to the source code at https://github.com/karlorz/auth-gateway*
