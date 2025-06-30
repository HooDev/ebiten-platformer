package main

import (
	"testing"
	"ebiten-platformer/level"
	"ebiten-platformer/entities"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestSpecificPosition192(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "Test Level")
	
	// Create a platform from tile X=2 to X=5 at Y=5
	for x := 2; x <= 5; x++ {
		testLevel.SetTile(x, 5, level.TileSolid)
	}
	
	// Create collision adapter
	adapter := level.NewCollisionAdapter(testLevel)
	
	// Create player exactly at X=192 where the issue occurs
	spriteSheet := ebiten.NewImage(32, 32)
	player := entities.NewPlayer(192, 128, spriteSheet)
	player.SetLevel(adapter)
	
	deltaTime := 1.0 / 60.0
	
	t.Logf("Player bounds: X=192 to 224")
	t.Logf("Platform tiles:")
	for x := 2; x <= 5; x++ {
		t.Logf("  Tile %d: X=%d to %d", x, x*32, (x+1)*32)
	}
	
	// Manual collision check
	result := adapter.CheckCollision(192, 128, 32, 32)
	t.Logf("Manual collision result: OnGround=%t, Collided=%t", result.OnGround, result.Collided)
	
	// Check what tiles the player overlaps
	leftTile := int(192 / 32)   // Should be 6
	rightTile := int((192 + 32 - 1) / 32)  // Should be 7
	t.Logf("Player overlaps tiles X=%d to %d", leftTile, rightTile)
	
	// Check what's in those tiles
	for x := leftTile; x <= rightTile; x++ {
		tile := testLevel.GetTile(x, 5)
		t.Logf("Tile (%d,5): Type=%v", x, tile.Type)
	}
	
	// Update player
	player.Update(deltaTime)
	t.Logf("After update: OnGround=%t", player.IsOnGround())
}
