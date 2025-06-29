package main

import (
	"os"
	"path/filepath"
	
	"ebiten-platformer/engine"
)

// setupTestGameAssets creates temporary test assets for the main game tests
func setupTestGameAssets(t engine.TempDirer) string {
	return engine.CreateTestAssets(t, []string{"player.png"})
}

// setupTestGameAssetsWithExtras creates test assets including additional ones
func setupTestGameAssetsWithExtras(t engine.TempDirer, extraAssets []string) string {
	allAssets := append([]string{"player.png"}, extraAssets...)
	return engine.CreateTestAssets(t, allAssets)
}

// addAssetsToDir adds more assets to an existing directory
func addAssetsToDir(t engine.TempDirer, dir string, assetNames []string) {
	for _, assetName := range assetNames {
		err := os.WriteFile(filepath.Join(dir, assetName), engine.TestPNG, 0644)
		if err != nil {
			t.Fatalf("Failed to add test asset %s: %v", assetName, err)
		}
	}
}
