package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewPlayer(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	if player == nil {
		t.Fatal("NewPlayer returned nil")
	}

	if player.X != 100 {
		t.Errorf("Expected X position 100, got %f", player.X)
	}

	if player.Y != 200 {
		t.Errorf("Expected Y position 200, got %f", player.Y)
	}

	if player.VelocityX != 0 {
		t.Errorf("Expected VelocityX 0, got %f", player.VelocityX)
	}

	if player.VelocityY != 0 {
		t.Errorf("Expected VelocityY 0, got %f", player.VelocityY)
	}

	if player.FacingRight != true {
		t.Error("Player should start facing right")
	}

	if player.Width != 32 {
		t.Errorf("Expected Width 32, got %f", player.Width)
	}

	if player.Height != 32 {
		t.Errorf("Expected Height 32, got %f", player.Height)
	}

	if player.AnimationController == nil {
		t.Error("AnimationController should be initialized")
	}

	if player.GetAnimationState() != AnimationIdle {
		t.Errorf("Expected initial state %v, got %v", AnimationIdle, player.GetAnimationState())
	}
}

func TestPlayer_Update(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	deltaTime := 1.0 / 60.0 // 60 FPS

	// Test basic update without input
	player.Update(deltaTime)

	// Player should apply gravity if not on ground
	if !player.OnGround && player.VelocityY <= 0 {
		t.Error("Player should have positive Y velocity due to gravity when not on ground")
	}
}

func TestPlayer_Movement(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test left movement
	player.MoveLeft()
	if player.VelocityX >= 0 {
		t.Error("Player should have negative velocity when moving left")
	}
	if player.FacingRight != false {
		t.Error("Player should be facing left when moving left")
	}

	// Test right movement
	player.MoveRight()
	if player.VelocityX <= 0 {
		t.Error("Player should have positive velocity when moving right")
	}
	if player.FacingRight != true {
		t.Error("Player should be facing right when moving right")
	}
}

func TestPlayer_Jump(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test jump from ground
	player.OnGround = true
	player.Jump()

	if player.VelocityY >= 0 {
		t.Error("Player should have negative Y velocity when jumping")
	}

	if player.IsJumping != true {
		t.Error("Player should be in jumping state")
	}

	// Test that player can't jump when not on ground
	player.OnGround = false
	prevVelY := player.VelocityY
	player.Jump()

	if player.VelocityY != prevVelY {
		t.Error("Player should not be able to jump when not on ground")
	}
}

func TestPlayer_Climbing(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test start climbing
	player.StartClimbing()
	if !player.IsClimbing {
		t.Error("Player should be in climbing state")
	}

	// Test climb up
	player.ClimbUp()
	if player.VelocityY >= 0 {
		t.Error("Player should have negative Y velocity when climbing up")
	}

	// Test climb down
	player.ClimbDown()
	if player.VelocityY <= 0 {
		t.Error("Player should have positive Y velocity when climbing down")
	}

	// Test stop climbing
	player.StopClimbing()
	if player.IsClimbing {
		t.Error("Player should not be in climbing state after stopping")
	}
}

func TestPlayer_Damage(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test taking damage
	player.TakeDamage()
	if !player.IsDamaged {
		t.Error("Player should be in damaged state")
	}

	if player.DamageTimer <= 0 {
		t.Error("Damage timer should be set")
	}

	// Test damage timer countdown
	deltaTime := 0.5
	player.Update(deltaTime)
	if player.DamageTimer <= 0 {
		t.Error("Damage timer should still be counting down")
	}

	// Test damage state clearing
	player.Update(1.0) // Should clear damage state
	if player.IsDamaged {
		t.Error("Player should no longer be damaged")
	}
}

func TestPlayer_GettersAndSetters(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test position getter
	x, y := player.GetPosition()
	if x != 100 || y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", x, y)
	}

	// Test position setter
	player.SetPosition(150, 250)
	x, y = player.GetPosition()
	if x != 150 || y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", x, y)
	}

	// Test velocity getter
	velX, velY := player.GetVelocity()
	if velX != player.VelocityX || velY != player.VelocityY {
		t.Errorf("Velocity getter returned incorrect values")
	}

	// Test facing direction
	player.FacingRight = true
	if !player.IsFacingRight() {
		t.Error("IsFacingRight should return true")
	}

	player.FacingRight = false
	if player.IsFacingRight() {
		t.Error("IsFacingRight should return false")
	}

	// Test on ground state
	player.OnGround = true
	if !player.IsOnGround() {
		t.Error("IsOnGround should return true")
	}

	player.OnGround = false
	if player.IsOnGround() {
		t.Error("IsOnGround should return false")
	}
}

func TestPlayer_GetBounds(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	x, y, width, height := player.GetBounds()

	expectedX := player.X
	expectedY := player.Y
	expectedWidth := player.Width
	expectedHeight := player.Height

	if x != expectedX || y != expectedY || width != expectedWidth || height != expectedHeight {
		t.Errorf("GetBounds returned incorrect values: got (%f, %f, %f, %f), expected (%f, %f, %f, %f)",
			x, y, width, height, expectedX, expectedY, expectedWidth, expectedHeight)
	}
}

func TestPlayer_AnimationStates(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	// Test initial state
	if player.GetAnimationState() != AnimationIdle {
		t.Errorf("Expected initial animation state %v, got %v", AnimationIdle, player.GetAnimationState())
	}

	// Test that animation controller is properly initialized
	if player.AnimationController == nil {
		t.Fatal("Animation controller should be initialized")
	}

	// Update animation state based on movement
	player.IsMoving = true
	player.updateAnimationState()
	// Note: The actual animation state changes depend on the implementation
	// This test just ensures the method doesn't panic
}

func TestPlayer_Physics(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)

	deltaTime := 1.0 / 60.0 // 60 FPS

	// Test gravity application
	player.OnGround = false
	initialVelY := player.VelocityY
	player.updatePhysics(deltaTime)

	if player.VelocityY <= initialVelY {
		t.Error("Gravity should increase Y velocity when not on ground")
	}

	// Test friction application when on ground
	player.OnGround = true
	player.VelocityX = 100.0 // Set some X velocity
	initialVelX := player.VelocityX
	player.updatePhysics(deltaTime)

	if player.VelocityX >= initialVelX {
		t.Error("Friction should reduce X velocity when on ground")
	}
}
