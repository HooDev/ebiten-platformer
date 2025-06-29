# Animation System Documentation

## Overview

The ROBO-9 Platformer uses a comprehensive animation system built on top of Ebitengine to handle sprite-based character animations. This system provides smooth, timing-based animations with support for multiple animation states and automatic transitions.

## Architecture

### Core Components

#### 1. Animation (`entities/animation.go`)
Represents a single animation sequence with multiple frames.

```go
type Animation struct {
    Frames      []*ebiten.Image // Individual frames
    FrameCount  int             // Total frames
    FrameTime   float64         // Time per frame (seconds)
    Loop        bool            // Whether to loop
    CurrentTime float64         // Current position in animation
    Finished    bool            // Completion status
}
```

#### 2. AnimationController (`entities/animation.go`)
Manages multiple animations for a single entity and handles state transitions.

```go
type AnimationController struct {
    animations    map[AnimationState]*Animation
    currentState  AnimationState
    previousState AnimationState
    spriteSheet   *ebiten.Image
    frameWidth    int
    frameHeight   int
}
```

#### 3. AnimationState (`entities/animation.go`)
Enumeration defining available animation states.

```go
const (
    AnimationIdle AnimationState = iota
    AnimationWalk
    AnimationJump
    AnimationFall
    AnimationClimb
    AnimationDamage
)
```

## Usage Guide

### Creating an Animation Controller

```go
// Create controller with sprite sheet and frame dimensions
controller := NewAnimationController(spriteSheet, 32, 32)

// Add animations with: state, startFrame, frameCount, frameTime, loop
controller.AddAnimation(AnimationIdle, 0, 4, 0.2, true)
controller.AddAnimation(AnimationWalk, 4, 4, 0.1, true)
controller.AddAnimation(AnimationJump, 8, 2, 0.1, false)
```

### Managing Animation States

```go
// Set animation state (automatically resets animation)
controller.SetState(AnimationWalk)

// Get current state
currentState := controller.GetCurrentState()

// Check if animation is finished (for non-looping animations)
if controller.IsCurrentAnimationFinished() {
    // Handle animation completion
}
```

### Updating and Rendering

```go
// Update animation (called each frame)
controller.Update(deltaTime)

// Get current frame for rendering
currentFrame := controller.GetCurrentFrame()
if currentFrame != nil {
    screen.DrawImage(currentFrame, drawOptions)
}
```

## Animation States

### State Descriptions

1. **AnimationIdle**
   - Triggered when: Player is not moving and on ground
   - Loop: Yes
   - Purpose: Resting state with subtle movement

2. **AnimationWalk**
   - Triggered when: Player is moving horizontally on ground
   - Loop: Yes
   - Purpose: Walking/running movement

3. **AnimationJump**
   - Triggered when: Player jumps (negative Y velocity)
   - Loop: No
   - Purpose: Launch animation for jumping

4. **AnimationFall**
   - Triggered when: Player is airborne with positive Y velocity
   - Loop: Yes
   - Purpose: Falling through air

5. **AnimationClimb**
   - Triggered when: Player is in climbing mode
   - Loop: Yes
   - Purpose: Wall climbing movement

6. **AnimationDamage**
   - Triggered when: Player takes damage
   - Loop: No
   - Purpose: Damage reaction and recovery

### State Transition Logic

The animation state is automatically determined based on player physics and status:

```go
func (p *Player) updateAnimationState() {
    // Priority 1: Damage state (overrides all others)
    if p.IsDamaged {
        p.AnimationController.SetState(AnimationDamage)
        return
    }
    
    // Priority 2: Climbing state
    if p.IsClimbing {
        p.AnimationController.SetState(AnimationClimb)
        return
    }
    
    // Priority 3: Airborne states
    if !p.OnGround {
        if p.VelocityY < 0 {
            p.AnimationController.SetState(AnimationJump)
        } else {
            p.AnimationController.SetState(AnimationFall)
        }
        return
    }
    
    // Priority 4: Ground-based movement
    if p.IsMoving {
        p.AnimationController.SetState(AnimationWalk)
    } else {
        p.AnimationController.SetState(AnimationIdle)
    }
}
```

## Technical Implementation

### Frame Extraction

The system automatically extracts frames from a sprite sheet using a grid-based approach:

