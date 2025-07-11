package entities

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Constants for collision detection tuning
const (
	// BinarySearchMaxIterations defines the maximum number of iterations
	// for binary search collision detection. Higher values provide more
	// precision but cost more CPU cycles.
	BinarySearchMaxIterations = 50

	// BinarySearchMaxIterationsY defines the maximum number of iterations
	// for Y-axis binary search collision detection. Uses fewer iterations
	// to prevent over-precision in vertical movement.
	BinarySearchMaxIterationsY = 30

	// BinarySearchTolerance defines the minimum distance threshold for
	// binary search convergence. Smaller values provide higher precision.
	BinarySearchTolerance = 0.01

	// BinarySearchToleranceY defines the tolerance for Y-axis binary search.
	// Slightly larger than X-axis to balance precision and stability.
	BinarySearchToleranceY = 0.1
)

// Player represents the ROBO-9 character
type Player struct {
	// Position and movement
	X, Y          float64
	VelocityX     float64
	VelocityY     float64
	Width, Height float64

	// Physics constants
	Speed     float64
	JumpSpeed float64
	Gravity   float64
	Friction  float64

	// State
	OnGround    bool
	FacingRight bool
	IsJumping   bool
	IsMoving    bool
	IsClimbing  bool
	IsDamaged   bool

	// Animation
	AnimationController *AnimationController

	// Timing
	DamageTimer float64
	DamageTime  float64

	// Coyote time for forgiving jumps
	CoyoteTime  float64 // Duration of coyote time window
	CoyoteTimer float64 // Current coyote time remaining
	WasOnGroundPhysics bool // Previous frame physics ground state

	// Collision
	level CollisionChecker
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
		CoyoteTime:  0.1, // 100ms of coyote time (standard for platform edge jumps)
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

	// Update coyote timer (after physics so we have correct OnGround state)
	p.updateCoyoteTime(deltaTime)

	// Update animation state based on current movement
	p.updateAnimationState()

	// Update animation controller
	p.AnimationController.Update(deltaTime)
}

// updateCoyoteTime manages the coyote time system for forgiving jumps
func (p *Player) updateCoyoteTime(deltaTime float64) {
	// Simple coyote time logic: activate when leaving ground, count down over time
	
	// If player was on ground last frame but isn't now (and not from jumping), start coyote time
	if p.WasOnGroundPhysics && !p.OnGround && !p.IsJumping {
		if p.CoyoteTimer <= 0 { // Only start if not already active
			p.CoyoteTimer = p.CoyoteTime
		}
	}

	// Count down coyote timer
	if p.CoyoteTimer > 0 {
		p.CoyoteTimer -= deltaTime
		if p.CoyoteTimer < 0 {
			p.CoyoteTimer = 0
		}
	}

	// Reset coyote timer when landing on ground
	if p.OnGround && p.VelocityY >= 0 {
		p.CoyoteTimer = 0
	}

	// Store current state for next frame comparison
	p.WasOnGroundPhysics = p.OnGround
}


