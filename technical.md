Pangolin CMS - Technical Documentation
Version: 1.0.0
Status: Draft
Author: Maidul
License: Proprietary (Phase 1), AGPL v3 (Phase 2 - Open Source)
1. Project Overview
Pangolin CMS is a high-performance, multi-tenant blogging platform built with Go (Golang). It is designed to operate as a SaaS (Pangolin Cloud) initially, with the capability to be packaged as a Self-Hosted Single Binary later.
1.1 Core Value Proposition
Performance: Compiled Go binary, handling 10x more concurrency than PHP/Python.
Multi-Tenancy: Single instance hosts thousands of blogs (user.pangolin.com or custom.com).
Hybrid Deployment: Works as a managed service OR a downloadable static binary.
Security: Armored against common web exploits (SQLi, XSS) via Go's type safety and strict middleware.
1.2 Business Strategy
Phase 1 (SaaS): Host pangolin.com. Users get free subdomains. Revenue via Premium Plans (Custom Domains, Analytics).
Phase 2 (Open Source): Release core code on GitHub. Users can self-host. Revenue via Managed Hosting & Enterprise Support.
2. System Architecture
2.1 High-Level Diagram
mermaid




graph TD
    User[User Browser] -->|HTTPS| LB[Load Balancer / Traefik]
    LB -->|Host Header| App[Pangolin Go Server]
    
    subgraph "Pangolin Server (Go Fiber)"
        Middleware[Tenant Resolver Middleware]
        Router[API & Web Routes]
        SSL[CertMagic SSL Manager]
    end
    
    App --> Middleware
    Middleware --> Router
    Router --> SSL
    
    subgraph "Data Layer"
        DB[(PostgreSQL / SQLite)]
        Cache[(Redis / Memory)]
        Storage[S3 / Local Disk]
    end
    
    Router --> DB
    Router --> Cache
    Router --> Storage








2.2 Deployment Modes
Feature
SaaS Mode (Cloud)
Self-Hosted Mode (Binary)
Database
PostgreSQL (Clustered)
SQLite (Embedded)
Storage
AWS S3 / Cloudflare R2
Local Filesystem
SSL
Automatic (CertMagic)
Manual or Off
Tenancy
Multi-Tenant (Subdomains)
Single-Tenant
Config
Env Vars + DB Settings
.env file
3. Technology Stack
Component
Technology
Justification
Language
Go (1.21+)
Performance, Concurrency, Single Binary
Web Framework
Fiber v2
Express-like syntax, fastest Go framework
ORM
GORM
Familiar to Laravel/Eloquent devs, robust
Database
PostgreSQL (SaaS) / SQLite (Self)
Relational integrity vs. Zero-config
Frontend
HTMX + Tailwind CSS
No build step, embedded in binary, fast
Templates
Go html/template
Secure SSR, prevents XSS
SSL/TLS
CertMagic
Automatic HTTPS management for custom domains
Config
Viper
Robust environment variable management
Logging
Zap
High-performance structured logging
Storage
S3 Protocol
Uniform interface for Local & Cloud storage
4. Database Schema (Multi-Tenant)
Critical Rule: Every content table MUST have a site_id to ensure data isolation.
4.1 Users Table (Global)
sql
1234567
4.2 Sites Table (Tenants)
sql
1234567891011
4.3 Posts Table (Content)
sql
1234567891011
4.4 Domains Table (SSL Management)
sql
12345678
5. Core Logic & Middleware
5.1 Tenant Resolution Middleware
This is the heart of the multi-tenant system. It runs on every request.
Extract Host: Get Host header (e.g., blog.pangolin.com or mysite.com).
Check Cache: Look for site_id in Redis/Memory cache.
Query DB: If not cached, query sites table where slug = host OR custom_domain = host.
Set Context: Store site object in c.Locals("site").
Scope Queries: All subsequent DB queries must use c.Locals("site_id").
5.2 Storage Interface
To support both SaaS (S3) and Self-Hosted (Local), use an interface:
go
12345
SaaS Implementation: Uploads to S3 Bucket pangolin-user-uploads.
Self-Hosted Implementation: Saves to ./uploads/{site_id}/.
5.3 SSL Automation (CertMagic)
Wildcard Cert: Used for *.pangolin.com (Managed by Platform).
Dynamic Cert: Used for Custom Domains.
User adds domain in Dashboard.
System generates DNS TXT challenge.
User adds TXT record.
System verifies & issues Cert via Let's Encrypt.
Cert stored in /data/certs (Self-Hosted) or S3 (SaaS).
6. API Design (Internal & Public)
6.1 Authentication
POST /api/auth/login (Returns JWT)
POST /api/auth/register (Creates User + Default Site)
POST /api/auth/logout
6.2 Sites (Tenants)
GET /api/sites/me (Get current site context)
PUT /api/sites/settings (Update theme, SEO)
POST /api/sites/domains (Add custom domain)
6.3 Content
GET /api/posts (List posts for current site)
POST /api/posts (Create draft)
PUT /api/posts/:id (Update/Publish)
DELETE /api/posts/:id
6.4 Public (SSR)
GET / (Home Page)
GET /:slug (Post Page)
GET /sitemap.xml
GET /rss.xml
7. Security & Isolation
7.1 Data Isolation
Rule: No query ever runs without WHERE site_id = ?.
Implementation: Create a GORM Scope func ScopeSite(db *gorm.DB) *gorm.DB that automatically injects site_id.
7.2 Rate Limiting
Per IP: Prevent DDoS.
Per Site: Prevent one free blog from consuming all server resources.
Tool: golang.org/x/time/rate or Redis-based limiter.
7.3 Input Sanitization
HTML: Use bluemonday to sanitize user content (prevent XSS in posts).
SQL: GORM prevents SQLi if used correctly (no raw strings).
7.4 Licensing (Future)
Phase 1: Closed Source (Proprietary).
Phase 2: AGPL v3. Ensures if someone modifies Pangolin and runs it as a service, they must release their changes. Protects SaaS revenue.
8. Configuration (Environment Variables)
bash
12345678910111213141516171819202122
# Server
PORT=8080
MODE=saas # or 'self-hosted'

