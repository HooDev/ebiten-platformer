package engine

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// AssetManager handles loading and caching of game assets
type AssetManager struct {
	images      map[string]*ebiten.Image
	mu          sync.RWMutex
	assetDir    string
	useEmbedded bool
	embeddedFS  embed.FS
}

// AssetConfig holds configuration for the asset manager
type AssetConfig struct {
	AssetDir    string
	UseEmbedded bool
	EmbeddedFS  embed.FS
}

// NewAssetManager creates a new asset manager instance
func NewAssetManager(config AssetConfig) *AssetManager {
	return &AssetManager{
		images:      make(map[string]*ebiten.Image),
		assetDir:    config.AssetDir,
		useEmbedded: config.UseEmbedded,
		embeddedFS:  config.EmbeddedFS,
	}
}

// LoadImage loads an image asset and caches it
func (am *AssetManager) LoadImage(path string) (*ebiten.Image, error) {
	am.mu.RLock()
	if img, exists := am.images[path]; exists {
		am.mu.RUnlock()
		return img, nil
	}
	am.mu.RUnlock()

	var imgData image.Image

	if am.useEmbedded {
		// Load from embedded filesystem
		file, err := am.embeddedFS.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open embedded image %s: %w", path, err)
		}
		defer file.Close()

		imgData, err = png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode embedded image %s: %w", path, err)
		}
	} else {
		// Load from filesystem
		fullPath := filepath.Join(am.assetDir, path)
		file, err := os.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open image %s: %w", fullPath, err)
		}
		defer file.Close()

		imgData, err = png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image %s: %w", fullPath, err)
		}
	}

	ebitenImg := ebiten.NewImageFromImage(imgData)

	am.mu.Lock()
	am.images[path] = ebitenImg
	am.mu.Unlock()

	log.Printf("Loaded image: %s", path)
	return ebitenImg, nil
}

// GetImage retrieves a cached image or loads it if not cached
func (am *AssetManager) GetImage(path string) (*ebiten.Image, error) {
	return am.LoadImage(path)
}

// LoadAudio loads an audio file and caches the raw data
func (am *AssetManager) LoadAudio(path string) ([]byte, error) {
	// Audio loading will be implemented in Phase 4
	return nil, fmt.Errorf("audio loading not yet implemented")
}

// CreateAudioPlayer creates an audio player from a loaded audio file
func (am *AssetManager) CreateAudioPlayer(path string) error {
	// Audio player creation will be implemented in Phase 4
	return fmt.Errorf("audio player creation not yet implemented")
}

// PreloadAssets loads a list of assets into cache
func (am *AssetManager) PreloadAssets(imagePaths []string, audioPaths []string) error {
	var errors []string

	// Preload images
	for _, path := range imagePaths {
		if _, err := am.LoadImage(path); err != nil {
			errors = append(errors, fmt.Sprintf("image %s: %v", path, err))
		}
	}

	// Audio preloading will be implemented in Phase 4
	if len(audioPaths) > 0 {
		log.Println("Audio preloading will be available in Phase 4")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to preload assets:\n%s", strings.Join(errors, "\n"))
	}

	log.Printf("Successfully preloaded %d images", len(imagePaths))
	return nil
}

// GetLoadedImageCount returns the number of cached images
func (am *AssetManager) GetLoadedImageCount() int {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return len(am.images)
}

// GetLoadedAudioCount returns the number of cached audio files
func (am *AssetManager) GetLoadedAudioCount() int {
	// Audio counting will be implemented in Phase 4
	return 0
}

// ClearCache removes all cached assets (useful for memory management)
func (am *AssetManager) ClearCache() {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	am.images = make(map[string]*ebiten.Image)
	log.Println("Asset cache cleared")
}

// ListCachedAssets returns lists of all cached asset paths
func (am *AssetManager) ListCachedAssets() (images []string, audio []string) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for path := range am.images {
		images = append(images, path)
	}

	return images, audio
}
