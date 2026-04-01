# Pangolin Phase 2 - Test Plan

## Overview
This document outlines the testing strategy for Phase 2 features: Netlify deployment, Image optimization, CSS/JS minification, and Incremental export.

---

## 1. Netlify Deployment

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| N01 | Deploy to Netlify with valid credentials | 1. Set NETLIFY_AUTH_TOKEN env var<br>2. Run `pangolin deploy --platform netlify --netlify-site <site-id>` | Files uploaded, deploy URL returned |
| N02 | Deploy fails without token | Run deploy without token | Error: "Netlify auth token required" |
| N03 | Deploy fails without site ID | Provide token but no site ID | Error: "Netlify site ID required" |
| N04 | List available sites | Run with token | Lists user's Netlify sites |
| N05 | Upload directory structure | Deploy with nested directories | All files uploaded maintaining structure |

### Prerequisites
```bash
export NETLIFY_AUTH_TOKEN="your-token"
```

---

## 2. Image Optimization

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| IO01 | Optimize images to WebP | 1. Create dist/images/ with .jpg/.png files<br>2. Run `pangolin export --optimize-images` | Images converted to .webp format |
| IO02 | Quality setting affects output | Run with `--image-quality 50` and `--image-quality 90` | Different file sizes |
| IO03 | No images directory | Run optimization on dir without images | Message: "No images to optimize" |
| IO04 | Unsupported format handling | Pass --image-format avif | Error or graceful fallback |

### Test Setup
```bash
mkdir -p test_dist/images
# Add test images to test_dist/images/
```

---

## 3. CSS/JS Minification

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| MJ01 | Minify CSS files | 1. Create dist/assets/style.css with comments/whitespace<br>2. Run `pangolin export --minify` | Minified .min.css created, smaller size |
| MJ02 | Minify JS files | 1. Create dist/assets/script.js with comments/whitespace<br>2. Run `pangolin export --minify` | Minified .min.js created, smaller size |
| MJ03 | Combined with image optimization | Run with both --optimize-images --minify | Both processes run sequentially |
| MJ04 | No assets directory | Run minify on dir without assets | Message: "No assets directory found" |

### Test Setup
```bash
mkdir -p test_dist/assets
echo "/* Comment */ .class { color: red; }" > test_dist/assets/style.css
echo "// Comment function test() { return true; }" > test_dist/assets/script.js
```

---

## 4. Incremental Export

### Test Cases

| ID | Description | Steps | Expected Result |
|----|-------------|-------|-----------------|
| IE01 | Full export first time | Run `pangolin export --incremental` | All pages exported |
| IE02 | Skip unchanged pages | Run export again without changes | Message: "No changes detected" |
| IE03 | Export changed page | Modify a page, run incremental | Only changed page re-exported |
| IE04 | Clear cache | Run `pangolin export --incremental --clear-cache` | Cache deleted, full re-export |
| IE05 | Cache file created | Run incremental export | .pangolin/export_cache.json created |

### Test Setup
```bash
# Requires configured WordPress site
pangolin init -u https://your-site.com -k YOUR_API_KEY
```

---

## 5. Integration Tests

### Full Pipeline Test
```bash
# Complete workflow
pangolin init -u https://example.com -k API_KEY
pangolin export -c 5 -d dist --optimize-images --minify
pangolin deploy -p s3 -b my-bucket -r us-east-1

# Incremental workflow
pangolin export --incremental
pangolin deploy -p netlify --netlify-site my-site
```

### Flag Combinations
| Command | Expected Behavior |
|---------|------------------|
| `export` | Basic full export |
| `export --optimize-images` | Export + image optimization |
| `export --minify` | Export + minification |
| `export --incremental` | Incremental export |
| `export --optimize-images --minify --incremental` | All features combined |

---

## 6. Manual Testing Checklist

- [ ] Sitemap discovery works with sitemap.xml
- [ ] Sitemap discovery works with wp-sitemap.xml
- [ ] Assets downloaded to /images/ and /assets/
- [ ] Links rewritten to relative paths
- [ ] Config saved to ~/.pangolin/pangolin.yaml
- [ ] S3 deploy uploads with correct Content-Type
- [ ] Netlify deploy creates new deploy
- [ ] Image optimization produces WebP files
- [ ] CSS minification removes comments/whitespace
- [ ] JS minification removes comments/whitespace
- [ ] Incremental export skips unchanged pages
- [ ] Clear cache forces full re-export

---

## 7. Performance Benchmarks (Optional)

| Operation | Metric | Target |
|-----------|--------|--------|
| Full export (100 pages) | Time | < 60s |
| Image optimization (50 images) | Time | < 30s |
| CSS minification (10 files) | Time | < 5s |
| Incremental (no changes) | Time | < 5s |

---

## 8. Error Handling Tests

| Scenario | Expected Behavior |
|----------|-------------------|
| Invalid sitemap URL | Clear error message |
| Network timeout | Retry with timeout message |
| Invalid Netlify token | 401 error from API |
| Disk full during export | Error with disk space message |
| Corrupted image file | Skip with warning, continue |

---

## 9. Test Data

### Sample WordPress Site
For testing, use a WordPress instance with:
- 10+ pages/posts
- Featured images
- Custom CSS/JS
- Multiple authors

### Sample Static Files
Create test_dist with:
```
test_dist/
├── index.html
├── about/
│   └── index.html
├── images/
│   ├── photo.jpg
│   └── logo.png
└── assets/
    ├── style.css
    └── script.js
```

---

## 10. Reporting

After testing, document:
1. All test cases passed/failed
2. Any bugs found
3. Performance metrics
4. Suggestions for improvement
