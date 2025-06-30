# Coyote Time Implementation

## Overview

Coyote time is a game design pattern that provides players with a short grace period to jump after leaving a platform edge. This creates more forgiving and responsive platformer controls, reducing frustration from pixel-perfect timing requirements.

## Implementation Details

### Core Concept

The player can jump for a brief window (100ms) after walking off a platform edge, even when they're already falling. This is named after the Road Runner cartoons where Wile E. Coyote would briefly float in mid-air before falling.

### Key Components

#### Player struct fields:
```go
// Coyote time for forgiving jumps
CoyoteTime  float64 // Duration of coyote time window (0.1 seconds)
CoyoteTimer float64 // Current coyote time remaining
WasOnGroundPhysics bool // Previous frame physics ground state
```

#### Core Logic (updateCoyoteTime):
1. **Activation**: When `WasOnGroundPhysics` is true but `OnGround` becomes false (and player isn't jumping)
2. **Countdown**: Timer decreases each frame by `deltaTime`
3. **Reset**: Timer resets to 0 when player lands on ground
4. **Consumption**: Timer resets to 0 when player jumps during coyote time

### Ground Detection Integration

The coyote time system relies on the robust ground detection logic in `level/level.go`:

- **50% Overlap Rule**: Player must have at least 50% horizontal overlap with a solid tile to be considered "on ground"
- **Ground Tolerance**: 2-pixel tolerance for vertical ground contact detection
- **Edge Detection**: Player falls off when overlap drops below 50%

### Jump Logic Enhancement

```go
func (p *Player) Jump() {
    // Allow jumping if on ground OR during coyote time window
    canJump := (p.OnGround || p.CoyoteTimer > 0) && !p.IsDamaged

    if canJump {
        p.VelocityY = -p.JumpSpeed
        p.IsJumping = true
        p.OnGround = false
        p.CoyoteTimer = 0 // Consume coyote time
    }
}
```

## Configuration

### Default Settings
- **Duration**: 100ms (0.1 seconds)
- **Tolerance**: Works with existing ground detection thresholds
- **Activation**: Only when walking/falling off platforms (not when jumping)

### Tuning Parameters
- `CoyoteTime`: Adjust duration (0.05-0.2s typical range)
- Ground detection tolerance in `level/level.go`
- Platform overlap requirements (currently 50%)

## Testing

### Automated Tests
- `TestCoyoteTime`: Basic functionality - jump during coyote window
- `TestCoyoteTimeExpires`: Timer expiration - cannot jump after timeout
- `TestCoyoteTimeForgivingDetection`: Edge case - gradual platform departure

### Test Scenarios
```go
// Small platform setup for reliable testing
testLevel.SetTile(2, 7, level.TileSolid)
testLevel.SetTile(3, 7, level.TileSolid)

// Player movement to trigger coyote time
for i := 0; i < 25; i++ {
    player.MoveRight()
    player.Update(deltaTime)
}
```

## Usage Examples

### Basic Coyote Time Check
```go
// In game loop
if input.JumpPressed() {
    player.Jump() // Will work if OnGround OR CoyoteTimer > 0
}

// Debug info
fmt.Printf("Coyote Timer: %.3f\n", player.GetCoyoteTimer())
```

### Manual Activation (for special cases)
```go
// Force activate coyote time (if needed for special mechanics)
player.CoyoteTimer = player.CoyoteTime
```

## Performance Considerations

- **Minimal Overhead**: Simple timer countdown per frame
- **No Complex Calculations**: Uses existing ground detection
- **Efficient State Tracking**: Single boolean for previous ground state

## Debugging

### Debug Information
```go
// Available debug info
debugInfo := player.GetDebugInfo()
// Output: "OnGround: false, WasOnGround: true, CoyoteTimer: 0.067, VelocityY: 25.0"
```

### Common Issues
1. **Timer Not Activating**: Check ground detection is working properly
2. **Timer Not Expiring**: Verify deltaTime is being passed correctly
3. **False Activation**: Ensure player isn't jumping when leaving ground

## Integration Notes

### Required Dependencies
- Ground detection system (`level/level.go`)
- Physics update cycle (`updatePhysics`)
- Input handling for jump commands

### Interaction with Other Systems
- **Animation**: Works with existing jump/fall animations
- **Physics**: Integrates with gravity and collision detection
- **Input**: Compatible with existing jump input handling

## Future Enhancements

### Potential Improvements
1. **Variable Duration**: Different coyote time for different platform types
2. **Visual Feedback**: Particle effects or animation cues
3. **Audio Feedback**: Sound cues for coyote time activation
4. **Difficulty Scaling**: Adjust duration based on game difficulty

### Related Features
- **Wall Jumping**: Could use similar grace period logic
- **Double Jumping**: Interaction with coyote time consumption
- **Moving Platforms**: Special handling for dynamic surfaces

## Technical Notes

### Frame Rate Independence
The system uses `deltaTime` for frame-rate independent timing, ensuring consistent behavior at different FPS.

### State Machine Integration
Coyote time integrates cleanly with the existing player state machine without requiring major architectural changes.

### Memory Usage
Minimal memory footprint - only adds 3 float64 fields to the Player struct.
