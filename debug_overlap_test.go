package main

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

func TestDebugOverlapCalculation(t *testing.T) {
	// Create a test level with a single tile platform
	testLevel := level.NewLevel(10, 10, 32, "debug_overlap_test")
	
	// Create a single tile platform at X=3, Y=7 (position 96-128, 224-256)
	testLevel.SetTile(3, 7, level.TileSolid)
	
	// Create a 1x1 pixel image for testing
	testImage := ebiten.NewImage(1, 1)
	
	// Test different overlap scenarios
	testCases := []struct{
		x float64
		expectedOnGround bool
		overlapPercent float64
	}{
		{75, false, 34}, // 75-107 vs 96-128, overlap 96-107 = 11px = 34%
		{80, true, 50},  // 80-112 vs 96-128, overlap 96-112 = 16px = 50%
		{64, false, 0},  // 64-96 vs 96-128, overlap = 0px = 0%
		{90, true, 72},  // 90-122 vs 96-128, overlap 96-122 = 26px = 81%
	}
	
	adapter := level.NewCollisionAdapter(testLevel)
	
	for _, tc := range testCases {
		player := entities.NewPlayer(tc.x, 192, testImage)
		player.SetLevel(adapter)
		
		deltaTime := 1.0 / 60.0
		player.Update(deltaTime)
		
		x, y := player.GetPosition()
		actualOnGround := player.IsOnGround()
		
		fmt.Printf("Position (%.0f, %.0f): Expected OnGround=%v, Actual OnGround=%v, Expected Overlap=%.0f%%\n", 
			x, y, tc.expectedOnGround, actualOnGround, tc.overlapPercent)
		
		if actualOnGround != tc.expectedOnGround {
			t.Errorf("Position %.0f: Expected OnGround=%v, got %v", tc.x, tc.expectedOnGround, actualOnGround)
		}
	}
}
