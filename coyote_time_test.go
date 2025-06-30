package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

// Main coyote time integration tests using realistic level collision detection.
// These tests verify coyote time works correctly when walking off platforms.

func TestCoyoteTime(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "coyote_test")
	
	// Create a smaller platform: solid tiles at Y=7 for X=2,3 (only 2 tiles wide)
	testLevel.SetTile(2, 7, level.TileSolid)
	testLevel.SetTile(3, 7, level.TileSolid)
	
	// Create a 1x1 pixel image for testing
	testImage := ebiten.NewImage(1, 1)
	
	// Create player on the platform, closer to the edge
	player := entities.NewPlayer(80, 192, testImage) // Position closer to edge
	adapter := level.NewCollisionAdapter(testLevel)
	player.SetLevel(adapter)
	
	// Simulate first frame - player should be on ground
	deltaTime := 1.0 / 60.0 // 60 FPS
	player.Update(deltaTime)
	
	if !player.IsOnGround() {
		t.Fatal("Player should be on ground initially")
	}
	
	// Move player off the edge of the platform (need enough frames)
	for i := 0; i < 25; i++ {
		player.MoveRight()
		player.Update(deltaTime)
	}
	
	// Player should now be off ground but still have coyote time
	if player.IsOnGround() {
		t.Error("Player should not be on ground after walking off platform")
	}
	
	if player.GetCoyoteTimer() <= 0 {
		t.Error("Player should have coyote time remaining after leaving platform")
	}
	
	// Player should be able to jump during coyote time
	player.Jump()
	
	if player.GetVelocityY() >= 0 {
		t.Error("Player should have negative Y velocity after jumping during coyote time")
	}
	
	if player.GetCoyoteTimer() > 0 {
		t.Error("Coyote timer should be consumed after jumping")
	}
}

func TestCoyoteTimeExpires(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "coyote_expire_test")
	
	// Create a platform
	testLevel.SetTile(2, 7, level.TileSolid)
	testLevel.SetTile(3, 7, level.TileSolid)
	
	// Create a 1x1 pixel image for testing
	testImage := ebiten.NewImage(1, 1)
	
	// Create player on the platform
	player := entities.NewPlayer(80, 192, testImage)
	adapter := level.NewCollisionAdapter(testLevel)
	player.SetLevel(adapter)
	
	deltaTime := 1.0 / 60.0 // 60 FPS
	
	// Initialize player on ground
	player.Update(deltaTime)
	
	// Move player off the platform
	for i := 0; i < 25; i++ {
		player.MoveRight()
		player.Update(deltaTime)
	}
	
	// Player should not be on ground but should have coyote time
	if player.IsOnGround() {
		t.Error("Player should not be on ground after moving off platform")
	}
	
	if player.GetCoyoteTimer() <= 0 {
		t.Error("Player should have coyote time after leaving platform")
	}
	
	// Wait for coyote time to expire (0.1 seconds = 6 frames at 60 FPS)
	for i := 0; i < 7; i++ {
		player.Update(deltaTime)
	}
	
	// Coyote time should now be expired
	if player.GetCoyoteTimer() > 0 {
		t.Error("Coyote time should have expired")
	}
	
	// Player should not be able to jump now
	initialVelocityY := player.GetVelocityY()
	player.Jump()
	
	if player.GetVelocityY() != initialVelocityY {
		t.Error("Player should not be able to jump after coyote time expires")
	}
}

func TestCoyoteTimeForgivingDetection(t *testing.T) {
	// Create a test level with a narrow platform
	testLevel := level.NewLevel(10, 10, 32, "coyote_forgiving_test")
	
	// Create a single tile platform
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
	
	if !player.IsOnGround() {
		t.Error("Player should be on ground with good overlap")
	}
	
	// Move player gradually to edge to trigger coyote time naturally
	// Move left to reduce overlap below 50%
	for i := 0; i < 10; i++ {
		player.MoveLeft()
		player.Update(deltaTime)
		
		if !player.IsOnGround() && player.GetCoyoteTimer() > 0 {
			// Successfully triggered coyote time
			return
		}
	}
	
	t.Error("Player should have triggered coyote time by moving to platform edge")
}