// updatePhysics handles movement and gravity with tile-based collision
func (p *Player) updatePhysics(deltaTime float64) {
	// Store previous position for collision resolution
	prevX := p.X
	prevY := p.Y

	// Apply gravity if not on ground
	if !p.OnGround {
		p.VelocityY += p.Gravity * deltaTime
	}

	// Apply friction to horizontal movement
	p.VelocityX *= p.Friction

	// Calculate intended movement
	deltaX := p.VelocityX * deltaTime
	deltaY := p.VelocityY * deltaTime

	// Check if moving horizontally
	p.IsMoving = math.Abs(p.VelocityX) > 10.0

	// If no level is set, use simple ground collision as fallback
	if p.level == nil {
		p.updateSimplePhysics(deltaTime, deltaX, deltaY)
		return
	}

	// Store previous ground state for smoother transitions
	prevOnGround := p.OnGround

	// Reset ground state
	p.OnGround = false

	// Use swept collision detection for more robust movement
	finalX, finalY := p.performSweptMovement(prevX, prevY, deltaX, deltaY)
	p.X = finalX
	p.Y = finalY

	// Final collision check to set ground state and handle any remaining issues
	result := p.level.CheckCollision(p.X, p.Y, p.Width, p.Height)

	// Update climbing state based on collision
	if result.ClimbableSurface && !p.IsDamaged {
		// Player can potentially climb here
		// The climbing mode is still controlled by input (C key for debug)
	}

	// Handle ground state (swept movement should have already set OnGround for most cases)
	if result.OnGround && !p.OnGround {
		p.OnGround = true
		p.IsJumping = false
		if p.VelocityY > 0 {
			p.VelocityY = 0
		}
	}

	// Use hysteresis for ground state to reduce jitter, but only for very small movements
	if !result.OnGround && !result.CollisionY {
		// Only apply hysteresis if the player is essentially stationary both horizontally and vertically
		// This prevents false ground detection when walking off platforms
		if prevOnGround && math.Abs(p.VelocityY) < 0.5 && math.Abs(p.VelocityX) < 5.0 {
			// Much stricter tolerance to prevent interference with edge walking
			p.OnGround = true
		} else {
			p.OnGround = false
		}
	}

	// Handle dangerous tiles
	if result.DangerousTile && !p.IsDamaged {
		p.TakeDamage()
	}

	// Update wall touching state for future wall jumping/climbing features
	if result.TouchingWall {
		// Could be used for wall jumping in the future
	}
}

// handleCollisionResult processes collision results and updates player state
func (p *Player) handleCollisionResult(result *CollisionResult, prevX, prevY float64) bool {
	if result == nil || !result.Collided {
		p.OnGround = false
		return false
	}

	// Handle horizontal collision
	if result.CollisionX {
		return true
	}

	// Handle vertical collision
	if result.CollisionY {
		return true
	}

	// Set ground state
	if result.OnGround {
		p.OnGround = true
		p.IsJumping = false
	}

	return false
}

