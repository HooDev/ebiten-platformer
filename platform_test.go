package main

import (
	"testing"
	"ebiten-platformer/level"
	"ebiten-platformer/entities"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestPlayerFallsOffPlatform(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "Test Level")
	
	// Create a platform from tile X=2 to X=5 at Y=5
	for x := 2; x <= 5; x++ {
		testLevel.SetTile(x, 5, level.TileSolid)
	}
	
	// Create collision adapter
	adapter := level.NewCollisionAdapter(testLevel)
	
	// Create player on the platform
	spriteSheet := ebiten.NewImage(32, 32)
	
	// Position player at pixel coordinates that put them on the platform
	// Tile Y=5 means pixel Y = 5 * 32 = 160
	// We want player to be on top of the platform, so Y = 160 - 32 = 128
	player := entities.NewPlayer(128, 128, spriteSheet) // X=128 (tile 4), Y=128 (on platform)
	player.SetLevel(adapter)
	
	deltaTime := 1.0 / 60.0
	
	// Update to establish ground state
	player.Update(deltaTime)
	
	t.Logf("Initial state: Position=(%f,%f), OnGround=%t", player.X, player.Y, player.OnGround)
	
	// Verify player is on ground (on the platform)
	if !player.OnGround {
		t.Error("Player should be on ground initially")
	}
	
	// Simulate walking to the right off the platform
	// Move player further past X=192 (platform edge) to ensure they fall off
	for i := 0; i < 40; i++ { // Even more frames to move further
		player.MoveRight()
		player.Update(deltaTime)
		
		if i == 20 { // Log state partway through
			t.Logf("Partway: Position=(%f,%f), OnGround=%t, CoyoteTimer=%f, VelocityY=%f", 
				player.X, player.Y, player.OnGround, player.CoyoteTimer, player.VelocityY)
		}
	}
	
	t.Logf("After walking right: Position=(%f,%f), OnGround=%t, CoyoteTimer=%f, VelocityY=%f", 
		player.X, player.Y, player.OnGround, player.CoyoteTimer, player.VelocityY)
	
	// Player should be at or past the platform edge (X should be >= 192, allowing for floating point precision)
	if player.X < 191.999 {
		t.Errorf("Player should have reached platform edge (X>=192), got X=%f", player.X)
	}
	
	// Continue updating without input to see if player falls
	for i := 0; i < 30; i++ { // 30 more frames to see if player falls
		player.Update(deltaTime)
		
		if i == 10 {
			t.Logf("Falling frame 10: Position=(%f,%f), OnGround=%t, CoyoteTimer=%f, VelocityY=%f", 
				player.X, player.Y, player.OnGround, player.CoyoteTimer, player.VelocityY)
		}
		if i == 20 {
			t.Logf("Falling frame 20: Position=(%f,%f), OnGround=%t, CoyoteTimer=%f, VelocityY=%f", 
				player.X, player.Y, player.OnGround, player.CoyoteTimer, player.VelocityY)
		}
	}
	
	t.Logf("Final state: Position=(%f,%f), OnGround=%t, CoyoteTimer=%f, VelocityY=%f", 
		player.X, player.Y, player.OnGround, player.CoyoteTimer, player.VelocityY)
	
	// Player should be falling (positive VelocityY) 
	if player.VelocityY <= 0 {
		t.Errorf("Player should be falling (VelocityY > 0), got %f", player.VelocityY)
	}
	
	// Player should have fallen below the platform level
	platformY := float64(5 * 32) // Platform is at tile Y=5 = pixel Y=160
	if player.Y <= platformY {
		t.Errorf("Player should have fallen below platform level %f, got Y=%f", platformY, player.Y)
	}
	
	// At this point, coyote time behavior depends on whether the player pressed jump
	// during the fall. In this test, they didn't jump, so coyote time should have 
	// either expired naturally or not activated due to the gradual movement off the platform.
}
