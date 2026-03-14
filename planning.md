# 🦔 Pangolin CMS — Master Planning Document

> **Version:** 1.0.0 | **Status:** Active Planning | **Author:** Maidul  
> **License:** Proprietary (Phase 1 SaaS) → AGPL v3 (Phase 2 Open Source)

---

## 🎯 Project Vision

**Pangolin CMS** is a high-performance, multi-tenant blogging platform built in Go. It bridges the gap between the ease of WordPress and the raw speed of a compiled binary.

### Two Products, One Codebase

| Product | Description |
| :--- | :--- |
| **Pangolin Cloud** | Managed SaaS. Users get a free `user.pangolin.com` subdomain instantly. Revenue via Premium Plans (custom domains, analytics, etc.). |
| **Pangolin Core** | Open-source self-hosted binary. Single `./pangolin` executable. No LAMP/LEMP stack required. |

### Core Value Propositions
1. **Performance** — Compiled Go binary handles 10× more concurrency than PHP/Python.
2. **Multi-Tenancy** — One server hosts thousands of blogs, isolated by tenant.
3. **Hybrid Mode** — SaaS or self-hosted; users migrate freely between them.
4. **Single Binary** — Download and run. Frontend embedded via `go:embed`.
5. **Security** — Go type safety + strict middleware guards against SQLi, XSS, CSRF.

---

## 🛠 Tech Stack

| Layer | Technology | Justification |
| :--- | :--- | :--- |
| Language | Go 1.21+ | Performance, concurrency, single binary output |
| Web Framework | Fiber v2 | Express-like syntax, fastest Go HTTP framework |
| ORM | GORM | Familiar ActiveRecord-style API |
| Database (SaaS) | PostgreSQL | Clustered, multi-tenant at scale |
| Database (Self-Hosted) | SQLite | Zero-config embedded database |
| Frontend | HTMX + Tailwind CSS | No build step, embedded in binary |
| Templates | Go `html/template` | Secure SSR, XSS-safe by default |
| SSL/TLS | CertMagic | Automatic Let's Encrypt for custom domains |
| Config | Viper | Robust `.env` / env var management |
| Logging | Zap | Structured, high-performance logging |
| Storage | S3-compatible API | Uniform interface for local disk & cloud (S3/R2) |
| DNS Automation | Cloudflare API | Domain verification and TXT record management |
| Security | `bluemonday` | HTML sanitization for user-submitted content |
| Rate Limiting | `golang.org/x/time/rate` or Redis | Per-IP and per-site DDoS protection |

---

## 🏗 System Architecture

```
User Browser ──HTTPS──▶ Load Balancer (Traefik)
                               │
                               ▼
                   ┌─── Pangolin Go Server ───┐
                   │  Tenant Resolver          │
                   │  Middleware               │
                   │  API & Web Router         │
                   │  CertMagic SSL Manager    │
                   └───────────┬──────────────┘
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                 ▼
      PostgreSQL/SQLite     Redis/Memory      S3 / Local Disk
         (Data)              (Cache)           (Storage)
```

### Tenant Resolution (The Heart of Multi-Tenancy)
Every request goes through this middleware pipeline:
1. Read `Host` header (e.g., `blog.pangolin.com` or `mysite.com`).
2. Check Redis/Memory cache for `site_id`.
3. If not cached → query `sites` table (`slug = host` OR `custom_domain = host`).
4. Inject `site` object into `c.Locals("site")`.
5. **All subsequent DB queries MUST include `WHERE site_id = ?`.**

### Deployment Mode Comparison

| Feature | SaaS (Cloud) | Self-Hosted (Binary) |
| :--- | :--- | :--- |
| Database | PostgreSQL (clustered) | SQLite (embedded) |
| Storage | AWS S3 / Cloudflare R2 | Local filesystem (`./uploads/`) |
| SSL | Automatic via CertMagic | Manual or off |
| Multi-Tenancy | ✅ Yes | ❌ Single site |
| Config | Env vars + DB settings | `.env` file |

---

## 📐 Database Schema

> **Critical Rule:** Every content table MUST have `site_id` to guarantee data isolation.

