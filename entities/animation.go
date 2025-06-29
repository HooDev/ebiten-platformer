package entities

import (
	"image"
	"github.com/hajimehoshi/ebiten/v2"
)

// AnimationState represents different animation states
type AnimationState int

const (
	AnimationIdle AnimationState = iota
	AnimationWalk
	AnimationJump
	AnimationFall
	AnimationClimb
	AnimationDamage
)

// Animation represents a single animation sequence
type Animation struct {
	Frames      []*ebiten.Image // Individual frames of the animation
	FrameCount  int             // Number of frames in the animation
	FrameTime   float64         // Time per frame in seconds
	Loop        bool            // Whether the animation should loop
	CurrentTime float64         // Current time in the animation
	Finished    bool            // Whether the animation has finished (for non-looping)
}

// NewAnimation creates a new animation from a sprite sheet
func NewAnimation(spriteSheet *ebiten.Image, frameWidth, frameHeight, frameCount int, frameTime float64, loop bool) *Animation {
	frames := make([]*ebiten.Image, frameCount)
	
	// Extract frames from sprite sheet (assumes horizontal layout)
	for i := 0; i < frameCount; i++ {
		x := i * frameWidth
		y := 0
		
		// Create a sub-image for each frame
		frame := spriteSheet.SubImage(image.Rect(x, y, x+frameWidth, y+frameHeight)).(*ebiten.Image)
		
		frames[i] = frame
	}
	
	return &Animation{
		Frames:     frames,
		FrameCount: frameCount,
		FrameTime:  frameTime,
		Loop:       loop,
	}
}

// Update updates the animation timing
func (a *Animation) Update(deltaTime float64) {
	if a.Finished && !a.Loop {
		return
	}
	
	a.CurrentTime += deltaTime
	
	// Check if we've completed the animation
	totalAnimationTime := float64(a.FrameCount) * a.FrameTime
	if a.CurrentTime >= totalAnimationTime {
		if a.Loop {
			// Reset for looping animations
			a.CurrentTime = 0
		} else {
			// Mark as finished for non-looping animations
			a.CurrentTime = totalAnimationTime - 0.001 // Keep at last frame
			a.Finished = true
		}
	}
}

// GetCurrentFrame returns the current frame image
func (a *Animation) GetCurrentFrame() *ebiten.Image {
	if a.FrameCount == 0 {
		return nil
	}
	
	frameIndex := int(a.CurrentTime / a.FrameTime)
	if frameIndex >= a.FrameCount {
		frameIndex = a.FrameCount - 1
	}
	
	return a.Frames[frameIndex]
}

// Reset resets the animation to the beginning
func (a *Animation) Reset() {
	a.CurrentTime = 0
	a.Finished = false
}

// IsFinished returns whether the animation has finished (for non-looping animations)
func (a *Animation) IsFinished() bool {
	return a.Finished
}

// AnimationController manages multiple animations for a single entity
type AnimationController struct {
	animations    map[AnimationState]*Animation
	currentState  AnimationState
	previousState AnimationState
	spriteSheet   *ebiten.Image
	frameWidth    int
	frameHeight   int
}

// NewAnimationController creates a new animation controller
func NewAnimationController(spriteSheet *ebiten.Image, frameWidth, frameHeight int) *AnimationController {
	return &AnimationController{
		animations:  make(map[AnimationState]*Animation),
		spriteSheet: spriteSheet,
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
	}
}

// AddAnimation adds an animation to the controller
func (ac *AnimationController) AddAnimation(state AnimationState, startFrame, frameCount int, frameTime float64, loop bool) {
	// Extract the specific frames for this animation from the sprite sheet
	frames := make([]*ebiten.Image, frameCount)
	
	// Calculate sprite sheet layout
	sheetWidth := ac.spriteSheet.Bounds().Dx()
	sheetHeight := ac.spriteSheet.Bounds().Dy()
	
	// Check if the sprite sheet is large enough for our frame size
	if sheetWidth < ac.frameWidth || sheetHeight < ac.frameHeight {
		// For test cases or invalid sprite sheets, create dummy frames
		for i := 0; i < frameCount; i++ {
			// Use the entire sprite sheet as a single frame (for 1x1 test images)
			frames[i] = ac.spriteSheet
		}
	} else {
		// Calculate frames per row
		framesPerRow := sheetWidth / ac.frameWidth
		if framesPerRow == 0 {
			framesPerRow = 1 // Prevent division by zero
		}
		
		for i := 0; i < frameCount; i++ {
			frameIndex := startFrame + i
			x := (frameIndex % framesPerRow) * ac.frameWidth
			y := (frameIndex / framesPerRow) * ac.frameHeight
			
			// Ensure we don't go out of bounds
			if x+ac.frameWidth > sheetWidth {
				x = sheetWidth - ac.frameWidth
				if x < 0 {
					x = 0
				}
			}
			if y+ac.frameHeight > sheetHeight {
				y = sheetHeight - ac.frameHeight
				if y < 0 {
					y = 0
				}
			}
			
			frame := ac.spriteSheet.SubImage(image.Rect(x, y, x+ac.frameWidth, y+ac.frameHeight)).(*ebiten.Image)
			frames[i] = frame
		}
	}
	
	ac.animations[state] = &Animation{
		Frames:     frames,
		FrameCount: frameCount,
		FrameTime:  frameTime,
		Loop:       loop,
	}
}

// SetState changes the current animation state
func (ac *AnimationController) SetState(state AnimationState) {
	if state != ac.currentState {
		ac.previousState = ac.currentState
		ac.currentState = state
		
		// Reset the new animation
		if animation, exists := ac.animations[state]; exists {
			animation.Reset()
		}
	}
}

// GetCurrentState returns the current animation state
func (ac *AnimationController) GetCurrentState() AnimationState {
	return ac.currentState
}

// Update updates the current animation
func (ac *AnimationController) Update(deltaTime float64) {
	if animation, exists := ac.animations[ac.currentState]; exists {
		animation.Update(deltaTime)
	}
}

// GetCurrentFrame returns the current frame of the active animation
func (ac *AnimationController) GetCurrentFrame() *ebiten.Image {
	if animation, exists := ac.animations[ac.currentState]; exists {
		return animation.GetCurrentFrame()
	}
	return nil
}

// IsCurrentAnimationFinished returns whether the current animation has finished
func (ac *AnimationController) IsCurrentAnimationFinished() bool {
	if animation, exists := ac.animations[ac.currentState]; exists {
		return animation.IsFinished()
	}
	return false
}
