package engine

import (
	"os"
	"path/filepath"
)

// TestPNG contains a valid minimal PNG (1x1 transparent pixel)
var TestPNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk (13 bytes)
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 image
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, // bit depth 8, colour type 6 (RGBA), CRC
	0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, // IDAT chunk (10 bytes)
	0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x01, // Compressed data
	0x0D, 0x0A, 0x2D, 0xB4, // CRC
	0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82, // IEND chunk
}

// TempDirer interface for both testing.T and testing.B
type TempDirer interface {
	TempDir() string
	Fatalf(format string, args ...interface{})
}

// CreateTestAssets creates a temporary directory with test assets
func CreateTestAssets(t TempDirer, assetNames []string) string {
	tmpDir := t.TempDir()
	
	for _, filename := range assetNames {
		err := os.WriteFile(filepath.Join(tmpDir, filename), TestPNG, 0644)
		if err != nil {
			t.Fatalf("Failed to create test asset %s: %v", filename, err)
		}
	}
	
	return tmpDir
}

// CreateTestAssetsWithMap creates test assets from a map of filename to data
func CreateTestAssetsWithMap(t TempDirer, assets map[string][]byte) string {
	tmpDir := t.TempDir()
	
	for filename, data := range assets {
		err := os.WriteFile(filepath.Join(tmpDir, filename), data, 0644)
		if err != nil {
			t.Fatalf("Failed to create test asset %s: %v", filename, err)
		}
	}
	
	return tmpDir
}

// StandardTestAssets returns a list of commonly used test asset names
func StandardTestAssets() []string {
	return []string{
		"test.png",
		"player.png",
		"heart.png",
		"cat_sad.png",
		"cat_happy.png",
	}
}
