package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Player represents the ROBO-9 character
type Player struct {
	// Position and movement
	X, Y          float64
	VelocityX     float64
	VelocityY     float64
	Width, Height float64
	
	// Physics constants
	Speed         float64
	JumpSpeed     float64
	Gravity       float64
	Friction      float64
	
	// State
	OnGround      bool
	FacingRight   bool
	IsJumping     bool
	IsMoving      bool
	IsClimbing    bool
	IsDamaged     bool
	
	// Animation
	AnimationController *AnimationController
	
	// Timing
	DamageTimer float64
	DamageTime  float64
}

// NewPlayer creates a new ROBO-9 player instance
func NewPlayer(x, y float64, spriteSheet *ebiten.Image) *Player {
	// Assuming 32x32 pixel sprites
	frameWidth := 32
	frameHeight := 32
	
	player := &Player{
		X:           x,
		Y:           y,
		Width:       float64(frameWidth),
		Height:      float64(frameHeight),
		Speed:       120.0, // pixels per second
		JumpSpeed:   200.0,
		Gravity:     500.0,
		Friction:    0.8,
		FacingRight: true,
		DamageTime:  1.0, // 1 second of damage immunity
	}
	
	// Initialize animation controller
	player.AnimationController = NewAnimationController(spriteSheet, frameWidth, frameHeight)
	player.setupAnimations()
	
	// Start with idle animation
	player.AnimationController.SetState(AnimationIdle)
	
	return player
}

// setupAnimations configures all the player animations
func (p *Player) setupAnimations() {
	// Define animation sequences (frame start, count, timing, loop)
	// These values would need to be adjusted based on your actual sprite sheet layout
	
	// Idle animation: frames 0-3, 0.2 seconds per frame, loops
	p.AnimationController.AddAnimation(AnimationIdle, 0, 4, 0.2, true)
	
	// Walk animation: frames 4-7, 0.1 seconds per frame, loops  
	p.AnimationController.AddAnimation(AnimationWalk, 4, 4, 0.1, true)
	
	// Jump animation: frames 8-9, 0.1 seconds per frame, doesn't loop
	p.AnimationController.AddAnimation(AnimationJump, 8, 2, 0.1, false)
	
	// Fall animation: frames 10-11, 0.15 seconds per frame, loops
	p.AnimationController.AddAnimation(AnimationFall, 10, 2, 0.15, true)
	
	// Climb animation: frames 12-15, 0.15 seconds per frame, loops
	p.AnimationController.AddAnimation(AnimationClimb, 12, 4, 0.15, true)
	
	// Damage animation: frames 16-17, 0.1 seconds per frame, doesn't loop
	p.AnimationController.AddAnimation(AnimationDamage, 16, 2, 0.1, false)
}

// Update updates the player's state and animation
func (p *Player) Update(deltaTime float64) {
	// Update damage timer
	if p.DamageTimer > 0 {
		p.DamageTimer -= deltaTime
		if p.DamageTimer <= 0 {
			p.IsDamaged = false
		}
	}
	
	// Apply physics
	p.updatePhysics(deltaTime)
	
	// Update animation state based on current movement
	p.updateAnimationState()
	
	// Update animation controller
	p.AnimationController.Update(deltaTime)
}

// updatePhysics handles movement and gravity
func (p *Player) updatePhysics(deltaTime float64) {
	// Apply gravity if not on ground
	if !p.OnGround {
		p.VelocityY += p.Gravity * deltaTime
	}
	
	// Apply friction to horizontal movement
	p.VelocityX *= p.Friction
	
	// Update position
	p.X += p.VelocityX * deltaTime
	p.Y += p.VelocityY * deltaTime
	
	// Check if moving horizontally
	p.IsMoving = math.Abs(p.VelocityX) > 10.0
	
	// Simple ground collision (would be replaced with proper collision detection)
	groundY := 300.0 // Temporary ground level (updated for larger screen)
	if p.Y >= groundY {
		p.Y = groundY
		p.VelocityY = 0
		p.OnGround = true
		p.IsJumping = false
	} else {
		p.OnGround = false
	}
}