```sql
-- Global Users Table
CREATE TABLE users (
  id         BIGSERIAL PRIMARY KEY,
  email      VARCHAR(255) UNIQUE NOT NULL,
  password   VARCHAR(255) NOT NULL, -- bcrypt
  role       VARCHAR(50) DEFAULT 'owner',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sites / Tenants Table
CREATE TABLE sites (
  id            BIGSERIAL PRIMARY KEY,
  user_id       BIGINT REFERENCES users(id),
  name          VARCHAR(255) NOT NULL,
  slug          VARCHAR(100) UNIQUE NOT NULL,   -- e.g., "myblog" → myblog.pangolin.com
  custom_domain VARCHAR(255) UNIQUE,
  plan          VARCHAR(50) DEFAULT 'free',
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- Posts Table
CREATE TABLE posts (
  id         BIGSERIAL PRIMARY KEY,
  site_id    BIGINT REFERENCES sites(id) NOT NULL,  -- tenant scoping
  title      VARCHAR(500) NOT NULL,
  slug       VARCHAR(500) NOT NULL,
  body       TEXT,
  status     VARCHAR(50) DEFAULT 'draft',           -- draft | published
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(site_id, slug)
);

-- Custom Domains SSL Management
CREATE TABLE domains (
  id          BIGSERIAL PRIMARY KEY,
  site_id     BIGINT REFERENCES sites(id) NOT NULL,
  domain      VARCHAR(255) UNIQUE NOT NULL,
  verified    BOOLEAN DEFAULT FALSE,
  cert_path   VARCHAR(500),
  created_at  TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 🗺 Development Roadmap

### Phase 1 — Core Foundation & Multi-Tenancy (Weeks 1–4)
**Goal:** A single running server that hosts multiple distinct blogs.

#### 1.1 Project Setup
- [ ] Initialize Go module: `go mod init pangolin`
- [ ] Define clean directory structure:
  ```
  /cmd          → Entry points (serve, export, import)
  /internal     → Private app code (handlers, services, middleware)
  /pkg          → Public reusable libraries
  /web          → HTMX templates + Tailwind assets
  /scripts      → DB migrations, build helpers
  ```
- [ ] Configure Viper for `.env` / env var management
- [ ] Set up structured logging with Zap

#### 1.2 Database & ORM
- [ ] Integrate GORM with PostgreSQL driver (SaaS) and SQLite driver (self-hosted)
- [ ] Implement schema migrations (`AutoMigrate` on startup)
- [ ] Create models: `User`, `Site`, `Post`, `Domain`, `Setting`, `Category`, `Tag`
- [ ] Implement GORM Scope: `ScopeSite(db)` — injects `WHERE site_id = ?` automatically

#### 1.3 Tenant Resolver Middleware
- [ ] Parse `Host` header on every request
- [ ] Query `sites` table (with Redis cache)
- [ ] Inject `site` into Fiber context (`c.Locals("site")`)
- [ ] Enforce that unauthenticated/unresolved hosts return `404`

#### 1.4 Authentication (Multi-Tenant Aware)
- [ ] Register → creates `User` + default `Site` (+ auto-subdomain)
- [ ] Login → returns JWT scoped to `user_id` + `site_id`
- [ ] Auth middleware for protected admin routes
- [ ] Password hashing with bcrypt

#### 1.5 Core API Endpoints
- [ ] `POST /api/auth/register` — Create user + site
- [ ] `POST /api/auth/login` — Return JWT
- [ ] `GET /api/sites/me` — Get current site context
- [ ] `PUT /api/sites/settings` — Update theme, SEO
- [ ] `GET /api/posts` — List posts (site-scoped)
- [ ] `POST /api/posts` — Create draft
- [ ] `PUT /api/posts/:id` — Update / publish
- [ ] `DELETE /api/posts/:id` — Delete post
- [ ] `GET /api/health` — Health check

**✅ Deliverable:** `site1.pangolin.com` and `site2.pangolin.com` serve different content from one server.

---

### Phase 2 — Domain & SSL Management (Weeks 5–8)
**Goal:** Users can connect custom domains with auto-SSL.

#### 2.1 Subdomain Logic
- [ ] Wildcard DNS `*.pangolin.com` configured at DNS level
- [ ] Logic to reserve/validate slugs (no squatting, no reserved words)
- [ ] Wildcard SSL cert for `*.pangolin.com`

#### 2.2 Custom Domain Logic
- [ ] UI for users to add `www.myblog.com`
- [ ] DNS TXT record verification flow (generate token → user sets record → system polls & verifies)
- [ ] CertMagic integration for automatic Let's Encrypt cert issuance
- [ ] Auto-renewal logic (renew certs 30 days before expiry)
- [ ] `domains` table updates (verified flag, cert path)

#### 2.3 Routing Middleware
- [ ] Priority: `Custom Domain > Subdomain > Default`
- [ ] HTTP → HTTPS redirect
- [ ] `404` for unregistered/unverified domains

**✅ Deliverable:** Users can connect `myblog.com` with HTTPS in under 5 minutes.

---

### Phase 3 — Admin Dashboard & Content (Weeks 9–12)
**Goal:** A polished UI for writing and managing content without touching code.

#### 3.1 Frontend Stack
- [ ] Embed Tailwind CSS + HTMX in binary via `//go:embed web/*`
- [ ] Alpine.js for minor client-side state (modals, dropdowns)
- [ ] Multi-site switcher (users who own multiple sites)

