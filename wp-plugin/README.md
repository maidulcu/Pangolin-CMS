# Pangolin WordPress Plugin

Provides API key authentication for the Pangolin CLI tool.

## Installation

1. Upload the `pangolin` folder to your `/wp-content/plugins/` directory
2. Activate the plugin through the 'Plugins' menu in WordPress
3. Go to Settings → Pangolin to generate an API key

## Usage

### Generate API Key

1. Go to **Settings → Pangolin**
2. Click **Generate API Key**
3. Copy the generated key

### CLI Setup

```bash
pangolin init -u https://your-site.com -k YOUR_API_KEY
```

### Test Connection

```bash
# Validate API key
curl -H "Authorization: Bearer YOUR_API_KEY" https://your-site.com/wp-json/pangolin/v1/validate

# Get site info
curl -H "Authorization: Bearer YOUR_API_KEY" https://your-site.com/wp-json/pangolin/v1/info
```

## REST API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/wp-json/pangolin/v1/key` | POST | Manage Options | Generate API key |
| `/wp-json/pangolin/v1/key` | DELETE | Manage Options | Delete API key |
| `/wp-json/pangolin/v1/validate` | GET | None | Validate API key |
| `/wp-json/pangolin/v1/info` | GET | Bearer Token | Get site info |

## Requirements

- WordPress 5.0+
- PHP 7.4+
- User with `edit_posts` capability

## Security

- API key requires `edit_posts` capability
- Keys are stored in WordPress options (consider encrypting for production)
- Use HTTPS in production

## Changelog

### 1.0.0
- Initial release
- API key generation
- REST API endpoints
- Admin settings page
