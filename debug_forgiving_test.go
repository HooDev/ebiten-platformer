package main

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

func TestDebugForgivingDetection(t *testing.T) {
	// Create a test level with a narrow platform
	testLevel := level.NewLevel(10, 10, 32, "debug_forgiving_test")
	
	// Create a single tile platform at X=3 (96-128)
	testLevel.SetTile(3, 7, level.TileSolid)
	
	// Create a 1x1 pixel image for testing
	testImage := ebiten.NewImage(1, 1)
	
	// Position player so they have good overlap with the platform initially
	player := entities.NewPlayer(90, 192, testImage) // Good overlap with tile at X=96-128
	adapter := level.NewCollisionAdapter(testLevel)
	player.SetLevel(adapter)
	
	deltaTime := 1.0 / 60.0
	
	// Update player - should be on ground with good overlap
	player.Update(deltaTime)
	x, y := player.GetPosition()
	fmt.Printf("Good overlap: OnGround=%v, Position=(%.1f, %.1f), %s\n", 
		player.IsOnGround(), x, y, player.GetDebugInfo())
	
	// Move player to have partial overlap (less than 50% but still some contact)
	player.SetPosition(75, 192) // Spans 75-107, tile spans 96-128, overlap is 96-107 = 11 pixels = 34%
	x, y = player.GetPosition()
	fmt.Printf("Before update with partial overlap: Position=(%.1f, %.1f)\n", x, y)
	
	player.Update(deltaTime)
	x, y = player.GetPosition()
	fmt.Printf("After update with partial overlap: OnGround=%v, Position=(%.1f, %.1f), %s\n", 
		player.IsOnGround(), x, y, player.GetDebugInfo())
}
