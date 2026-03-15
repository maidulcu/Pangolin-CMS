# Pangolin Dashboard Plan

## Overview
Web-based dashboard for managing Pangolin exports, deployments, and site monitoring.

## Technology Stack
- **Framework:** Fiber (Go)
- **Frontend:** HTMX + TailwindCSS
- **Auth:** Session-based with password

## Features

### Core Features (Implemented)
- [x] Export management (trigger, monitor progress)
- [x] Deployment history
- [x] Site configuration
- [x] Preview exported site

### Pro Features
- [ ] Auto-sync on content change (webhooks)
- [ ] CDN cache invalidation
- [ ] Multi-site support
- [ ] Scheduled exports

## Dashboard Structure

```
dashboard/
├── main.go           # Dashboard entry point
├── handlers/         # HTTP handlers
├── views/           # HTMX templates
├── static/          # CSS, JS assets
└── middleware/      # Auth, logging
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Dashboard home |
| `/export` | POST | Trigger export |
| `/deploy` | POST | Trigger deploy |
| `/settings` | GET/POST | Site settings |
| `/logs` | GET | View export logs |

## UI Layout

```
┌─────────────────────────────────────────┐
│  Pangolin            [Site] [Logout]   │
├─────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌────────┐ │
│  │ Exported│  │  Pages  │  │ Status │ │
│  │  142    │  │  1,234  │  │  ✓ OK  │ │
│  └─────────┘  └─────────┘  └────────┘ │
├─────────────────────────────────────────┤
│  Recent Exports                         │
│  ─────────────────                      │
│  ✓ 2 min ago - 142 pages              │
│  ✓ 1 hour ago - 140 pages              │
├─────────────────────────────────────────┤
│  [Export]  [Deploy]  [Preview]         │
└─────────────────────────────────────────┘
```

## Implementation Phases

### Phase 1: Basic Dashboard (Done)
1. Setup Fiber app ✅
2. Basic HTML templates ✅
3. Export trigger via API call ✅
4. Progress display ✅

### Phase 2: Enhanced Features
1. Real-time updates (HTMX)
2. Deployment history
3. Settings management

### Phase 3: Pro Features
1. Authentication
2. Multi-site support
3. Webhook handlers
4. Scheduled jobs
