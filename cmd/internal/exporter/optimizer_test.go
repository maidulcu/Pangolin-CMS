package exporter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewImageOptimizer(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled:     true,
		Format:      "webp",
		Quality:     80,
		Parallelism: 4,
	})

	if optimizer == nil {
		t.Fatal("Expected non-nil ImageOptimizer")
	}

	if optimizer.quality != 80 {
		t.Errorf("Expected quality 80, got %d", optimizer.quality)
	}

	if optimizer.format != "webp" {
		t.Errorf("Expected format 'webp', got '%s'", optimizer.format)
	}

	if !optimizer.enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestNewImageOptimizer_Defaults(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled: true,
		Format:  "",
		Quality: 0,
	})

	if optimizer.quality != 80 {
		t.Errorf("Expected default quality 80, got %d", optimizer.quality)
	}

	if optimizer.format != "webp" {
		t.Errorf("Expected default format 'webp', got '%s'", optimizer.format)
	}
}

func TestOptimizeDirectory_NoImagesDir(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled: true,
		Format:  "webp",
	})
	tmpDir := t.TempDir()

	err := optimizer.OptimizeDirectory(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOptimizeDirectory_EmptyImagesDir(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled: true,
		Format:  "webp",
	})
	tmpDir := t.TempDir()
	imagesDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(imagesDir, 0755)

	err := optimizer.OptimizeDirectory(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOptimizeDirectory_Disabled(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled: false,
	})
	tmpDir := t.TempDir()
	imagesDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(imagesDir, 0755)

	err := optimizer.OptimizeDirectory(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOptimizeDirectory_WithNonImageFiles(t *testing.T) {
	optimizer := NewImageOptimizer(OptimizeOptions{
		Enabled: true,
		Format:  "webp",
	})
	tmpDir := t.TempDir()
	imagesDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(imagesDir, 0755)

	os.WriteFile(filepath.Join(imagesDir, "readme.txt"), []byte("text file"), 0644)

	err := optimizer.OptimizeDirectory(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
	}

	for _, tt := range tests {
		result := formatFileSize(tt.input)
		if result != tt.expected {
			t.Errorf("formatFileSize(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestGetFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!")
	os.WriteFile(testFile, content, 0644)

	size := getFileSize(testFile)
	if size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), size)
	}
}

func TestGetFileSize_NonExistent(t *testing.T) {
	size := getFileSize("/nonexistent/file.txt")
	if size != 0 {
		t.Errorf("Expected 0 for non-existent file, got %d", size)
	}
}

func TestCalculateSavings(t *testing.T) {
	tmpDir := t.TempDir()
	original := filepath.Join(tmpDir, "original.txt")
	optimized := filepath.Join(tmpDir, "optimized.txt")

	os.WriteFile(original, []byte("AAAAAAAABBBBBBBB"), 0644)
	os.WriteFile(optimized, []byte("AB"), 0644)

	savings, err := CalculateSavings(original, optimized)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if savings <= 0 {
		t.Error("Expected positive savings percentage")
	}
}

func TestCalculateSavings_ZeroOriginal(t *testing.T) {
	tmpDir := t.TempDir()
	original := filepath.Join(tmpDir, "original.txt")
	optimized := filepath.Join(tmpDir, "optimized.txt")

	os.WriteFile(original, []byte(""), 0644)
	os.WriteFile(optimized, []byte("AB"), 0644)

	_, err := CalculateSavings(original, optimized)
	if err == nil {
		t.Error("Expected error for zero original size")
	}
}

func TestOptimizeOptions_Validation(t *testing.T) {
	tests := []struct {
		name     string
		opts     OptimizeOptions
		expected struct {
			quality     int
			format      string
			parallelism int
		}
	}{
		{
			name: "valid options",
			opts: OptimizeOptions{
				Quality:     90,
				Format:      "webp",
				Parallelism: 8,
			},
			expected: struct {
				quality     int
				format      string
				parallelism int
			}{90, "webp", 8},
		},
		{
			name: "zero quality defaults to 80",
			opts: OptimizeOptions{
				Quality: 0,
			},
			expected: struct {
				quality     int
				format      string
				parallelism int
			}{80, "webp", 4},
		},
		{
			name: "zero parallelism defaults to 4",
			opts: OptimizeOptions{
				Parallelism: 0,
			},
			expected: struct {
				quality     int
				format      string
				parallelism int
			}{80, "webp", 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewImageOptimizer(tt.opts)

			if optimizer.quality != tt.expected.quality {
				t.Errorf("Expected quality %d, got %d", tt.expected.quality, optimizer.quality)
			}
			if optimizer.format != tt.expected.format {
				t.Errorf("Expected format %s, got %s", tt.expected.format, optimizer.format)
			}
			if optimizer.parallelism != tt.expected.parallelism {
				t.Errorf("Expected parallelism %d, got %d", tt.expected.parallelism, optimizer.parallelism)
			}
		})
	}
}