// updateSimplePhysics provides fallback physics when no level is set
func (p *Player) updateSimplePhysics(deltaTime, deltaX, deltaY float64) {
	// Update position
	p.X += deltaX
	p.Y += deltaY

	// Simple ground collision (fallback)
	groundY := 300.0 // Temporary ground level
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

// Jump makes the player jump (if on ground or during coyote time)
func (p *Player) Jump() {
	// Allow jumping if on ground or during coyote time window
	canJump := (p.OnGround || p.CoyoteTimer > 0) && !p.IsDamaged

	if canJump {
		p.VelocityY = -p.JumpSpeed
		p.IsJumping = true
		p.OnGround = false
		p.CoyoteTimer = 0 // Consume coyote time
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

// SetLevel sets the level for collision detection
func (p *Player) SetLevel(level CollisionChecker) {
	p.level = level
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

// GetCoyoteTimer returns the current coyote time remaining
func (p *Player) GetCoyoteTimer() float64 {
	return p.CoyoteTimer
}

// IsOnGround returns whether the player is currently on ground
func (p *Player) IsOnGround() bool {
	return p.OnGround
}

// GetVelocityY returns the player's Y velocity
func (p *Player) GetVelocityY() float64 {
	return p.VelocityY
}

// IsFacingRight returns whether the player is facing right
func (p *Player) IsFacingRight() bool {
	return p.FacingRight
}

// GetAnimationState returns the current animation state
func (p *Player) GetAnimationState() AnimationState {
	return p.AnimationController.GetCurrentState()
}

// performSweptMovement performs swept collision detection to prevent tunneling
// and ensure accurate collision response regardless of movement speed
func (p *Player) performSweptMovement(startX, startY, deltaX, deltaY float64) (float64, float64) {
	if p.level == nil {
		return startX + deltaX, startY + deltaY
	}

	finalX := startX
	finalY := startY

	// Step 1: Handle horizontal movement with swept collision
	if deltaX != 0 {
		finalX = p.sweptHorizontalMovement(startX, startY, deltaX)
	}

	// Step 2: Handle vertical movement with swept collision
	if deltaY != 0 {
		finalY = p.sweptVerticalMovement(finalX, startY, deltaY)
	}

	return finalX, finalY
}

// sweptHorizontalMovement handles horizontal movement with collision detection
func (p *Player) sweptHorizontalMovement(startX, y, deltaX float64) float64 {
	targetX := startX + deltaX

	// Check if the target position would cause collision
	result := p.level.CheckCollision(targetX, y, p.Width, p.Height)
	if !result.CollisionX {
		// No collision, move to target position
		return targetX
	}

	// Use binary search for precise collision point detection
	return p.binarySearchCollisionX(startX, y, deltaX)
}

// binarySearchCollisionX uses binary search to find the exact collision point on X axis
func (p *Player) binarySearchCollisionX(startX, y, deltaX float64) float64 {
	left := startX
	right := startX + deltaX

	// Ensure left is valid, right causes collision
	if deltaX > 0 {
		// Moving right
		for i := 0; i < BinarySearchMaxIterations; i++ {
			mid := (left + right) / 2

			if math.Abs(right-left) < BinarySearchTolerance {
				p.VelocityX = 0
				return left
			}

			result := p.level.CheckCollision(mid, y, p.Width, p.Height)
			if result.CollisionX {
				right = mid
			} else {
				left = mid
			}
		}
	} else {
		// Moving left
		for i := 0; i < BinarySearchMaxIterations; i++ {
			mid := (left + right) / 2

			if math.Abs(right-left) < BinarySearchTolerance {
				p.VelocityX = 0
				return right
			}

			result := p.level.CheckCollision(mid, y, p.Width, p.Height)
			if result.CollisionX {
				left = mid
			} else {
				right = mid
			}
		}
	}

	p.VelocityX = 0
	if deltaX > 0 {
		return left
	} else {
		return right
	}
}

// sweptVerticalMovement handles vertical movement with collision detection
func (p *Player) sweptVerticalMovement(x, startY, deltaY float64) float64 {
	targetY := startY + deltaY

	// Check if the target position would cause collision
	result := p.level.CheckCollision(x, targetY, p.Width, p.Height)
	if !result.CollisionY && !result.OnGround {
		// No collision, move to target position
		return targetY
	}

	// Use binary search for precise collision point detection
	return p.binarySearchCollisionY(x, startY, deltaY)
}

// binarySearchCollisionY uses binary search to find the exact collision point on Y axis
func (p *Player) binarySearchCollisionY(x, startY, deltaY float64) float64 {
	if deltaY > 0 { // Moving down (falling)
		left := startY
		right := startY + deltaY
		bestValidY := startY

		for i := 0; i < BinarySearchMaxIterationsY; i++ {
			if math.Abs(right-left) < BinarySearchToleranceY {
				break
			}

			mid := (left + right) / 2
			result := p.level.CheckCollision(x, mid, p.Width, p.Height)

			if result.CollisionY || result.OnGround {
				right = mid
			} else {
				left = mid
				bestValidY = mid
			}
		}

		// Final position check and state setting
		finalY := bestValidY

		// Check for ground detection in a small range around the final position
		for offset := 0.0; offset <= 1.0; offset += 0.1 {
			testY := finalY + offset
			testResult := p.level.CheckCollision(x, testY, p.Width, p.Height)
			if testResult.OnGround || testResult.CollisionY {
				if testResult.OnGround {
					p.OnGround = true
					p.IsJumping = false
				}
				p.VelocityY = 0
				return testY
			}
		}

		return finalY

	} else { // Moving up (jumping)
		left := startY + deltaY
		right := startY

		for i := 0; i < BinarySearchMaxIterationsY; i++ {
			if math.Abs(right-left) < BinarySearchToleranceY {
				break
			}

			mid := (left + right) / 2
			result := p.level.CheckCollision(x, mid, p.Width, p.Height)

			if result.CollisionY {
				left = mid
			} else {
				right = mid
			}
		}

		p.VelocityY = 0
		return right
	}
}

// GetDebugInfo returns debug information about the player's state
func (p *Player) GetDebugInfo() string {
	onGroundStr := "false"
	if p.OnGround {
		onGroundStr = "true"
	}
	
	wasOnGroundStr := "false"
	if p.WasOnGroundPhysics {
		wasOnGroundStr = "true"
	}

	return fmt.Sprintf("OnGround: %s, WasOnGround: %s, CoyoteTimer: %.3f, VelocityY: %.1f",
		onGroundStr, wasOnGroundStr, p.CoyoteTimer, p.VelocityY)
}
