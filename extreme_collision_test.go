package main

import (
	"testing"
	"fmt"
	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

// TestExtremeSpeedCollision verifies tunneling prevention at extremely high speeds
func TestExtremeSpeedCollision(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "Extreme Speed Test")
	testLevel.SetTile(2, 5, level.TileSolid) // Platform at tile (2,5) = world position (64, 160) to (96, 192)
	
	// Create level adapter
	levelAdapter := level.NewCollisionAdapter(testLevel)
	
	// Create a test sprite sheet to avoid the nil pointer issue
	testSpriteSheet := entities.CreateTestSpriteSheet()
	
	// Test extreme speeds
	testCases := []struct {
		name string
		startY, velocityY float64
		expectedMaxY float64 // Player should not go below this Y position
	}{
		{"Extremely high speed", 10, 100.0, 135},   // Should not tunnel through
		{"Ludicrous speed", 0, 500.0, 135},         // Should not tunnel through  
		{"Teleportation speed", 0, 2000.0, 135},   // Should not tunnel through
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh player for each test
			player := entities.NewPlayer(64, tc.startY, testSpriteSheet)
			player.SetLevel(levelAdapter)
			
			// Set the falling velocity
			player.VelocityY = tc.velocityY
			
			maxYReached := tc.startY
			landed := false
			
			// Simulate until player lands or we hit max frames
			for i := 0; i < 200; i++ { // More frames for extreme speeds
				prevOnGround := player.OnGround
				player.Update(1.0/60.0) // 60 FPS
				
				_, y := player.GetPosition()
				if y > maxYReached {
					maxYReached = y
				}
				
				// Check if player just landed
				if !prevOnGround && player.OnGround {
					landed = true
					fmt.Printf("%s: Landed at Y=%.1f (max Y=%.1f) after %d frames\n", tc.name, y, maxYReached, i+1)
					break
				}
				
				// Check for tunneling - player should never go too far below the platform
				if y > tc.expectedMaxY {
					t.Errorf("%s: Player tunneled! Y=%.1f exceeds expected max Y=%.1f", tc.name, y, tc.expectedMaxY)
					return
				}
			}
			
			if !landed {
				_, finalY := player.GetPosition()
				t.Errorf("%s: Player didn't land after 200 frames, final position Y=%.1f", tc.name, finalY)
			}
			
			// Verify final position is reasonable
			if maxYReached > tc.expectedMaxY {
				t.Errorf("%s: Player went too far down (Y=%.1f), expected to stop around Y=%.1f", tc.name, maxYReached, tc.expectedMaxY)
			}
		})
	}
}

// TestFallingFromHeight verifies no variable sinking when falling from extreme heights
func TestFallingFromHeight(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(20, 20, 32, "Height Test")
	testLevel.SetTile(5, 15, level.TileSolid) // Platform at tile (5,15) = world position (160, 480) to (192, 512)
	
	// Create level adapter and test sprite sheet
	levelAdapter := level.NewCollisionAdapter(testLevel)
	testSpriteSheet := entities.CreateTestSpriteSheet()
	
	// Test falling from various extreme heights
	testCases := []struct {
		name string
		startY float64
	}{
		{"Fall from height 1", 50},
		{"Fall from height 2", 100}, 
		{"Fall from height 3", 200},
		{"Fall from height 4", 300},
		{"Fall from very high", 50},
	}
	
	expectedLandingY := 448.0 // Platform at Y=480, player height=32, so landing at Y=448
	tolerance := 1.0 // Allow 1 pixel variation
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh player for each test
			player := entities.NewPlayer(160, tc.startY, testSpriteSheet)
			player.SetLevel(levelAdapter)
			
			// Let gravity take effect
			landed := false
			
			// Simulate until player lands
			for i := 0; i < 500; i++ { // Lots of frames for high falls
				prevOnGround := player.OnGround
				player.Update(1.0/60.0) // 60 FPS
				
				// Check if player just landed
				if !prevOnGround && player.OnGround {
					_, y := player.GetPosition()
					landed = true
					fmt.Printf("%s: Landed at Y=%.1f after %d frames\n", tc.name, y, i+1)
					
					// Check landing consistency
					if y < expectedLandingY - tolerance || y > expectedLandingY + tolerance {
						t.Errorf("%s: Inconsistent landing Y=%.1f, expected around %.1f ±%.1f", 
							tc.name, y, expectedLandingY, tolerance)
					}
					
					break
				}
			}
			
			if !landed {
				_, finalY := player.GetPosition()
				t.Errorf("%s: Player didn't land after 500 frames, final position Y=%.1f", tc.name, finalY)
			}
		})
	}
}
