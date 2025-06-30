package main

import (
	"testing"
	"ebiten-platformer/level"
)

func TestPositionRange(t *testing.T) {
	// Test a range of positions around the platform edge
	testLevel := level.NewLevel(10, 10, 32, "Test Level")
	
	for x := 2; x <= 5; x++ {
		testLevel.SetTile(x, 5, level.TileSolid)
	}
	
	adapter := level.NewCollisionAdapter(testLevel)
	
	// Test positions from X=180 to X=210 in steps of 2
	for x := 180.0; x <= 210.0; x += 2.0 {
		result := adapter.CheckCollision(x, 128, 32, 32)
		
		// Calculate expected overlap
		playerLeft := x
		playerRight := x + 32
		platformLeft := 64.0
		platformRight := 192.0
		
		overlapLeft := maxVal(playerLeft, platformLeft)
		overlapRight := minVal(playerRight, platformRight) 
		overlapWidth := maxVal(0, overlapRight - overlapLeft)
		overlapPercent := (overlapWidth / 32) * 100
		
		t.Logf("X=%.0f: OnGround=%t, Overlap=%.1f pixels (%.1f%%)", 
			x, result.OnGround, overlapWidth, overlapPercent)
	}
}

func maxVal(a, b float64) float64 {
	if a > b { return a }
	return b
}

func minVal(a, b float64) float64 {
	if a < b { return a }
	return b
}