#### 3.2 Admin UI Pages
- [ ] Login / Register page
- [ ] Dashboard overview (posts count, recent activity)
- [ ] **Post Editor** — Markdown with live preview
- [ ] **Media Manager** — Image upload to S3/local; auto-convert to WebP
- [ ] **SEO Settings** per post — meta title, description, OpenGraph
- [ ] **Site Settings** — name, branding, permalink structure

#### 3.3 Public Theme Engine
- [ ] Default theme — minimalist, fast, SEO-100 Lighthouse ready
- [ ] Go `html/template` SSR rendering
- [ ] Dynamic routes: `/`, `/:slug`, `/category/:name`, `/tag/:name`
- [ ] `sitemap.xml` and `rss.xml` generation

**✅ Deliverable:** Full CMS — sign up, write a post, view it publicly on subdomain.

---

### Phase 4 — Security Hardening & Isolation (Weeks 13–15)
**Goal:** Tenants cannot access or affect each other in any way.

#### 4.1 Tenant Isolation Audit
- [ ] Audit ALL DB queries for missing `site_id` — use `ScopeSite()` everywhere
- [ ] Storage isolation: enforced `/uploads/{site_id}/` prefix
- [ ] Per-site rate limiting (one blog cannot DDoS the platform)
- [ ] File access validation — prevent path traversal attacks

#### 4.2 Platform Security
- [ ] Security headers: HSTS, CSP, X-Frame-Options, X-Content-Type
- [ ] CORS configuration for API routes
- [ ] `bluemonday` HTML sanitization for all user-submitted post content
- [ ] CSRF protection for admin form submissions
- [ ] 2FA for platform admins (TOTP)
- [ ] WAF middleware (block common attack patterns)

**✅ Deliverable:** Secure multi-tenant environment — passes internal security audit.

---

### Phase 5 — Distribution & Open Source Launch (Weeks 16–18)
**Goal:** v1.0.0 released publicly on GitHub.

#### 5.1 Self-Hosted Binary
- [ ] `go build -o pangolin` with `//go:embed web/*` for all frontend assets
- [ ] CLI commands: `pangolin serve`, `pangolin export`, `pangolin import`
- [ ] SQLite as default DB for self-hosted (zero configuration)
- [ ] `.env.example` template

#### 5.2 Static Export Engine
- [ ] `./pangolin export` command
- [ ] Crawl all public routes and render static `.html` files into `/dist`
- [ ] Copy CSS/JS/image assets to `/dist`
- [ ] Benefit: Deploy free to Netlify, S3, GitHub Pages

#### 5.3 WordPress Importer
- [ ] `./pangolin import wp.xml` command
- [ ] Parse WordPress WXR XML
- [ ] Map: WP Users → Pangolin Users
- [ ] Map: WP Posts → Pangolin Posts (preserve slugs, dates)
- [ ] Map: WP Images → Pangolin Media (download & reattach)

#### 5.4 CI/CD & Distribution
- [ ] GitHub Actions: build binaries for Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- [ ] Docker: minimal `Dockerfile` (distroless or scratch base)
- [ ] Docker Hub + GitHub Container Registry (GHCR) publish

