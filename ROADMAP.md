# 🦔 Pangolin CMS - Project Roadmap

## 🚀 Project Vision
To build a **high-performance, secure, and distributable Blog CMS** using Go that bridges the gap between the ease of WordPress and the speed of Static Site Generators.

**Core Value Proposition:**
1.  **Single Binary:** Download & Run. No LAMP/LEMP stack required.
2.  **Hybrid Engine:** Run as a dynamic server OR export to static HTML.
3.  **WordPress Migration:** Seamless import from WP XML.
4.  **Modern Stack:** Go (Fiber), SQLite/Postgres, HTMX, Tailwind CSS.
5.  **Armored Security:** Built with Go's type safety and memory management to prevent common web exploits.

---

## 📅 Phase 1: Foundation & Core (Weeks 1-3)
**Goal:** A runnable binary that can serve a basic blog post.

### 1.1 Project Setup
- [ ] Initialize Go Module (`go mod init pangolin`).
- [ ] Setup Directory Structure (Clean Architecture).
- [ ] Configure `.env` management (Viper).
- [ ] Setup Logging (Zap or Logrus).

### 1.2 Database & ORM
- [ ] Integrate **GORM**.
- [ ] Implement **SQLite** (default) and **PostgreSQL** (configurable).
- [ ] Create Models: `User`, `Post`, `Category`, `Tag`, `Setting`.
- [ ] Implement Auto-Migration on startup.

### 1.3 Authentication
- [ ] JWT-based Authentication for Admin API.
- [ ] Session-based Auth for Admin Dashboard (HTMX friendly).
- [ ] Password Hashing (bcrypt).
- [ ] Middleware for Protected Routes.

### 1.4 Basic API
- [ ] CRUD APIs for Posts (`/api/posts`).
- [ ] CRUD APIs for Users (`/api/users`).
- [ ] Health Check Endpoint (`/api/health`).

**✅ Deliverable:** A compiled binary that runs on `localhost:8080` and allows API creation of posts.

---

## 📅 Phase 2: Admin Dashboard & Content (Weeks 4-6)
**Goal:** A usable UI for writing and managing content without touching code.

### 2.1 Frontend Stack
- [ ] Integrate **Tailwind CSS** (via CDN or embedded).
- [ ] Integrate **HTMX** for dynamic interactions (no React build step).
- [ ] Integrate **Alpine.js** for minor client-side state.

### 2.2 Admin UI
- [ ] Login/Register Page.
- [ ] Dashboard Overview (Stats: Total Posts, Views).
- [ ] Post Editor (Markdown support + Live Preview).
- [ ] Media Manager (Image upload to local disk/S3).
- [ ] Settings Page (Site Title, SEO, Permalinks).

### 2.3 Public Theme Engine
- [ ] Default Theme (Minimalist, Fast, SEO optimized).
- [ ] Go Template Rendering (`html/template`).
- [ ] Dynamic Routing (`/blog/:slug`, `/category/:name`).
- [ ] Sitemap.xml & RSS Feed Generation.

**✅ Deliverable:** A fully functional dynamic blog where users can log in, write posts, and view them publicly.

---

## 📅 Phase 3: The "Killer Features" (Weeks 7-9)
**Goal:** Differentiate from WordPress with Performance & Portability.

### 3.1 Static Export Engine
- [ ] Command: `./pangolin export`.
- [ ] Crawl all public routes.
- [ ] Generate static `.html` files in `/dist` folder.
- [ ] Copy assets (CSS/JS/Images) to `/dist`.
- [ ] **Benefit:** Host on Netlify/S3/GitHub Pages for free.

### 3.2 WordPress Importer
- [ ] Command: `./pangolin import wp.xml`.
- [ ] Parse WordPress WXR XML format.
- [ ] Map WP Users → Pangolin Users.
- [ ] Map WP Posts → Pangolin Posts (preserve dates/slugs).
- [ ] Map WP Images → Pangolin Media (download & reattach).
- [ ] **Benefit:** Zero friction migration for WP users.

### 3.3 Performance Optimization
- [ ] Implement Caching (Redis or In-Memory).
- [ ] Image Optimization (WebP conversion on upload).
- [ ] Minify HTML/CSS/JS on the fly.

**✅ Deliverable:** A hybrid CMS that can switch between Dynamic and Static modes and accepts WP migrations.

---

## 📅 Phase 4: Security & Distribution (Weeks 10-12)
**Goal:** Prepare for public release and real-world usage.

### 4.1 Security Hardening
- [ ] Rate Limiting (Middleware).
- [ ] CORS Configuration.
- [ ] SQL Injection Prevention (GORM handles most, but audit raw queries).
- [ ] XSS Protection (HTMX/Template escaping).
- [ ] Security Headers (HSTS, CSP, X-Frame-Options).

### 4.2 Distribution Pipeline
- [ ] **Go Embed:** Embed all frontend assets into the binary.
- [ ] **CI/CD:** GitHub Actions to build binaries for:
    - [ ] Linux (amd64, arm64)
    - [ ] macOS (amd64, arm64)
    - [ ] Windows (amd64)
- [ ] **Docker:** Create a minimal `Dockerfile` (scratch or distroless).

### 4.3 Documentation
- [ ] `README.md` (Installation, Usage, Config).
- [ ] Deployment Guides (VPS, Docker, Static Hosting).
- [ ] Theme Development Guide.

**✅ Deliverable:** v1.0.0 Release on GitHub.

---

## 📅 Phase 5: Growth & Ecosystem (Post-Launch)
**Goal:** Build a community and sustain the project.

- [ ] **Plugin System:** Explore Go Plugin support or Webhook-based extensions.
- [ ] **Theme Marketplace:** Allow users to submit themes.
- [ ] **Managed Hosting:** Offer a paid "One-Click Deploy" service.
- [ ] **Analytics:** Built-in privacy-focused analytics (no Google Analytics needed).

---

## 🛠️ Tech Stack Summary

| Component | Technology | Why? |
| :--- | :--- | :--- |
| **Language** | Go (Golang) | Performance, Single Binary, Concurrency |
| **Framework** | Fiber v2 | Express-like, Fast, Easy for PHP/Laravel devs |
| **Database** | SQLite / Postgres | SQLite for zero-config, Postgres for scale |
| **ORM** | GORM | Familiar to Laravel/Eloquent users |
| **Frontend** | HTMX + Tailwind | No build step, fast, fits in binary |
| **Templates** | Go `html/template` | Secure, Server-Side Rendering |
| **Config** | Viper | Robust environment management |
| **Deployment** | Docker / Binary | Universal compatibility |

---

## 🎯 Success Metrics (KPIs)
1.  **Performance:** Lighthouse Score > 95 (Performance, SEO, Accessibility).
2.  **Speed:** Time to First Byte (TTFB) < 50ms.
3.  **Adoption:** 100+ GitHub Stars in first month.
4.  **Migration:** Successfully import 100+ WordPress sites during beta.

---

## 📝 Notes for Vibe Coding
- **Prompt Strategy:** When using AI, specify "Use Fiber v2", "Use GORM", "Use HTMX".
- **Code Quality:** Enforce `golangci-lint` in CI pipeline.
- **Error Handling:** Go requires explicit error handling. Don't ignore `err`.
- **Context:** Always pass `context.Context` in database and HTTP handlers.
- **Naming:** Use "Pangolin" in CLI commands (e.g., `pangolin serve`, `pangolin export`).