// updateAnimationState determines which animation should be playing
func (p *Player) updateAnimationState() {
	if p.IsDamaged {
		p.AnimationController.SetState(AnimationDamage)
		return
	}
	
	if p.IsClimbing {
		p.AnimationController.SetState(AnimationClimb)
		return
	}
	
	if !p.OnGround {
		if p.VelocityY < 0 {
			p.AnimationController.SetState(AnimationJump)
		} else {
			p.AnimationController.SetState(AnimationFall)
		}
		return
	}
	
	if p.IsMoving {
		p.AnimationController.SetState(AnimationWalk)
	} else {
		p.AnimationController.SetState(AnimationIdle)
	}
}

// MoveLeft makes the player move left
func (p *Player) MoveLeft() {
	if !p.IsDamaged {
		p.VelocityX = -p.Speed
		p.FacingRight = false
	}
}

// MoveRight makes the player move right
func (p *Player) MoveRight() {
	if !p.IsDamaged {
		p.VelocityX = p.Speed
		p.FacingRight = true
	}
}

// Jump makes the player jump (if on ground)
func (p *Player) Jump() {
	if p.OnGround && !p.IsDamaged {
		p.VelocityY = -p.JumpSpeed
		p.IsJumping = true
		p.OnGround = false
	}
}

// StartClimbing puts the player in climbing mode
func (p *Player) StartClimbing() {
	if !p.IsDamaged {
		p.IsClimbing = true
		p.VelocityY = 0
		p.Gravity = 0
	}
}

// StopClimbing exits climbing mode
func (p *Player) StopClimbing() {
	p.IsClimbing = false
	p.Gravity = 500.0
}

// ClimbUp makes the player climb upward
func (p *Player) ClimbUp() {
	if p.IsClimbing && !p.IsDamaged {
		p.VelocityY = -p.Speed * 0.7 // Climb slower than walking
	}
}

// ClimbDown makes the player climb downward
func (p *Player) ClimbDown() {
	if p.IsClimbing && !p.IsDamaged {
		p.VelocityY = p.Speed * 0.7
	}
}

// TakeDamage puts the player in damage state
func (p *Player) TakeDamage() {
	if !p.IsDamaged {
		p.IsDamaged = true
		p.DamageTimer = p.DamageTime
		p.VelocityX = 0 // Stop movement when damaged
	}
}

// GetBounds returns the player's collision rectangle
func (p *Player) GetBounds() (float64, float64, float64, float64) {
	return p.X, p.Y, p.Width, p.Height
}

// Draw renders the player
func (p *Player) Draw(screen *ebiten.Image) {
	currentFrame := p.AnimationController.GetCurrentFrame()
	if currentFrame == nil {
		return
	}
	
	op := &ebiten.DrawImageOptions{}
	
	// Flip sprite horizontally if facing left
	if !p.FacingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(p.Width, 0)
	}
	
	// Position the sprite
	op.GeoM.Translate(p.X, p.Y)
	
	// Add damage effect (flashing)
	if p.IsDamaged {
		// Make the sprite flash by adjusting alpha
		flashCycle := math.Sin(p.DamageTimer * 20) // Fast flashing
		if flashCycle > 0 {
			op.ColorM.Scale(1, 1, 1, 0.5) // Semi-transparent
		}
	}
	
	screen.DrawImage(currentFrame, op)
}

// GetPosition returns the player's current position
func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}

// SetPosition sets the player's position
func (p *Player) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

// GetVelocity returns the player's current velocity
func (p *Player) GetVelocity() (float64, float64) {
	return p.VelocityX, p.VelocityY
}

// IsOnGround returns whether the player is on the ground
func (p *Player) IsOnGround() bool {
	return p.OnGround
}

// IsFacingRight returns whether the player is facing right
func (p *Player) IsFacingRight() bool {
	return p.FacingRight
}

// GetAnimationState returns the current animation state
func (p *Player) GetAnimationState() AnimationState {
	return p.AnimationController.GetCurrentState()
}
