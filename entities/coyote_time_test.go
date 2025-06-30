package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestCoyoteTime(t *testing.T) {
	// Create a mock sprite sheet for testing
	mockSpriteSheet := ebiten.NewImage(256, 32) // 8 frames of 32x32
	
	// Create player at ground level (300.0 is the hardcoded ground in updateSimplePhysics)
	player := NewPlayer(100, 300, mockSpriteSheet)
	
	// Initially player should be on ground due to the position
	deltaTime := 1.0 / 60.0 // 60 FPS
	player.Update(deltaTime)
	
	// Verify player is on ground and coyote timer is 0
	if !player.OnGround {
		t.Error("Player should be on ground")
	}
	if player.CoyoteTimer != 0 {
		t.Errorf("Coyote timer should be 0 when on ground, got %f", player.CoyoteTimer)
	}
	
	// Simulate player leaving ground (not from jumping) by moving above ground
	player.Y = 250 // Move above ground level
	player.IsJumping = false
	
	// Update to trigger coyote time
	player.Update(deltaTime)
	
	// Verify coyote timer is active
	if player.CoyoteTimer <= 0 {
		t.Errorf("Coyote timer should be active after leaving ground, got %f", player.CoyoteTimer)
	}
	if player.CoyoteTimer > player.CoyoteTime {
		t.Errorf("Coyote timer should not exceed coyote time limit, got %f, limit %f", 
			player.CoyoteTimer, player.CoyoteTime)
	}
	
	// Test that jumping during coyote time works
	initialVelocityY := player.VelocityY
	player.Jump()
	
	// Verify jump occurred
	if player.VelocityY >= initialVelocityY {
		t.Error("Player should have negative Y velocity after jumping during coyote time")
	}
	if !player.IsJumping {
		t.Error("Player should be in jumping state after coyote jump")
	}
	if player.CoyoteTimer != 0 {
		t.Error("Coyote timer should be consumed after jumping")
	}
}

func TestCoyoteTimeExpiry(t *testing.T) {
	// Create a mock sprite sheet for testing
	mockSpriteSheet := ebiten.NewImage(256, 32)
	
	// Create player at ground level
	player := NewPlayer(100, 300, mockSpriteSheet)
	
	deltaTime := 1.0 / 60.0 // 60 FPS
	
	// Update to establish ground state
	player.Update(deltaTime)
	
	// Verify player is on ground
	if !player.OnGround {
		t.Error("Player should be on ground initially")
	}
	
	// Simulate player leaving ground by moving above ground level
	player.Y = 250
	player.IsJumping = false
	
	// Update to start coyote time
	player.Update(deltaTime)
	
	// Verify coyote timer started
	if player.CoyoteTimer <= 0 {
		t.Error("Coyote timer should be active")
	}
	
	// Simulate time passing beyond coyote time limit
	// Run updates for longer than coyote time duration
	totalTime := 0.0
	for totalTime < player.CoyoteTime + 0.05 { // Add small buffer
		player.Update(deltaTime)
		totalTime += deltaTime
	}
	
	// Verify coyote timer expired
	if player.CoyoteTimer > 0 {
		t.Errorf("Coyote timer should have expired, got %f", player.CoyoteTimer)
	}
	
	// Test that jumping after coyote time expiry doesn't work
	initialVelocityY := player.VelocityY
	player.Jump()
	
	// Verify jump did not occur
	if player.VelocityY != initialVelocityY {
		t.Error("Player should not be able to jump after coyote time expires")
	}
	if player.IsJumping {
		t.Error("Player should not be in jumping state after failed coyote jump")
	}
}

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
