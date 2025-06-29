package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler manages player input and controls
type InputHandler struct {
	player *Player
}

// NewInputHandler creates a new input handler for the player
func NewInputHandler(player *Player) *InputHandler {
	return &InputHandler{
		player: player,
	}
}

// Update processes input and updates player accordingly
func (ih *InputHandler) Update() {
	if ih.player == nil {
		return
	}
	
	// Horizontal movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ih.player.MoveLeft()
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		ih.player.MoveRight()
	}
	
	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ih.player.Jump()
	}
	
	// Climbing controls (when near climbable surfaces)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if ih.player.IsClimbing {
			ih.player.ClimbUp()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if ih.player.IsClimbing {
			ih.player.ClimbDown()
		}
	}
	
	// Debug controls (remove in final version)
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		// Toggle climbing mode for testing
		if ih.player.IsClimbing {
			ih.player.StopClimbing()
		} else {
			ih.player.StartClimbing()
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		// Test damage state
		ih.player.TakeDamage()
	}
}

// GetPlayer returns the player instance
func (ih *InputHandler) GetPlayer() *Player {
	return ih.player
}
