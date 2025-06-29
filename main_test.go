package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"ebiten-platformer/engine"
)

func TestNewRoboGame(t *testing.T) {
	game := NewRoboGame()
	
	if game == nil {
		t.Fatal("NewRoboGame returned nil")
	}
	
	if game.Game == nil {
		t.Error("RoboGame.Game is nil")
	}
	
	if game.playerImage != nil {
		t.Error("Expected playerImage to be nil before loading assets")
	}
	
	// Check initial state
	if game.GetState() != engine.StateLoading {
		t.Errorf("Expected initial state to be StateLoading, got %v", game.GetState())
	}
	
	// Check asset manager is available
	assetManager := game.GetAssetManager()
	if assetManager == nil {
		t.Error("Asset manager is nil")
	}
}

func TestRoboGame_LoadAssets_Success(t *testing.T) {
	// Create test game with temporary assets
	tmpDir := setupTestGameAssets(t)
	
	// Create game with custom asset directory
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    tmpDir,
			UseEmbedded: false,
		},
	}
	
	baseGame := engine.NewGame(config)
	game := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(480, 360),
	}
	
	// Load assets
	err := game.LoadAssets()
	if err != nil {
		t.Fatalf("LoadAssets failed: %v", err)
	}
	
	// Check that assets were loaded
	if game.playerImage == nil {
		t.Error("playerImage is nil after LoadAssets")
	}
	
	// Check that state changed to menu
	if game.GetState() != engine.StateMenu {
		t.Errorf("Expected state to be StateMenu after LoadAssets, got %v", game.GetState())
	}
	
	// Check that asset manager has the loaded image
	assetManager := game.GetAssetManager()
	if assetManager.GetLoadedImageCount() != 1 {
		t.Errorf("Expected 1 loaded image, got %d", assetManager.GetLoadedImageCount())
	}
}

func TestRoboGame_LoadAssets_Failure(t *testing.T) {
	// Create game with non-existent asset directory
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    "/nonexistent/directory",
			UseEmbedded: false,
		},
	}
	
	baseGame := engine.NewGame(config)
	game := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(480, 360),
	}
	
	// Load assets should succeed due to fallback
	err := game.LoadAssets()
	if err != nil {
		t.Errorf("LoadAssets should succeed with fallback, got error: %v", err)
	}
	
	// Player image should be the fallback test sprite sheet
	if game.playerImage == nil {
		t.Error("playerImage should not be nil after LoadAssets with fallback")
	}
	
	// State should transition to menu even with fallback
	if game.GetState() != engine.StateMenu {
		t.Errorf("Expected state to be StateMenu after LoadAssets with fallback, got %v", game.GetState())
	}
	
	// Player should be created successfully
	if game.player == nil {
		t.Error("player should not be nil after LoadAssets with fallback")
	}
	
	// Input handler should be created
	if game.inputHandler == nil {
		t.Error("inputHandler should not be nil after LoadAssets with fallback")
	}
}

func TestRoboGame_Update(t *testing.T) {
	tmpDir := setupTestGameAssets(t)
	
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    tmpDir,
			UseEmbedded: false,
		},
	}
	
	baseGame := engine.NewGame(config)
	game := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(480, 360),
	}
	
	// Load assets first
	err := game.LoadAssets()
	if err != nil {
		t.Fatalf("LoadAssets failed: %v", err)
	}
	
	// Test update in playing state
	err = game.Update()
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
}

func TestRoboGame_GameStateTransitions(t *testing.T) {
	tmpDir := setupTestGameAssets(t)
	
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    tmpDir,
			UseEmbedded: false,
		},
	}
	
	baseGame := engine.NewGame(config)
	game := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(480, 360),
	}
	
	// Initial state should be loading
	if game.GetState() != engine.StateLoading {
		t.Errorf("Expected initial state StateLoading, got %v", game.GetState())
	}
		// Load assets - should transition to menu
	err := game.LoadAssets()
	if err != nil {
		t.Fatalf("LoadAssets failed: %v", err)
	}

	if game.GetState() != engine.StateMenu {
		t.Errorf("Expected state StateMenu after LoadAssets, got %v", game.GetState())
	}
	
	// Test state changes
	game.SetState(engine.StatePaused)
	if game.GetState() != engine.StatePaused {
		t.Errorf("Expected state StatePaused after SetState, got %v", game.GetState())
	}
}

func TestRoboGame_AssetManagerIntegration(t *testing.T) {
	// Create test assets with additional ones
	additionalAssets := []string{"heart.png", "cat_sad.png", "cat_happy.png"}
	tmpDir := setupTestGameAssetsWithExtras(t, additionalAssets)
	
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    tmpDir,
			UseEmbedded: false,
		},
	}
	
	baseGame := engine.NewGame(config)
	game := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(480, 360),
	}
	
	// Load initial assets
	err := game.LoadAssets()
	if err != nil {
		t.Fatalf("LoadAssets failed: %v", err)
	}
	
	// Test loading additional assets through asset manager
	assetManager := game.GetAssetManager()
	
	for _, asset := range additionalAssets {
		img, err := assetManager.LoadImage(asset)
		if err != nil {
			t.Errorf("Failed to load additional asset %s: %v", asset, err)
		}
		if img == nil {
			t.Errorf("Additional asset %s is nil", asset)
		}
	}
	
	// Check total loaded assets (player.png + additional assets)
	expectedCount := 1 + len(additionalAssets)
	if assetManager.GetLoadedImageCount() != expectedCount {
		t.Errorf("Expected %d loaded assets, got %d", expectedCount, assetManager.GetLoadedImageCount())
	}
	
	// Test preloading - add more assets to the existing directory
	moreAssets := []string{"platform.png", "spike.png"}
	addAssetsToDir(t, tmpDir, moreAssets)
	
	err = assetManager.PreloadAssets(moreAssets, nil)
	if err != nil {
		t.Errorf("Failed to preload assets: %v", err)
	}
	
	// Check final count
	finalExpectedCount := expectedCount + len(moreAssets)
	if assetManager.GetLoadedImageCount() != finalExpectedCount {
		t.Errorf("Expected %d total loaded assets, got %d", finalExpectedCount, assetManager.GetLoadedImageCount())
	}
}

// BenchmarkRoboGame_LoadAssets benchmarks the asset loading performance
func BenchmarkRoboGame_LoadAssets(b *testing.B) {
	// Create test assets once
	tmpDir := engine.CreateTestAssets(b, []string{"player.png"})
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		config := engine.GameConfig{
			ScreenWidth:  480,
			ScreenHeight: 360,
			AssetConfig: engine.AssetConfig{
				AssetDir:    tmpDir,
				UseEmbedded: false,
			},
		}
		
		baseGame := engine.NewGame(config)
		game := &RoboGame{
			Game: baseGame,
			overlayImage: ebiten.NewImage(480, 360),
		}
		
		err := game.LoadAssets()
		if err != nil {
			b.Fatalf("LoadAssets failed: %v", err)
		}
	}
}
