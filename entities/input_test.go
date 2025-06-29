package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewInputHandler(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := NewInputHandler(player)

	if input == nil {
		t.Fatal("NewInputHandler returned nil")
	}

	if input.player != player {
		t.Error("InputHandler should store reference to player")
	}
}

func TestInputHandler_WithNilPlayer(t *testing.T) {
	input := NewInputHandler(nil)

	if input == nil {
		t.Fatal("NewInputHandler returned nil")
	}

	// Should not panic when updating with nil player
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update with nil player panicked: %v", r)
		}
	}()

	input.Update()
}

func TestInputHandler_Update(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := NewInputHandler(player)

	// Test basic update (no keys pressed)
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update method panicked: %v", r)
		}
	}()

	input.Update()
}

func TestInputHandler_PlayerReference(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := NewInputHandler(player)

	// Test that the input handler maintains the correct player reference
	if input.player != player {
		t.Error("InputHandler should maintain reference to the correct player")
	}

	// Test that we can create multiple input handlers for different players
	player2 := NewPlayer(200, 300, img)
	input2 := NewInputHandler(player2)

	if input2.player != player2 {
		t.Error("Second InputHandler should reference the second player")
	}

	if input.player == input2.player {
		t.Error("Different InputHandlers should reference different players")
	}
}

func TestInputHandler_StructIntegrity(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := &InputHandler{player: player}

	// Test that we can read the player field
	if input.player != player {
		t.Error("Failed to set player field directly")
	}

	// Test field modification
	player2 := NewPlayer(200, 300, img)
	input.player = player2

	if input.player != player2 {
		t.Error("Failed to modify player field")
	}
}

// Mock test to verify the Update method exists and can be called
func TestInputHandler_UpdateMethodExists(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := NewInputHandler(player)

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update method panicked: %v", r)
		}
	}()

	input.Update()
}

func TestInputHandler_Integration(t *testing.T) {
	img := ebiten.NewImage(320, 320)
	player := NewPlayer(100, 200, img)
	input := NewInputHandler(player)

	// Store initial player state
	_ = player.X      // Mark as used
	_ = player.Y      // Mark as used
	_ = player.VelocityX  // Mark as used
	_ = player.VelocityY  // Mark as used

	// Update input handler (simulates one frame of input processing)
	input.Update()

	// Player state may or may not have changed depending on input,
	// but the update should not cause crashes or invalid states
	if player.X < -1000 || player.X > 2000 {
		t.Error("Player X position should remain reasonable after input update")
	}

	if player.Y < -1000 || player.Y > 2000 {
		t.Error("Player Y position should remain reasonable after input update")
	}
}
