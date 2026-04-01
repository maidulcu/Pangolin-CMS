# Pangolin Phase 1 (MVP) - Test Plan

## Overview
This document outlines the testing strategy for Phase 1 (MVP) features: CLI commands, sitemap crawling, asset handling, config management, S3 deployment, and WordPress plugin.

---

## 1. Init Command

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| I01 | Initialize with valid site | Run `pangolin init -u https://example.com -k API_KEY` | Config saved to ~/.pangolin/pangolin.yaml |
| I02 | Initialize without URL | Run `pangolin init -k API_KEY` | Error: URL required |
| I03 | Initialize without API key | Run `pangolin init -u https://example.com` | Error: API key required |
| I04 | Initialize with custom config path | Run init, verify config location | Config at ~/.pangolin/ |
| I05 | Re-initialize overwrites | Run init twice | Second run overwrites first |

---

## 2. Sitemap Fetching

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| S01 | Fetch standard sitemap.xml | Configure site with sitemap.xml, run export | URLs discovered from sitemap |
| S02 | Fetch WordPress sitemap.xml | Configure site with wp-sitemap.xml, run export | URLs discovered correctly |
| S03 | No sitemap available | Site without sitemap | Error or empty result with message |
| S04 | Malformed sitemap | Invalid XML sitemap | Graceful error handling |
| S05 | Sitemap with large URL count | 1000+ URLs | All URLs discovered |

---

## 3. Export Command - Crawling

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| E01 | Basic full export | Run `pangolin export` | All pages exported to dist/ |
| E02 | Custom concurrency | Run `pangolin export -c 10` | Uses 10 concurrent requests |
| E03 | Custom output directory | Run `pangolin export -d output` | Files in output/ directory |
| E04 | Network timeout handling | Simulate slow/failing site | Errors logged, export continues |
| E05 | Handle redirects | Site with HTTP→HTTPS redirects | Follows redirects correctly |
| E06 | Concurrent request limit | Test with -c 1 vs -c 20 | Respects concurrency setting |

---

## 4. Asset Handling

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| A01 | Download images | Page with remote images | Images saved to dist/images/ |
| A02 | Download CSS | Page with external stylesheets | CSS saved to dist/assets/ |
| A03 | Download JS | Page with external scripts | JS saved to dist/assets/ |
| A04 | Rewrite image src | Exported page | src points to local /images/ |
| A05 | Rewrite CSS href | Exported page | href points to local /assets/ |
| A06 | Rewrite JS src | Exported page | src points to local /assets/ |
| A07 | Skip duplicate assets | Multiple pages share images | Asset downloaded once |
| A08 | Handle missing assets | Image returns 404 | Error logged, continue export |

---

## 5. Link Rewriting

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| L01 | Rewrite internal links | WordPress internal URLs | Converted to relative paths |
| L02 | Preserve external links | Links to google.com | Remain unchanged |
| L03 | Rewrite home page link | Link to site root | Rewritten to "/" |
| L04 | Handle query strings | URLs with ?param=value | Query strings preserved |
| L05 | Handle anchors | URLs with #section | Anchor preserved |

---

## 6. Config Management

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| C01 | Load existing config | Run any command after init | Config loaded successfully |
| C02 | Missing config | Run command without init | Error: "run init first" |
| C03 | Invalid config format | Manually corrupt config | Graceful error handling |
| C04 | Config precedence | Set via flag vs env vs config | Correct priority order |
| C05 | Config file permissions | Check file is not world-readable | Secure permissions (600) |

---

## 7. Serve Command

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| SV01 | Basic serve | Run `pangolin serve -d dist` | Serves files on port 8080 |
| SV02 | Custom port | Run `pangolin serve -p 3000` | Serves on port 3000 |
| SV03 | Serve non-existent directory | Run with invalid path | Error: directory not found |
| SV04 | Index file served | Request / | Returns index.html |
| SV05 | Subdirectory files | Request /about/ | Returns /about/index.html |
| SV06 | Static file types | Request /images/photo.jpg | Correct Content-Type |

---

## 8. S3 Deployment

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| S3-01 | Deploy to S3 | Run `pangolin deploy -b bucket -r region` | Files uploaded to S3 |
| S3-02 | Custom region | Deploy to eu-west-1 | Files in correct region |
| S3-03 | Missing bucket | Run without -b flag | Error: bucket required |
| S3-04 | Missing AWS credentials | No env vars set | Error: credentials required |
| S3-05 | Correct Content-Type | Upload .html, .css, .js | Each with proper MIME type |
| S3-06 | Large file upload | Deploy with large assets | Upload succeeds |
| S3-07 | Overwrite existing | Deploy twice | Files overwritten |

---

## 9. WordPress Plugin

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| WP01 | Plugin activation | Activate in WP Admin | No errors |
| WP02 | Generate API key | Click generate | Key displayed/copied |
| WP03 | API key auth | Use key with CLI | Authentication succeeds |
| WP04 | Invalid API key | Use wrong key | 401 Unauthorized |
| WP05 | Minimal permissions | Check user capabilities | Only edit_posts required |
| WP06 | Settings page | Navigate to Settings > Pangolin | Page loads without errors |
| WP07 | REST endpoint | GET /wp-json/pangolin/v1/posts | Returns post data |

---

## 10. Dashboard Command

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| D01 | Start dashboard | Run `pangolin dashboard` | Server starts on port 3000 |
| D02 | Access dashboard | Open http://localhost:3000 | Dashboard UI loads |
| D03 | Trigger export | Click export button | Export runs |
| D04 | View export history | Navigate to history | Past exports listed |
| D05 | Settings page | Access settings | Configurable options |

---

## 11. Error Handling

| Scenario | Expected Behavior |
|----------|------------------|
| Site unreachable | Clear error message with URL |
| Invalid sitemap URL | "Failed to fetch sitemap" |
| Disk full | "Failed to write file: no space" |
| Permission denied | "Permission denied" on config/file |
| Malformed HTML | Parser continues, logs warning |
| Large sitemap | Process without memory issues |

---

## 12. Performance Tests

| Operation | Metric | Target |
|-----------|--------|--------|
| 100 page export | Time | < 60 seconds |
| 1000 page export | Time | < 5 minutes |
| 50MB assets | Download time | < 2 minutes |
| Serve static files | Response time | < 100ms |
| S3 upload (100 files) | Time | < 2 minutes |

---

## 13. Test Data

### WordPress Test Sites
Use these for testing:
- WordPress.com test site
- Local WP installation with WP-CLI
- Docker WordPress image

### Sample Content
Create test site with:
- 10+ posts with featured images
- Custom CSS in theme
- JavaScript enqueued
- Multiple pages
- Categories and tags
- Custom post types (if applicable)

---

## 14. Manual Testing Checklist

### Core Functionality
- [ ] `pangolin init` saves config correctly
- [ ] `pangolin export` discovers all URLs
- [ ] Images downloaded and paths rewritten
- [ ] CSS/JS downloaded and paths rewritten
- [ ] Internal links converted to relative
- [ ] External links remain unchanged
- [ ] Export summary shows correct counts

### Deployment
- [ ] S3 upload with correct file types
- [ ] S3 upload overwrites existing files
- [ ] Local serve serves all file types

### Integration
- [ ] WP plugin generates valid API key
- [ ] CLI authenticates with WP API
- [ ] Dashboard displays export status

---

## 15. Reporting

After testing, document:
1. All test cases passed/failed
2. Any bugs found with severity
3. Performance metrics
4. Compatibility issues
5. Suggestions for improvement
