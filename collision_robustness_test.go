package main

import (
	"testing"
	"fmt"
	"math"
	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

// TestHighSpeedCollision verifies that high-speed movement doesn't cause tunneling
func TestHighSpeedCollision(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "High Speed Test")
	testLevel.SetTile(2, 5, level.TileSolid) // Platform at tile (2,5) = world position (64, 160) to (96, 192)
	
	// Create level adapter
	levelAdapter := level.NewCollisionAdapter(testLevel)
	
	// Create a test sprite sheet to avoid the nil pointer issue
	testSpriteSheet := entities.CreateTestSpriteSheet()
	
	// Test different speeds to ensure no tunneling occurs
	testCases := []struct {
		name string
		startY, velocityY float64
		expectedMaxY float64 // Player should not go below this Y position
	}{
		{"Normal speed", 100, 8.0, 135}, // Should stop around Y=128-135
		{"High speed", 80, 15.0, 135},   // Should not tunnel through
		{"Very high speed", 50, 25.0, 135}, // Should not tunnel through
		{"Extreme speed", 20, 40.0, 135},   // Should not tunnel through
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
			for i := 0; i < 100; i++ {
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
					t.Errorf("%s: Player tunneled too far! Y=%.1f exceeds expected max Y=%.1f", tc.name, y, tc.expectedMaxY)
					return
				}
			}
			
			if !landed {
				t.Errorf("%s: Player didn't land after 100 frames", tc.name)
			}
			
			// Verify final position is reasonable
			if maxYReached > tc.expectedMaxY {
				t.Errorf("%s: Player went too far down (Y=%.1f), expected to stop around Y=%.1f", tc.name, maxYReached, tc.expectedMaxY)
			}
		})
	}
}

// TestConsistentLanding verifies that player landing is consistent regardless of velocity
func TestConsistentLanding(t *testing.T) {
	// Create a test level with a platform
	testLevel := level.NewLevel(10, 10, 32, "Landing Test")
	testLevel.SetTile(2, 5, level.TileSolid) // Platform at tile (2,5) = world position (64, 160) to (96, 192)
	
	// Create level adapter and test sprite sheet
	levelAdapter := level.NewCollisionAdapter(testLevel)
	testSpriteSheet := entities.CreateTestSpriteSheet()
	
	// Test different landing scenarios
	testCases := []struct {
		name string
		startY, velocityY float64
		expectedYRange [2]float64 // [min, max] expected final Y position
	}{
		{"Slow landing", 120, 2.0, [2]float64{125, 135}},
		{"Medium landing", 100, 5.0, [2]float64{125, 135}}, 
		{"Fast landing", 80, 10.0, [2]float64{125, 135}},
		{"Very fast landing", 50, 20.0, [2]float64{125, 135}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh player for each test
			player := entities.NewPlayer(64, tc.startY, testSpriteSheet)
			player.SetLevel(levelAdapter)
			
			// Set the falling velocity
			player.VelocityY = tc.velocityY
			
			// Simulate until player lands (max 100 frames)
			for i := 0; i < 100; i++ {
				prevOnGround := player.OnGround
				player.Update(1.0/60.0) // 60 FPS
				
				// Check if player just landed
				if !prevOnGround && player.OnGround {
					_, y := player.GetPosition()
					fmt.Printf("%s: Landed at Y=%.1f after %d frames\n", tc.name, y, i+1)
					
					// Check if landing position is within expected range
					if y < tc.expectedYRange[0] || y > tc.expectedYRange[1] {
						t.Errorf("%s: Landing Y position %.1f outside expected range [%.1f, %.1f]", 
							tc.name, y, tc.expectedYRange[0], tc.expectedYRange[1])
					}
					
					// Check if player is properly on ground
					if !player.OnGround {
						t.Errorf("%s: Player should be on ground after landing", tc.name)
					}
					
					// Check if velocity has been stopped
					_, vy := player.GetVelocity()
					if math.Abs(vy) > 0.1 {
						t.Errorf("%s: Player velocity should be near zero after landing, got %.2f", tc.name, vy)
					}
					
					return
				}
			}
			
			// If we reach here, player didn't land
			_, y := player.GetPosition()
			t.Errorf("%s: Player didn't land after 100 frames, final position Y=%.1f", tc.name, y)
		})
	}
}

// TestNoSinking verifies that players don't sink into platforms over time
func TestNoSinking(t *testing.T) {
	// Create a test level with a platform  
	testLevel := level.NewLevel(10, 10, 32, "Sinking Test")
	testLevel.SetTile(2, 5, level.TileSolid) // Platform at tile (2,5)
	
	// Create level adapter, sprite sheet, and player
	levelAdapter := level.NewCollisionAdapter(testLevel)
	testSpriteSheet := entities.CreateTestSpriteSheet()
	player := entities.NewPlayer(64, 128, testSpriteSheet) // Position player on platform
	player.SetLevel(levelAdapter)
	
	initialY := 128.0
	maxAllowedY := 135.0 // Allow some tolerance but shouldn't sink much
	
	// Simulate standing on platform for many frames
	for i := 0; i < 300; i++ { // 5 seconds at 60 FPS
		player.Update(1.0/60.0)
		
		_, y := player.GetPosition()
		
		// Player should not sink too far below initial position
		if y > maxAllowedY {
			t.Errorf("Player sank too far into platform at frame %d: Y=%.1f (started at %.1f, max allowed %.1f)", 
				i, y, initialY, maxAllowedY)
			break
		}
		
		// Player should remain on ground
		if !player.OnGround {
			t.Errorf("Player lost ground state at frame %d, Y=%.1f", i, y)
			break
		}
	}
	
	// Final position check
	_, finalY := player.GetPosition()
	fmt.Printf("Final position after 300 frames: Y=%.1f (started at %.1f)\n", finalY, initialY)
}
