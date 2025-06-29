# ROBO-9 Player Entity Implementation Guide

## Overview

This guide covers the implementation of the ROBO-9 player entity, including movement, physics, animation integration, and control handling. The player entity serves as the main character for the platformer game.

## Entity Structure

### Player Entity (`entities/player.go`)

```go
type Player struct {
    // Position and Movement
    X, Y          float64
    VelocityX     float64
    VelocityY     float64
    Width, Height float64
    
    // Physics Constants
    Speed         float64  // Horizontal movement speed
    JumpSpeed     float64  // Initial jump velocity
    Gravity       float64  // Downward acceleration
    Friction      float64  // Velocity reduction factor
    
    // State Flags
    OnGround      bool
    FacingRight   bool
    IsJumping     bool
    IsMoving      bool
    IsClimbing    bool
    IsDamaged     bool
    
    // Animation System
    AnimationController *AnimationController
    
    // Timing
    DamageTimer float64
    DamageTime  float64
}
```

## Physics System

### Movement Constants

```go
Speed:       120.0,  // pixels per second horizontal
JumpSpeed:   200.0,  // pixels per second upward
Gravity:     500.0,  // pixels per second² downward
Friction:    0.8,    // velocity multiplier per frame
```

### Physics Update Loop

```go
func (p *Player) updatePhysics(deltaTime float64) {
    // Apply gravity when airborne
    if !p.OnGround {
        p.VelocityY += p.Gravity * deltaTime
    }
    
    // Apply friction to horizontal movement
    p.VelocityX *= p.Friction
    
    // Update position based on velocity
    p.X += p.VelocityX * deltaTime
    p.Y += p.VelocityY * deltaTime
    
    // Simple ground collision (replace with proper collision detection)
    if p.Y >= groundLevel {
        p.Y = groundLevel
        p.VelocityY = 0
        p.OnGround = true
        p.IsJumping = false
    }
}
```

## Movement System

### Basic Movement Methods

#### Horizontal Movement
```go
func (p *Player) MoveLeft() {
    if !p.IsDamaged {
        p.VelocityX = -p.Speed
        p.FacingRight = false
    }
}

func (p *Player) MoveRight() {
    if !p.IsDamaged {
        p.VelocityX = p.Speed
        p.FacingRight = true
    }
}
```

#### Jumping
```go
func (p *Player) Jump() {
    if p.OnGround && !p.IsDamaged {
        p.VelocityY = -p.JumpSpeed
        p.IsJumping = true
        p.OnGround = false
    }
}
```

#### Climbing System
```go
func (p *Player) StartClimbing() {
    if !p.IsDamaged {
        p.IsClimbing = true
        p.VelocityY = 0
        p.Gravity = 0  // Disable gravity while climbing
    }
}

func (p *Player) ClimbUp() {
    if p.IsClimbing && !p.IsDamaged {
        p.VelocityY = -p.Speed * 0.7  // Slower than walking
    }
}
```

## Input System

### Input Handler (`entities/input.go`)

The input handler translates keyboard input into player actions:

```go
type InputHandler struct {
    player *Player
}

func (ih *InputHandler) Update() {
    // Horizontal movement
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || 
       ebiten.IsKeyPressed(ebiten.KeyA) {
        ih.player.MoveLeft()
    }
    
    // Jumping (key just pressed, not held)
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) || 
       inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
        ih.player.Jump()
    }
}
```

### Control Mapping

| Action | Primary Keys | Alternative Keys |
|--------|-------------|------------------|
| Move Left | Left Arrow | A |
| Move Right | Right Arrow | D |
| Jump | Space | Up Arrow, W |
| Climb Up | Up Arrow | W |
| Climb Down | Down Arrow | S |

### Debug Controls

| Action | Key | Purpose |
|--------|-----|---------|
| Toggle Climb | C | Test climbing mode |
| Test Damage | X | Trigger damage state |

## State Management

### Player States

The player can be in multiple states simultaneously:

1. **Ground States**: `OnGround` (true/false)
2. **Movement States**: `IsMoving`, `IsJumping`, `IsClimbing`
3. **Direction State**: `FacingRight` (true/false)
4. **Special States**: `IsDamaged`

### State Priority System

```go
func (p *Player) updateAnimationState() {
    // 1. Damage (highest priority)
    if p.IsDamaged {
        p.AnimationController.SetState(AnimationDamage)
        return
    }
    
    // 2. Climbing
    if p.IsClimbing {
        p.AnimationController.SetState(AnimationClimb)
        return
    }
    
    // 3. Airborne
    if !p.OnGround {
        if p.VelocityY < 0 {
            p.AnimationController.SetState(AnimationJump)
        } else {
            p.AnimationController.SetState(AnimationFall)
        }
        return
    }
    
    // 4. Ground movement
    if p.IsMoving {
        p.AnimationController.SetState(AnimationWalk)
    } else {
        p.AnimationController.SetState(AnimationIdle)
    }
}
```

## Damage System

### Damage State Management

```go
func (p *Player) TakeDamage() {
    if !p.IsDamaged {
        p.IsDamaged = true
        p.DamageTimer = p.DamageTime  // 1 second immunity
        p.VelocityX = 0  // Stop movement
    }
}
```

