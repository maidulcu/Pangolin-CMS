package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewBundler(t *testing.T) {
	bundler := NewBundler(BundlerOptions{
		Minify:      true,
		Parallelism: 4,
	})

	if bundler == nil {
		t.Fatal("Expected non-nil Bundler")
	}

	if bundler.minify != true {
		t.Error("Expected minify to be true")
	}

	if bundler.parallelism != 4 {
		t.Error("Expected parallelism to be 4")
	}
}

func TestNewBundler_DefaultParallelism(t *testing.T) {
	bundler := NewBundler(BundlerOptions{
		Minify:      true,
		Parallelism: 0,
	})

	if bundler.parallelism != 4 {
		t.Errorf("Expected default parallelism 4, got %d", bundler.parallelism)
	}
}

func TestMinifyCSS(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})

	css := `
	/* This is a comment */
	.class {
		color: red;
		background: blue;
	}
	`

	minified := bundler.minifyCSS(css)

	if strings.Contains(minified, "/*") {
		t.Error("Expected comments to be removed")
	}

	if strings.Contains(minified, "\n") {
		t.Error("Expected newlines to be removed")
	}

	if strings.Contains(minified, "  ") {
		t.Error("Expected multiple spaces to be collapsed")
	}
}

func TestMinifyCSS_MultipleSelectors(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})

	css := `
	.class1 { color: red; }
	.class2 { color: blue; }
	`

	minified := bundler.minifyCSS(css)

	if !strings.Contains(minified, "}") {
		t.Error("Expected closing braces to be present")
	}

	if strings.Count(minified, "{") != 2 {
		t.Error("Expected 2 opening braces")
	}
}

func TestMinifyJS(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})

	js := `
	// This is a single line comment
	/* This is a
	multi-line comment */
	function test() {
		return true;
	}
	`

	minified := bundler.minifyJS(js)

	if strings.Contains(minified, "//") {
		t.Error("Expected single line comments to be removed")
	}

	if strings.Contains(minified, "/*") {
		t.Error("Expected multi-line comments to be removed")
	}
}

func TestMinifyJS_RemoveWhitespace(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})

	js := `
	function  test()  {
		return  true;
	}
	`

	minified := bundler.minifyJS(js)

	if strings.Contains(minified, "\n") {
		t.Error("Expected newlines to be removed")
	}
}

func TestBundleDirectory_NoAssetsDir(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})
	tmpDir := t.TempDir()

	err := bundler.BundleDirectory(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBundleDirectory_WithCSSFiles(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true, Parallelism: 2})
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	cssContent := `/* Comment */ .class { color: red; }`
	os.WriteFile(filepath.Join(assetsDir, "style.css"), []byte(cssContent), 0644)

	err := bundler.BundleDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minifiedFile := filepath.Join(assetsDir, "style.min.css")
	if _, err := os.Stat(minifiedFile); os.IsNotExist(err) {
		t.Error("Expected minified CSS file to be created")
	}
}

func TestBundleDirectory_WithJSFiles(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true, Parallelism: 2})
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	jsContent := `// Comment function test() { return true; }`
	os.WriteFile(filepath.Join(assetsDir, "script.js"), []byte(jsContent), 0644)

	err := bundler.BundleDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minifiedFile := filepath.Join(assetsDir, "script.min.js")
	if _, err := os.Stat(minifiedFile); os.IsNotExist(err) {
		t.Error("Expected minified JS file to be created")
	}
}

func TestBundleDirectory_NoMinify(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: false, Parallelism: 2})
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	cssContent := `/* Comment */ .class { color: red; }`
	os.WriteFile(filepath.Join(assetsDir, "style.css"), []byte(cssContent), 0644)

	err := bundler.BundleDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minifiedFile := filepath.Join(assetsDir, "style.min.css")
	if _, err := os.Stat(minifiedFile); os.IsNotExist(err) {
		t.Error("Expected minified CSS file to be created even without minify flag")
	}
}

func TestProcessCSS(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})
	tmpDir := t.TempDir()

	cssContent := `/* Comment */ .class { color: red; }`
	cssPath := filepath.Join(tmpDir, "style.css")
	os.WriteFile(cssPath, []byte(cssContent), 0644)

	err := bundler.processCSS(cssPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minifiedFile := strings.TrimSuffix(cssPath, ".css") + ".min.css"
	content, err := os.ReadFile(minifiedFile)
	if err != nil {
		t.Fatalf("Failed to read minified file: %v", err)
	}

	if strings.Contains(string(content), "/*") {
		t.Error("Comments should be removed from minified CSS")
	}
}

func TestProcessJS(t *testing.T) {
	bundler := NewBundler(BundlerOptions{Minify: true})
	tmpDir := t.TempDir()

	jsContent := `// Comment function test() { return true; }`
	jsPath := filepath.Join(tmpDir, "script.js")
	os.WriteFile(jsPath, []byte(jsContent), 0644)

	err := bundler.processJS(jsPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minifiedFile := strings.TrimSuffix(jsPath, ".js") + ".min.js"
	content, err := os.ReadFile(minifiedFile)
	if err != nil {
		t.Fatalf("Failed to read minified file: %v", err)
	}

	if strings.Contains(string(content), "//") {
		t.Error("Comments should be removed from minified JS")
	}
}

func TestCombineCSSFiles(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "style1.css")
	file2 := filepath.Join(tmpDir, "style2.css")
	output := filepath.Join(tmpDir, "combined.css")

	os.WriteFile(file1, []byte(".class1 { color: red; }"), 0644)
	os.WriteFile(file2, []byte(".class2 { color: blue; }"), 0644)

	err := CombineCSSFiles([]string{file1, file2}, output)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("Failed to read combined file: %v", err)
	}

	combined := string(content)
	if !strings.Contains(combined, ".class1") {
		t.Error("Expected .class1 in combined file")
	}
	if !strings.Contains(combined, ".class2") {
		t.Error("Expected .class2 in combined file")
	}
}

func TestCombineJSFiles(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "script1.js")
	file2 := filepath.Join(tmpDir, "script2.js")
	output := filepath.Join(tmpDir, "combined.js")

	os.WriteFile(file1, []byte("function one() { return 1; }"), 0644)
	os.WriteFile(file2, []byte("function two() { return 2; }"), 0644)

	err := CombineJSFiles([]string{file1, file2}, output)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("Failed to read combined file: %v", err)
	}

	combined := string(content)
	if !strings.Contains(combined, "function one()") {
		t.Error("Expected function one in combined file")
	}
	if !strings.Contains(combined, "function two()") {
		t.Error("Expected function two in combined file")
	}
}