#### 5.5 Documentation
- [ ] `README.md` — installation, usage, config reference
- [ ] Deployment guides: VPS, Docker, static hosting
- [ ] Theme development guide
- [ ] API documentation

**✅ Deliverable:** v1.0.0 GitHub Release with binary downloads and Docker image.

---

### Phase 6 — Pangolin Cloud Live (Week 19+)
**Goal:** Go live with the managed SaaS product.

#### 6.1 Infrastructure
- [ ] Load balancer: Traefik in front of Go instances
- [ ] PostgreSQL clustering with read replicas
- [ ] Object storage: AWS S3 or Cloudflare R2 for media
- [ ] Redis for caching tenant resolution and sessions

#### 6.2 Billing (Stripe Integration)
- [ ] Free tier limits (storage, bandwidth, custom domains: 0)
- [ ] Premium plan: custom domain, remove branding, analytics
- [ ] Stripe Billing Portal for subscription management

#### 6.3 Open Source Release
- [ ] Re-license core to AGPL v3 (forces SaaS forks to publish changes)
- [ ] GitHub repository public launch

#### 6.4 Marketing
- [ ] ProductHunt launch
- [ ] "Migrate from WordPress" landing page & campaign
- [ ] Dev blog / changelog

**✅ Deliverable:** Pangolin Cloud v1.0 live. Open-source core published.

---

## 🔐 Security Design Principles

1. **No query without `site_id`** — Use `ScopeSite(db *gorm.DB) *gorm.DB` GORM scope enforced globally.
2. **No `panic()`** — All errors returned and handled gracefully.
3. **Sanitize all inputs** — `bluemonday` for HTML; GORM parameterized queries for SQL.
4. **Context everywhere** — Always pass `context.Context` in DB and HTTP handlers.
5. **Minimal attack surface** — HTMX + Go templates avoid JS injection vectors.

---

## 🎯 Success KPIs

| Metric | Target |
| :--- | :--- |
| Onboarding Time | < 2 minutes from sign-up to live blog |
| TTFB | < 100ms with 1,000+ concurrent active sites |
| SSL Provisioning | 100% of custom domains get HTTPS automatically |
| Lighthouse Score | ≥ 95 (Performance, SEO, Accessibility) |
| Data Isolation | Zero cross-tenant data leaks in security audit |
| GitHub Stars | 100+ in first month post-launch |

---

## 📁 Directory Structure

```
pangolin/
├── cmd/
│   ├── serve/          # Main HTTP server entrypoint
│   ├── export/         # Static site export
│   └── import/         # WordPress importer
├── internal/
│   ├── handler/        # HTTP handlers (admin, public, api)
│   ├── middleware/      # Tenant resolver, auth, rate limiter
│   ├── model/          # GORM models (User, Site, Post, Domain…)
│   ├── repository/     # DB query layer (enforces ScopeSite)
│   ├── service/        # Business logic
│   └── storage/        # Storage interface (S3 + local implementations)
├── pkg/
│   └── certmagic/      # CertMagic wrapper
├── web/
│   ├── templates/      # Go html/template files
│   └── assets/         # Tailwind CSS, HTMX, Alpine.js
├── scripts/
│   └── migrate/        # DB migration scripts
├── .env.example
├── Dockerfile
└── README.md
```

---

## 🤖 AI Coding Guidelines (Vibe Coding Rules)

When generating code with AI for Pangolin:

| Concern | Rule |
| :--- | :--- |
| Context | Always state: *"This is for a Multi-Tenant Go Fiber App"* |
| Database | Always require `site_id` in all queries |
| Error Handling | No `panic()`. Return `error` and handle explicitly. |
| Security | Sanitize all user inputs before storage or rendering |
| Linting | Follow `golangci-lint` standards |
| Context | Pass `context.Context` in all DB and HTTP handler calls |

**Example AI Prompts:**
- *"Write a Fiber middleware that extracts the hostname, queries the `sites` table for a matching domain or subdomain, and attaches the Site object to the context. This is for a multi-tenant Go Fiber app."*
- *"Show me how to integrate CertMagic with Fiber to handle automatic HTTPS for dynamic custom domains."*
- *"How do I ensure GORM queries always include `WHERE site_id = ?` using a GORM Scope?"*

---

*Last Updated: 2026-02-25*