# Database
DB_DRIVER=postgres # or 'sqlite'
DB_HOST=localhost
DB_USER=pangolin
DB_PASS=secret
DB_NAME=pangolin_db

9. Development Guidelines (Vibe Coding)
9.1 Folder Structure
text
12345
/cmd            # Entry points
/internal       # Private app code (Not for import)
/pkg            # Public libraries (Safe for open source)
/web            # Frontend assets (HTMX, Tailwind)
/scripts        # DB Migrations, Build scripts
9.2 AI Prompting Rules
When using AI to generate code for Pangolin:
Context: Always state "This is for a Multi-Tenant Go Fiber App".
Database: Always require site_id in queries.
Error Handling: No panic(). Return error and handle gracefully.
Security: Sanitize all user inputs.
Style: Follow golangci-lint standards.
9.3 Testing Strategy
Unit Tests: For Services and Utils (_test.go).
Integration Tests: For API Endpoints (using httptest).
Tenant Tests: Ensure Site A cannot access Site B's data.
10. Deployment Pipeline
10.1 SaaS (Cloud)
CI/CD: GitHub Actions builds Docker image.
Registry: Push to Docker Hub / GHCR.
Deploy: Kubernetes or Docker Swarm on VPS.
DB: Managed PostgreSQL (AWS RDS / DigitalOcean).
Storage: S3 / R2.
10.2 Self-Hosted (Binary)
Build: go build -o pangolin
Embed: go:embed web/* (Includes frontend in binary).
Release: Upload binary to GitHub Releases.
Run: User downloads & executes ./pangolin.
11. Future Roadmap (Post-Launch)
Plugin System: Webhooks first, Go Plugins later.
Theme Store: Marketplace for paid themes.
Analytics: Built-in privacy-friendly analytics (replace Google Analytics).
Membership: Paid subscriptions for readers (Stripe integration).
12. Glossary
Tenant: A single blog/site within the platform.
Host: The domain name used to access the site.
Slug: The URL-friendly identifier for a subdomain (e.g., slug.pangolin.com).
CertMagic: The Go library used for automatic HTTPS.
HTMX: The frontend library used for dynamic UI without JavaScript frameworks.