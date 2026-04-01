# Pangolin - Project Plan

## Overview
Pangolin is a CLI tool to export WordPress sites to static HTML for deployment to S3, Netlify, or other static hosting providers.

## Architecture

### Components
1. **CLI (Go)** - Main engine: `init`, `export`, `deploy`, `serve`, `dashboard`
2. **Connector (PHP)** - WordPress Plugin for Auth & Webhooks
3. **Dashboard (Go/HTMX)** - Web UI for export management

### Technology Stack
- **Language:** Go 1.21+
- **CLI Framework:** Cobra
- **HTML Parsing:** goquery
- **Config:** Viper
- **AWS SDK:** AWS SDK v2 for S3
- **Dashboard:** Fiber + HTMX + TailwindCSS
- **WP Plugin:** PHP (WordPress)

## Current Status: Phase 2 Complete ✅

### Completed Features
- [x] CLI with Cobra (init, export, deploy, serve, dashboard)
- [x] Sitemap fetching (sitemap.xml / wp-sitemap.xml)
- [x] Concurrent page crawling with goroutines
- [x] Link rewriting (absolute → relative URLs)
- [x] Static HTML export to local folder
- [x] Config management with Viper (saves to ~/.pangolin/)
- [x] S3 deployment with content-type detection
- [x] WordPress Plugin for API key auth
- [x] Netlify deployment support
- [x] Image optimization (WebP conversion)
- [x] CSS/JS minification
- [x] Incremental export (ETag/Last-Modified)

### Usage
```bash
# Initialize with WordPress site
pangolin init -u https://example.com -k YOUR_API_KEY

# Export to static HTML
pangolin export -c 5 -d dist

# Deploy to S3
pangolin deploy -b my-bucket -r us-east-1
```

## Future Enhancements

### Phase 2: Enhanced Features (COMPLETE)
- [x] Netlify deployment support
- [x] Image optimization (WebP)
- [x] CSS/JS minification
- [x] Incremental export (only changed pages)

### Phase 3: Pro Features
- [ ] Auto-sync on content change (webhooks)
- [ ] CDN cache invalidation (Cloudflare/Fastly/CloudFront)
- [ ] Multi-site support
- [ ] Scheduled exports
- [ ] Real-time dashboard updates
- [ ] Priority support

## Security Model
- No admin access required
- Binary runs locally (not on WP server)
- API key requires only `edit_posts` capability

## Project Structure
```
├── main.go                 # Entry point
├── go.mod                  # Go dependencies
├── cmd/
│   ├── init.go            # init command
│   ├── export.go          # export command
│   ├── deploy.go          # deploy command
│   ├── serve.go           # preview server
│   ├── dashboard.go       # dashboard command
│   └── internal/
│       ├── config/        # Config management
│       ├── sitemap/       # Sitemap fetching
│       ├── crawler/       # Page fetching & link rewriting
│       └── exporter/      # Export & S3 upload
├── dashboard/             # Web dashboard
│   ├── main.go           # Fiber server
│   └── views/            # HTMX templates
└── wp-plugin/            # WordPress plugin
```

## Configuration
Config is saved to `~/.pangolin/pangolin.yaml`:
```yaml
site_url: "https://example.com"
api_key: "your-api-key"
s3_bucket: ""
s3_region: "us-east-1"
```

## Environment Variables (for S3 deploy)
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

## Roadmap
1. **MVP** - CLI export to local folder ✅
2. **Phase 1** - Enhanced features (Netlify, image optimization)
3. **Phase 2** - Pro features (auto-sync, multi-site)
