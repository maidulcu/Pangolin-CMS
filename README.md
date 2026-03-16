# Pangolin

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/maidulcu/Pangolin-CMS)
[![GitHub Release](https://img.shields.io/github/v/release/maidulcu/Pangolin-CMS?label=release)](https://github.com/maidulcu/Pangolin-CMS/releases)

> **Export WordPress sites to blazing-fast static HTML** — Deploy anywhere. Zero PHP.

---

## Free Features ✅

All features below are **100% free** and open source (MIT License).

### 🔧 Core Engine

- **Concurrent Crawling** — High-performance page fetching using Go goroutines
- **Smart Asset Handling** — Auto-download & localize images, CSS, JS, fonts
- **URL Rewriting** — Seamlessly convert absolute WordPress URLs to relative paths
- **Sitemap Integration** — Auto-discover via `sitemap.xml` or `wp-sitemap.xml`
- **Detailed Reporting** — Export summaries with success/failure metrics and logs

### 🎛 Dashboard (Web UI)

- Real-time export progress with visual feedback
- Persistent history of exports and deployments
- Centralized settings management (site URL, API keys, S3 config)
- Analytics dashboard: page counts, asset totals, export duration

### 🔌 WordPress Plugin Companion

- Secure REST API integration with API key authentication
- Minimal permissions: requires only `edit_posts` capability
- Intuitive admin settings page for key management
- No frontend impact — runs silently in the background

---

## 📦 Installation

### Option 1: Build from Source

```bash
git clone https://github.com/maidulcu/Pangolin-CMS.git
cd Pangolin-CMS
go build -o pangolin .
```

### Option 2: Pre-built Binaries

Download the latest release for your platform:

👉 [Releases Page](https://github.com/maidulcu/Pangolin-CMS/releases)

### Option 3: Go Install (Go 1.21+)

```bash
go install github.com/maidulcu/Pangolin-CMS@latest
```

---

## 🏁 Quick Start

### Step 1: Install & Configure WordPress Plugin

1. Upload `wp-plugin/` to `/wp-content/plugins/` on your WordPress site
2. Activate **Pangolin Connector** via WordPress Admin → Plugins
3. Navigate to **Settings → Pangolin**
4. Click **Generate API Key** and copy the key

### Step 2: Initialize Pangolin CLI

```bash
pangolin init -u https://example.com -k YOUR_API_KEY
```

✅ Configuration is saved to `~/.pangolin/pangolin.yaml`

### Step 3: Export to Static HTML

```bash
pangolin export -d dist
```

📁 Output: Clean, self-contained static files in `./dist`

### Step 4: Deploy to AWS S3

```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"

pangolin deploy -b my-bucket -r us-east-1
```

### Step 5: Preview Locally (Optional)

```bash
pangolin serve -p 8080
```

---

## 📚 CLI Commands Reference

### `init` — Connect to WordPress

```bash
pangolin init [flags]
```

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--url` | `-u` | ✅ | WordPress site URL (e.g., `https://example.com`) |
| `--api-key` | `-k` | ✅ | API key generated from WP plugin |

### `export` — Generate Static Site

```bash
pangolin export [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--concurrency` | `-c` | `5` | Max concurrent page requests |
| `--dist` | `-d` | `"dist"` | Output directory for static files |

### `deploy` — Upload to S3

```bash
pangolin deploy [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--bucket` | `-b` | ✅ | Target S3 bucket name |
| `--region` | `-r` | `"us-east-1"` | AWS region |
| `--dist` | `-d` | `"dist"` | Directory to upload |

### `serve` — Local Preview Server

```bash
pangolin serve [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dist` | `-d` | `"dist"` | Directory to serve |
| `--port` | `-p` | `8080` | Local server port |

### `dashboard` — Web Management UI

```bash
pangolin dashboard
```

Launches at `http://localhost:3000` — manage exports, view logs, and configure settings via browser.

---

## ⚙️ Configuration

Pangolin uses a YAML config file at `~/.pangolin/pangolin.yaml`:

```yaml
site_url: "https://example.com"
api_key: "your-api-key"
s3_bucket: "my-static-site"
s3_region: "us-east-1"
```

🔐 Sensitive values like `api_key` can also be set via environment variables:

```bash
export PANGOLIN_API_KEY="your_key"
```

**Environment Variables (S3 Deployment):**

```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

---

## 🔍 How It Works

1. **Discovery** — Fetches `sitemap.xml` or `wp-sitemap.xml` to build the URL queue
2. **Crawling** — Concurrently fetches pages using configurable goroutines
3. **Processing** — Parses HTML, downloads remote assets, rewrites links to relative paths
4. **Exporting** — Writes clean, self-contained HTML files to the output directory
5. **Deploying** — Uploads to S3 with correct `Content-Type` headers and cache controls

---

## 🗂 Project Structure

```
Pangolin-CMS/
├── main.go                 # CLI entry point
├── cmd/
│   ├── init.go             # init command
│   ├── export.go           # export command
│   ├── deploy.go           # deploy command
│   ├── serve.go            # preview server
│   ├── dashboard.go        # dashboard command
│   └── internal/
│       ├── config/         # Config management
│       ├── sitemap/        # Sitemap fetching
│       ├── crawler/        # Page fetching & link rewriting
│       └── exporter/       # Export & S3 upload
├── dashboard/              # Web dashboard
│   ├── main.go             # Fiber server
│   └── views/              # HTMX templates
└── wp-plugin/              # WordPress plugin
```

---

## 🔐 Security Best Practices

- ✅ **No admin access required** — Plugin uses minimal `edit_posts` capability
- ✅ **Local execution** — CLI runs on your machine, no remote code execution on WordPress
- ✅ **API key isolation** — Keys stored locally in `~/.pangolin/`, never transmitted unnecessarily
- ✅ **REST API hardening** — Nonce verification + capability checks on all endpoints
- ✅ **No persistent connections** — Short-lived HTTP requests with timeout controls

> 🛡️ For production: rotate API keys periodically and restrict S3 bucket policies to least privilege.

---

## 💎 Pro Features (Coming Soon)

Unlock advanced capabilities with **Pangolin Pro**:

| Feature | Benefit |
|---------|---------|
| 🚀 One-Click Netlify Deploy | `pangolin deploy --platform netlify` with auto-site creation |
| 🖼 Smart Image Optimization | Auto-compress & convert images to WebP/AVIF during export |
| 🔄 Incremental Exports | Detect changed content via last-modified headers — export only deltas |
| 🪝 Webhook Auto-Sync | Trigger exports automatically on WordPress post publish/update |
| 🌍 CDN Cache Invalidation | Auto-purge Cloudflare, Fastly, or CloudFront after deploy |
| 🌐 Multi-Site Management | Handle multiple WordPress instances from one CLI config |
| ⏱ Scheduled Exports | Cron-compatible scheduling for automated static builds |
| 🎯 Priority Support | Dedicated Slack channel + 24h SLA for bug resolution |

🔔 **Join the waitlist:** Subscribe for Pro Updates *(coming soon)*

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Commit changes with clear messages: `git commit -m 'feat: add image optimization hook'`
4. Push and open a Pull Request

📖 See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, testing guidelines, and code standards.

**Development Requirements:**
- Go 1.21+
- Node.js 18+ (for dashboard assets)
- WordPress 5.8+ (for plugin testing)

---

## 🆘 Support & Community

- 🐛 **Report a Bug:** [GitHub Issues](https://github.com/maidulcu/Pangolin-CMS/issues)
- 💡 **Request a Feature:** [Feature Requests](https://github.com/maidulcu/Pangolin-CMS/issues)
- 💬 **Discussions:** [GitHub Discussions](https://github.com/maidulcu/Pangolin-CMS/discussions)

---

## 📜 License

Distributed under the [MIT License](LICENSE).

✅ Free for personal and commercial use. No attribution required, but stars are appreciated! ⭐

---

*Built with ❤️ by [Maidul](https://github.com/maidulcu) — Empowering WordPress to go static, one export at a time.*
