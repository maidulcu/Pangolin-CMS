# Pangolin

A CLI tool to export WordPress sites to static HTML for deployment to S3, Netlify, or other static hosting providers.

## Why Pangolin?

- **Performance**: Serve static HTML instead of dynamic PHP
- **Security**: No WordPress database or plugins exposed
- **Cost**: Host on S3, Cloudflare Pages, or Netlify for free/minimal cost
- **Simplicity**: No WordPress maintenance, updates, or security patches

## Free Features ✅

All features below are **100% free** and open source (MIT License).

### CLI Commands
- `init` - Initialize with WordPress site
- `export` - Export site to static HTML
- `deploy` - Deploy to S3
- `serve` - Local preview server
- `dashboard` - Web UI for management

### Core Functionality
- Concurrent page crawling (goroutines)
- **Automatic asset downloading** (images, CSS, JS)
- **Automatic link rewriting** (absolute → relative URLs)
- Sitemap discovery (sitemap.xml / wp-sitemap.xml)
- Export summary with success/fail counts

### Dashboard Features
- Real-time progress updates
- Export/deploy history with persistence
- Settings management (site URL, API key, S3 config)
- Stats dashboard (pages, assets, totals)

### WordPress Plugin
- API key authentication
- REST API endpoints
- Admin settings page
- Requires only `edit_posts` capability

## Pro Features 🔒

These features will be available in the paid version.

- **Netlify deployment** - One-click deploy to Netlify
- **Image optimization** - Compress images during export
- **Incremental exports** - Only export changed pages
- **Auto-sync** - Webhook triggers automatic exports
- **CDN cache invalidation** - Auto-clear CloudFlare/Fastly cache
- **Multi-site support** - Manage multiple WordPress sites
- **Scheduled exports** - cron-based automatic exports
- **Priority support** - Faster issue resolution

[Subscribe for Pro →](#) *(coming soon)*

## Installation

### From Source

```bash
git clone https://github.com/pangolin-cms/staticpress.git
cd staticpress
go build -o pangolin .
```

### Pre-built Binaries

Download from [Releases](https://github.com/pangolin-cms/staticpress/releases)

## Quick Start

### 1. Install WordPress Plugin

1. Upload `wp-plugin/` to your `/wp-content/plugins/` directory
2. Activate in WordPress admin
3. Go to **Settings → Pangolin**
4. Click **Generate API Key**

### 2. Initialize CLI

```bash
pangolin init -u https://example.com -k YOUR_API_KEY
```

### 3. Export Site

```bash
pangolin export -d dist
```

### 4. Deploy to S3

```bash
pangolin deploy -b my-bucket -r us-east-1
```

## Commands

### init

Initialize Pangolin with your WordPress site.

```bash
pangolin init [flags]
```

Flags:
- `-u, --url` - WordPress site URL (required)
- `-k, --api-key` - API key from WP plugin (required)

### export

Export WordPress site to static HTML.

```bash
pangolin export [flags]
```

Flags:
- `-c, --concurrency` - Number of concurrent requests (default: 5)
- `-d, --dist` - Output directory (default: "dist")

### deploy

Deploy static files to S3.

```bash
pangolin deploy [flags]
```

Flags:
- `-b, --bucket` - S3 bucket name (required)
- `-r, --region` - AWS region (default: "us-east-1")
- `-d, --dist` - Directory to deploy (default: "dist")

### serve

Start a local server to preview the exported site.

```bash
pangolin serve [flags]
```

Flags:
- `-d, --dist` - Directory to serve (default: "dist")
- `-p, --port` - Port to listen on (default: 8080)

### dashboard

Start the web dashboard for managing exports.

```bash
pangolin dashboard
```

Starts on http://localhost:3000

## Configuration

Config is stored at `~/.pangolin/pangolin.yaml`:

```yaml
site_url: "https://example.com"
api_key: "your-api-key"
s3_bucket: ""
s3_region: "us-east-1"
```

## Environment Variables

For S3 deployment:

```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

## How It Works

1. **Sitemap Discovery**: finds sitemap.xml or wp-sitemap.xml
2. **Crawling**: fetches pages concurrently using goroutines
3. **Asset Download**: downloads images, CSS, JS locally
4. **Link Rewriting**: converts absolute URLs to relative paths
5. **Export**: saves as static HTML files
6. **Deploy**: uploads to S3 with correct MIME types

## Project Structure

```
pangolin/
├── main.go                 # CLI entry point
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

## Security

- No admin access required (uses REST API)
- Binary runs locally, not on WordPress server
- API key requires only `edit_posts` capability
- Config stored in user's home directory

## License

[MIT](LICENSE) - Free for personal and commercial use.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Support

- [Report Bugs](https://github.com/pangolin-cms/staticpress/issues)
- [Request Features](https://github.com/pangolin-cms/staticpress/issues)
