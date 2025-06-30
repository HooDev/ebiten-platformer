package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// Note: Basic coyote time functionality is tested in the root-level coyote_time_test.go
// with realistic level collision. These tests focus on edge cases and unit-level behavior.

func TestCoyoteTimeResetOnGround(t *testing.T) {
	// Create a mock sprite sheet for testing
	mockSpriteSheet := ebiten.NewImage(256, 32)
	
	// Create player above ground initially
	player := NewPlayer(100, 250, mockSpriteSheet)
	
	deltaTime := 1.0 / 60.0
	
	// Update to establish not-on-ground state
	player.Update(deltaTime)
	
	// Manually set coyote timer as if it was active
	player.CoyoteTimer = player.CoyoteTime / 2 // Half expired
	
	// Simulate landing back on ground by moving to ground level
	player.Y = 300
	player.Update(deltaTime)
	
	// Verify coyote timer is reset
	if player.CoyoteTimer != 0 {
		t.Errorf("Coyote timer should be reset when landing on ground, got %f", player.CoyoteTimer)
	}
}

func TestCoyoteTimeNotActivatedByJump(t *testing.T) {
	// Create a mock sprite sheet for testing
	mockSpriteSheet := ebiten.NewImage(256, 32)
	
	// Create player at ground level
	player := NewPlayer(100, 300, mockSpriteSheet)
	
	deltaTime := 1.0 / 60.0
	
	// Update to establish ground state
	player.Update(deltaTime)
	
	// Verify player is on ground
	if !player.OnGround {
		t.Error("Player should be on ground initially")
	}
	
	// Simulate jumping (which should set IsJumping = true)
	player.Jump()
	
	// Verify jump was successful
	if !player.IsJumping {
		t.Error("Player should be jumping after Jump() call")
	}
	
	// Update physics which will move player and update ground state
	player.Update(deltaTime)
	
	// Verify coyote timer is NOT activated when leaving ground due to jumping
	if player.CoyoteTimer > 0 {
		t.Errorf("Coyote timer should not activate when leaving ground due to jumping, got %f", player.CoyoteTimer)
	}
}

func TestCoyoteTimeDuration(t *testing.T) {
	// Create a mock sprite sheet for testing
	mockSpriteSheet := ebiten.NewImage(256, 32)
	
	// Create player
	player := NewPlayer(100, 100, mockSpriteSheet)
	
	// Verify coyote time is set to expected duration (100ms)
	expectedCoyoteTime := 0.1
	if player.CoyoteTime != expectedCoyoteTime {
		t.Errorf("Expected coyote time to be %f seconds, got %f", expectedCoyoteTime, player.CoyoteTime)
	}
}

// Benchmark the coyote time update performance
func BenchmarkCoyoteTimeUpdate(b *testing.B) {
	mockSpriteSheet := ebiten.NewImage(256, 32)
	player := NewPlayer(100, 100, mockSpriteSheet)
	
	deltaTime := 1.0 / 60.0
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		player.updateCoyoteTime(deltaTime)
	}
}