```go
// Calculate frame position in sprite sheet
frameIndex := startFrame + i
x := (frameIndex % (spriteSheet.Width() / frameWidth)) * frameWidth
y := (frameIndex / (spriteSheet.Width() / frameWidth)) * frameHeight

// Extract sub-image
frame := spriteSheet.SubImage(image.Rect(x, y, x+frameWidth, y+frameHeight))
```

### Timing System

- **Delta Time Based**: Animations use real-time progression rather than frame counting
- **Smooth Playback**: Maintains consistent speed regardless of framerate
- **Precision**: Uses floating-point timing for smooth transitions

### Memory Management

- **Efficient Storage**: Frames are stored as sub-images, not duplicated pixels
- **Minimal Overhead**: Animation state uses minimal memory
- **GPU Optimized**: Leverages Ebitengine's efficient rendering

## Performance Characteristics

### CPU Usage
- **Low Impact**: Animation updates are O(1) operations
- **Efficient State Transitions**: Only resets when state changes
- **Minimal Allocations**: Reuses existing frame references

### Memory Usage
- **Shared Frames**: Multiple animations can reference same frames
- **Sub-Image Efficiency**: No pixel data duplication
- **Controller Overhead**: ~200 bytes per animation controller

### GPU Usage
- **Batch Friendly**: Compatible with Ebitengine's draw call batching
- **Texture Efficient**: Single sprite sheet minimizes texture switches

## Extension Points

### Adding New Animation States

1. **Define State**: Add new constant to `AnimationState` enum
2. **Add Animation**: Call `AddAnimation()` with frame data
3. **Update Logic**: Modify state transition logic if needed

```go
// Add new state
const (
    // ... existing states ...
    AnimationSpecialMove
)

// Configure animation
controller.AddAnimation(AnimationSpecialMove, 18, 6, 0.08, false)

// Add transition logic
if player.IsPerformingSpecialMove {
    controller.SetState(AnimationSpecialMove)
}
```

### Custom Animation Timing

```go
// Variable frame timing (future enhancement)
type Animation struct {
    // ... existing fields ...
    FrameTimes []float64  // Individual frame durations
}
```

### Animation Events

```go
// Callback system (future enhancement)
type Animation struct {
    // ... existing fields ...
    OnComplete func()           // Called when animation finishes
    OnFrame    func(int)        // Called on each frame change
}
```

## Debugging and Testing

### Debug Information

The player entity provides debug information for animation states:

```go
// Get current animation state for debugging
animState := player.GetAnimationState()
fmt.Printf("Current animation: %v", animState)
```

### Test Sprite Sheet

The system includes a test sprite sheet generator for development:

```go
testSheet := entities.CreateTestSpriteSheet()
// Creates colored rectangles for each animation state
```

### Visual Debugging

Enable debug overlays to visualize:
- Current animation state
- Frame timing
- State transition triggers

## Best Practices

### Animation Design
1. **Consistent Timing**: Keep similar actions at similar speeds
2. **Smooth Transitions**: Design frames to flow naturally
3. **Loop Seamlessly**: Ensure looping animations connect properly
4. **Express Character**: Use animations to show personality

### Performance Optimization
1. **Minimize States**: Don't create unnecessary animation states
2. **Efficient Sheets**: Pack frames tightly in sprite sheets
3. **Appropriate Timing**: Balance smoothness with performance
4. **Reuse Frames**: Share frames between similar animations when possible

### Integration Guidelines
1. **Physics First**: Let physics drive animation states
2. **Priority System**: Implement clear state priority (damage > climb > air > ground)
3. **Smooth Transitions**: Avoid jarring state changes
4. **Consistent Direction**: Handle sprite flipping consistently

## Future Enhancements

### Planned Features
- **Animation Blending**: Smooth transitions between states
- **Variable Frame Timing**: Different durations per frame
- **Animation Events**: Callbacks for specific frames
- **Composite Animations**: Multiple sprite layers
- **Dynamic Timing**: Speed modifications based on game state

### Integration Opportunities
- **Sound System**: Sync audio cues with animation frames
- **Particle Effects**: Trigger effects on specific frames
- **Game Events**: Link gameplay events to animation completion
- **AI Systems**: Use animation states for AI decision making
