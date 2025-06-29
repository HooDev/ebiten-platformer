package engine

import (
	"testing"
)

// setupTestAssets creates temporary test assets for testing
func setupTestAssets(t *testing.T) string {
	return CreateTestAssets(t, StandardTestAssets())
}

func TestNewAssetManager(t *testing.T) {
	config := AssetConfig{
		AssetDir:    "test",
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	if am == nil {
		t.Fatal("NewAssetManager returned nil")
	}
	
	if am.assetDir != "test" {
		t.Errorf("Expected assetDir to be 'test', got '%s'", am.assetDir)
	}
	
	if am.useEmbedded != false {
		t.Errorf("Expected useEmbedded to be false, got %v", am.useEmbedded)
	}
	
	if am.images == nil {
		t.Error("Expected images map to be initialised")
	}
}

func TestLoadImage_Success(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// Test loading a valid image
	img, err := am.LoadImage("test.png")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	if img == nil {
		t.Error("LoadImage returned nil image")
	}
	
	// Verify image is cached
	if len(am.images) != 1 {
		t.Errorf("Expected 1 cached image, got %d", len(am.images))
	}
	
	// Test loading the same image again (should come from cache)
	img2, err := am.LoadImage("test.png")
	if err != nil {
		t.Fatalf("Failed to load cached image: %v", err)
	}
	
	if img != img2 {
		t.Error("Expected cached image to be the same instance")
	}
}

func TestLoadImage_FileNotFound(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// Test loading a non-existent image
	_, err := am.LoadImage("nonexistent.png")
	if err == nil {
		t.Error("Expected error when loading non-existent image")
	}
}

func TestGetImage(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// GetImage should load and cache the image
	img, err := am.GetImage("test.png")
	if err != nil {
		t.Fatalf("GetImage failed: %v", err)
	}
	
	if img == nil {
		t.Error("GetImage returned nil image")
	}
	
	// Verify it's cached
	if am.GetLoadedImageCount() != 1 {
		t.Errorf("Expected 1 cached image, got %d", am.GetLoadedImageCount())
	}
}

func TestPreloadAssets(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	imagePaths := []string{
		"player.png",
		"heart.png",
		"cat_sad.png",
		"cat_happy.png",
	}
	
	err := am.PreloadAssets(imagePaths, nil)
	if err != nil {
		t.Fatalf("PreloadAssets failed: %v", err)
	}
	
	// Verify all images are cached
	expectedCount := len(imagePaths)
	if am.GetLoadedImageCount() != expectedCount {
		t.Errorf("Expected %d cached images, got %d", expectedCount, am.GetLoadedImageCount())
	}
	
	// Verify we can get each preloaded asset
	for _, path := range imagePaths {
		img, err := am.GetImage(path)
		if err != nil {
			t.Errorf("Failed to get preloaded image %s: %v", path, err)
		}
		if img == nil {
			t.Errorf("Preloaded image %s is nil", path)
		}
	}
}

func TestPreloadAssets_PartialFailure(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	imagePaths := []string{
		"player.png",      // exists
		"nonexistent.png", // doesn't exist
		"heart.png",       // exists
	}
	
	err := am.PreloadAssets(imagePaths, nil)
	if err == nil {
		t.Error("Expected error when preloading non-existent assets")
	}
	
	// Should still have loaded the valid assets
	if am.GetLoadedImageCount() != 2 {
		t.Errorf("Expected 2 cached images after partial failure, got %d", am.GetLoadedImageCount())
	}
}

func TestCacheManagement(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// Load some images
	_, err := am.LoadImage("player.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}
	
	_, err = am.LoadImage("heart.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}
	
	// Check cache count
	if am.GetLoadedImageCount() != 2 {
		t.Errorf("Expected 2 cached images, got %d", am.GetLoadedImageCount())
	}
	
	// List cached assets
	images, audio := am.ListCachedAssets()
	if len(images) != 2 {
		t.Errorf("Expected 2 cached images in list, got %d", len(images))
	}
	if len(audio) != 0 {
		t.Errorf("Expected 0 cached audio files, got %d", len(audio))
	}
	
	// Clear cache
	am.ClearCache()
	
	if am.GetLoadedImageCount() != 0 {
		t.Errorf("Expected 0 cached images after clear, got %d", am.GetLoadedImageCount())
	}
}

func TestAudioFunctions_NotImplemented(t *testing.T) {
	config := AssetConfig{
		AssetDir:    "test",
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// Test audio loading (should return error as not implemented)
	_, err := am.LoadAudio("test.wav")
	if err == nil {
		t.Error("Expected error for unimplemented audio loading")
	}
	
	// Test audio player creation (should return error as not implemented)
	err = am.CreateAudioPlayer("test.wav")
	if err == nil {
		t.Error("Expected error for unimplemented audio player creation")
	}
	
	// Audio count should always be 0
	if am.GetLoadedAudioCount() != 0 {
		t.Errorf("Expected 0 audio files, got %d", am.GetLoadedAudioCount())
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := setupTestAssets(t)
	
	config := AssetConfig{
		AssetDir:    tmpDir,
		UseEmbedded: false,
	}
	
	am := NewAssetManager(config)
	
	// Test concurrent loading
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			// Each goroutine tries to load the same image
			_, err := am.LoadImage("test.png")
			if err != nil {
				t.Errorf("Concurrent LoadImage failed: %v", err)
			}
			
			// And check cache count
			count := am.GetLoadedImageCount()
			if count < 1 {
				t.Errorf("Expected at least 1 cached image, got %d", count)
			}
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Should only have cached the image once
	if am.GetLoadedImageCount() != 1 {
		t.Errorf("Expected 1 cached image after concurrent access, got %d", am.GetLoadedImageCount())
	}
}