### Visual Feedback

```go
func (p *Player) Draw(screen *ebiten.Image) {
    // ... get current frame ...
    
    // Add damage effect (flashing)
    if p.IsDamaged {
        flashCycle := math.Sin(p.DamageTimer * 20)
        if flashCycle > 0 {
            op.ColorM.Scale(1, 1, 1, 0.5)  // Semi-transparent
        }
    }
    
    screen.DrawImage(currentFrame, op)
}
```

## Rendering System

### Sprite Flipping

```go
func (p *Player) Draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    
    // Flip sprite horizontally when facing left
    if !p.FacingRight {
        op.GeoM.Scale(-1, 1)
        op.GeoM.Translate(p.Width, 0)
    }
    
    // Position the sprite
    op.GeoM.Translate(p.X, p.Y)
    
    screen.DrawImage(currentFrame, op)
}
```

### Visual Effects

- **Damage Flashing**: Semi-transparent rendering during damage immunity
- **Direction Flipping**: Automatic sprite mirroring for left movement
- **Animation Smoothness**: Delta-time based animation updates

## Integration Points

### Game Loop Integration

```go
// In main game update loop
func (g *RoboGame) Update() error {
    if g.inputHandler != nil {
        g.inputHandler.Update()  // Handle input
    }
    if g.player != nil {
        g.player.Update(g.deltaTime)  // Update physics and animation
    }
    return g.Game.Update()
}

// In main game draw loop
func (g *RoboGame) drawGameScreen(screen *ebiten.Image) {
    // ... draw background ...
    
    if g.player != nil {
        g.player.Draw(screen)  // Render player
    }
}
```

### Collision System Integration

```go
// Future collision system integration
func (p *Player) checkCollisions(world *World) {
    bounds := p.GetBounds()
    
    // Check platform collisions
    if platform := world.GetPlatformAt(bounds); platform != nil {
        p.HandlePlatformCollision(platform)
    }
    
    // Check climbable surface collisions
    if surface := world.GetClimbableSurfaceAt(bounds); surface != nil {
        p.CanClimb = true
    }
}
```

## Performance Considerations

### Optimization Strategies

1. **Efficient Updates**: Only update physics when state changes
2. **Animation Caching**: Reuse animation frames across instances
3. **State Batching**: Group similar state updates
4. **Collision Optimization**: Spatial partitioning for collision detection

### Memory Usage

- **Player Instance**: ~300 bytes
- **Animation Controller**: ~200 bytes
- **Input Handler**: ~50 bytes
- **Total per Player**: ~550 bytes

## Testing and Debugging

### Debug Information Display

```go
// Debug info rendering
x, y := player.GetPosition()
vx, vy := player.GetVelocity()
state := player.GetAnimationState()

debugInfo := fmt.Sprintf(
    "Pos: (%.1f, %.1f)\nVel: (%.1f, %.1f)\nGround: %v\nAnim: %v",
    x, y, vx, vy, player.IsOnGround(), state
)
```

### Unit Testing

```go
func TestPlayerJump(t *testing.T) {
    player := entities.NewPlayer(0, 0, testSpriteSheet)
    
    // Test jump from ground
    player.OnGround = true
    player.Jump()
    
    assert.True(t, player.IsJumping)
    assert.False(t, player.OnGround)
    assert.Equal(t, -200.0, player.VelocityY)
}
```

## Future Enhancements

### Planned Features

1. **Double Jump**: Secondary jump ability
2. **Wall Sliding**: Sliding down walls before climbing
3. **Dash Attack**: High-speed horizontal movement
4. **Scanner Ability**: Special vision mode
5. **Energy Shield**: Temporary damage immunity

### Extension Points

```go
// Special abilities interface
type Ability interface {
    Activate(player *Player)
    Update(deltaTime float64)
    IsActive() bool
    GetCooldown() float64
}

// Player enhancement
type Player struct {
    // ... existing fields ...
    Abilities []Ability
    Energy    float64
}
```

### Advanced Physics

```go
// Enhanced movement
type Player struct {
    // ... existing fields ...
    Acceleration  float64
    MaxSpeed      float64
    GroundFriction float64
    AirFriction   float64
    WallJumpForce float64
}
```

## Best Practices

### Code Organization

1. **Separation of Concerns**: Keep physics, input, and rendering separate
2. **State Machines**: Use clear state management patterns
3. **Data-Driven Design**: Make values configurable and tweakable
4. **Testing**: Unit test individual components

### Performance Guidelines

1. **Minimize Allocations**: Reuse objects where possible
2. **Efficient Collision**: Use spatial data structures
3. **Animation Optimization**: Cache frame references
4. **Update Frequency**: Only update what changed

### Design Principles

1. **Responsive Controls**: Immediate input response
2. **Predictable Physics**: Consistent movement behavior
3. **Visual Clarity**: Clear animation states
4. **Smooth Experience**: Delta-time based updates

This implementation provides a solid foundation for the ROBO-9 character while remaining extensible for future gameplay enhancements.
