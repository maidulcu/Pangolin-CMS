# 🦔 Pangolin CMS - Project Roadmap (v2.0)

## 🚀 Project Vision
A **Multi-Tenant Blogging Platform** built in Go. 
1.  **Pangolin Cloud:** Managed service with free subdomains (`user.pangolin.com`).
2.  **Pangolin Core:** Open-source binary for self-hosting with custom domains.
3.  **Hybrid:** Users can migrate from Cloud to Self-Hosted anytime.

**Core Value Proposition:**
1.  **Free Tier:** Free subdomain hosting for everyone.
2.  **Custom Domains:** Bring your own domain (BYOD) with auto-SSL.
3.  **Performance:** Go-based multi-tenant engine (10x faster than WP.com).
4.  **Security:** Isolated tenant data, automated HTTPS.

---

## 📅 Phase 1: Core & Multi-Tenancy (Weeks 1-4)
**Goal:** A single server that can host multiple distinct blogs.

### 1.1 Project Setup
- [ ] Initialize Go Module (`go mod init pangolin`).
- [ ] Setup Directory Structure (Clean Architecture).
- [ ] Configure `.env` management (Viper).

### 1.2 Database & Multi-Tenancy
- [ ] Integrate **GORM**.
- [ ] Create `Site` Model (ID, Name, Subdomain, CustomDomain, Status).
- [ ] Create `User` Model (linked to Site).
- [ ] Create `Post`, `Category`, `Setting` Models (all linked to `site_id`).
- [ ] **Critical:** Implement **Global Middleware** to resolve `Site` from `Host` header on every request.
- [ ] Implement **Row-Level Security** (ensure Site A cannot access Site B's data).

### 1.3 Authentication
- [ ] Multi-tenant Auth (Users belong to specific Sites).
- [ ] JWT/Session handling scoped to Site.
- [ ] Admin Dashboard per Site.

**✅ Deliverable:** A server where `site1.pangolin.com` and `site2.pangolin.com` show different content from the same database.

---

## 📅 Phase 2: Domain & SSL Management (Weeks 5-8)
**Goal:** Allow users to connect custom domains securely.

### 2.1 Subdomain Logic
- [ ] Configure Wildcard DNS (`*.pangolin.com`).
- [ ] Logic to reserve subdomains (prevent squatting).
- [ ] Auto-generate SSL for subdomains (Wildcard Cert or Dynamic).

### 2.2 Custom Domain Logic
- [ ] UI for users to add Custom Domain (`www.mysite.com`).
- [ ] **Domain Verification:** Require DNS TXT record check before activation.
- [ ] **Auto-SSL:** Integrate **CertMagic** or **LEGO** to issue Let's Encrypt certs automatically.
- [ ] Renewal Logic (Auto-renew certs before expiry).

### 2.3 Routing Middleware
- [ ] Priority Logic: Custom Domain > Subdomain > Default.
- [ ] 404 Handling for unregistered domains.
- [ ] Redirect Logic (HTTP → HTTPS).

**✅ Deliverable:** Users can connect `myblog.com` and have it work securely with HTTPS automatically.

---

## 📅 Phase 3: Admin Dashboard & Content (Weeks 9-12)
**Goal:** A usable UI for writing and managing content.

### 3.1 Frontend Stack
- [ ] **HTMX + Tailwind CSS** (Embedded in binary).
- [ ] Multi-site aware Dashboard (Switch between sites if user owns multiple).

### 3.2 Content Features
- [ ] Markdown Editor + Live Preview.
- [ ] Media Manager (S3 compatible storage for multi-tenant files).
- [ ] SEO Settings per Post (Meta tags, OpenGraph).
- [ ] **Static Export:** `./pangolin export` (for self-hosted users).

### 3.3 Public Theme Engine
- [ ] Global Theme Repository.
- [ ] Site-specific Theme Customization.
- [ ] Server-Side Rendering (SSR) for speed.

**✅ Deliverable:** Full CMS functionality for end-users.

---

## 📅 Phase 4: Security & Isolation (Weeks 13-15)
**Goal:** Ensure tenants cannot hack each other.

### 4.1 Tenant Isolation
- [ ] Audit all DB queries for `site_id` leakage.
- [ ] File Storage Isolation (Folder per Site: `/uploads/site_id/...`).
- [ ] Rate Limiting per Site (prevent one blog from DDoSing the platform).

### 4.2 Platform Security
- [ ] WAF (Web Application Firewall) Middleware.
- [ ] DDoS Protection (integration with Cloudflare or custom rate limits).
- [ ] Admin Protection (2FA for Platform Admins).

**✅ Deliverable:** Secure multi-tenant environment ready for public beta.

---

## 📅 Phase 5: Launch & Scaling (Weeks 16+)
**Goal:** Go live with Pangolin Cloud.

### 5.1 Infrastructure
- [ ] Load Balancer (Traefik or Nginx) in front of Go instances.
- [ ] Database Clustering (Postgres Read Replicas).
- [ ] Object Storage (AWS S3 / Cloudflare R2) for media.

### 5.2 Billing (Optional)
- [ ] Integrate Stripe for Premium Plans (Custom Domains, Removal of Branding).
- [ ] Free Tier Limits (Storage, Bandwidth).

### 5.3 Marketing
- [ ] Launch on ProductHunt.
- [ ] "Migrate from WordPress" Campaign.
- [ ] Documentation & API Docs.

**✅ Deliverable:** Pangolin Cloud Live (v1.0).

---

## 🛠️ Updated Tech Stack

| Component | Technology | Why? |
| :--- | :--- | :--- |
| **Language** | Go (Golang) | Concurrency for multi-tenancy |
| **Framework** | Fiber v2 | Fast routing, middleware support |
| **Database** | PostgreSQL | Better for multi-tenant than SQLite |
| **SSL/TLS** | **CertMagic** | **Critical:** Auto-HTTPS for custom domains |
| **DNS** | Cloudflare API | Automate DNS verification |
| **Storage** | S3 Compatible | Isolate user uploads securely |
| **Frontend** | HTMX + Tailwind | Fast, embedded, no build pipeline |

---

## 🎯 Success Metrics (KPIs)
1.  **Onboarding:** User can sign up and get a live subdomain in < 2 minutes.
2.  **SSL:** 100% of custom domains have valid HTTPS automatically.
3.  **Isolation:** Zero cross-site data leaks during security audit.
4.  **Performance:** < 100ms TTFB even with 1,000 concurrent sites.

---

## 📝 Notes for Vibe Coding
- **Prompt Strategy:** Always specify "Multi-tenant context".
- **Example Prompt:** "Write a Fiber middleware that extracts the hostname, queries the 'sites' table for a matching domain or subdomain, and attaches the Site object to the context."
- **SSL Prompt:** "Show me how to integrate CertMagic with Fiber to handle automatic HTTPS for dynamic custom domains."
- **Security:** "How to ensure GORM queries always include `WHERE site_id = ?`?"
