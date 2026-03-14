# Pangolin

A CLI tool to export WordPress sites to static HTML for deployment to S3, Netlify, or other static hosting providers.

## Why Pangolin?

- **Performance**: Serve static HTML instead of dynamic PHP
- **Security**: No WordPress database or plugins exposed
- **Cost**: Host on S3, Cloudflare Pages, or Netlify for free/minimal cost
- **Simplicity**: No WordPress maintenance, updates, or security patches

## Features

- Concurrent page crawling with goroutines
- Automatic link rewriting (absolute → relative URLs)
- Sitemap discovery (supports WordPress native sitemaps)
- S3 deployment with automatic content-type detection
- Local preview server
- Web dashboard for export management
- Config stored in `~/.pangolin/`

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

### 1. Initialize

```bash
pangolin init -u https://example.com -k YOUR_API_KEY
```

### 2. Export to Static HTML

```bash
pangolin export -d dist
```

### 3. Deploy to S3

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

## WordPress Plugin

The WordPress plugin provides API key authentication. Install from [wp-plugin/](wp-plugin/) directory.

### Requirements

- WordPress 5.0+
- PHP 7.4+

### Installation

1. Upload the plugin to your WordPress
2. Activate the plugin
3. Generate an API key in Settings → Pangolin

## How It Works

1. **Sitemap Discovery**: Finds sitemap.xml or wp-sitemap.xml
2. **Crawling**: Fetches pages concurrently using goroutines
3. **Link Rewriting**: Converts absolute URLs to relative paths
4. **Export**: Saves as static HTML files
5. **Deploy**: Uploads to S3 with correct MIME types

## Project Structure

```
staticpress/
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

## Roadmap

- [x] MVP - CLI export to local folder
- [x] WordPress Plugin for auth
- [x] Preview mode (local server)
- [x] Dashboard for export management
- [ ] Netlify deployment
- [ ] Image optimization
- [ ] Incremental exports
- [ ] Auto-sync on content change

## License

MIT

## Contributing

1. Fork the repo
2. Create a feature branch
3. Submit a PR
