# ROBO-9 Sprite Sheet Specification

## Overview

This document describes the sprite sheet requirements for the ROBO-9 player character, including frame layouts, animation sequences, and technical specifications.

## Sprite Sheet Layout

### Dimensions
- **Frame Size**: 32×32 pixels per frame
- **Sheet Layout**: 6 columns × 3 rows (18 frames total)
- **Total Size**: 192×96 pixels
- **Format**: PNG with transparency support

### Frame Organization

The sprite sheet is organized in a horizontal layout with frames numbered from left to right, top to bottom (0-17):

```
Row 1: [ 0] [ 1] [ 2] [ 3] [ 4] [ 5]
Row 2: [ 6] [ 7] [ 8] [ 9] [10] [11]
Row 3: [12] [13] [14] [15] [16] [17]
```

## Animation Sequences

### 1. Idle Animation (Frames 0-3)
- **Purpose**: ROBO-9 standing still
- **Frame Count**: 4 frames
- **Duration**: 0.2 seconds per frame (0.8s total cycle)
- **Loop**: Yes
- **Description**: Subtle breathing or power indicator animation

**Frame Details:**
- Frame 0: Base idle pose
- Frame 1: Slight variation (breathing/power pulse)
- Frame 2: Return to base or slight variation
- Frame 3: Back to starting position

### 2. Walk Animation (Frames 4-7)
- **Purpose**: ROBO-9 walking/running
- **Frame Count**: 4 frames
- **Duration**: 0.1 seconds per frame (0.4s total cycle)
- **Loop**: Yes
- **Description**: Side-scrolling walk cycle

**Frame Details:**
- Frame 4: Left foot forward, right foot back
- Frame 5: Both feet together (mid-stride)
- Frame 6: Right foot forward, left foot back
- Frame 7: Both feet together (mid-stride)

### 3. Jump Animation (Frames 8-9)
- **Purpose**: ROBO-9 jumping upward
- **Frame Count**: 2 frames
- **Duration**: 0.1 seconds per frame
- **Loop**: No (plays once)
- **Description**: Launch sequence for jumping

**Frame Details:**
- Frame 8: Crouch/prepare position (legs bent)
- Frame 9: Extended position (legs straight, arms up)

### 4. Fall Animation (Frames 10-11)
- **Purpose**: ROBO-9 falling through air
- **Frame Count**: 2 frames
- **Duration**: 0.15 seconds per frame
- **Loop**: Yes
- **Description**: Falling/airborne animation

**Frame Details:**
- Frame 10: Arms slightly raised for balance
- Frame 11: Arms in different position, showing air resistance

### 5. Climb Animation (Frames 12-15)
- **Purpose**: ROBO-9 climbing on metallic surfaces
- **Frame Count**: 4 frames
- **Duration**: 0.15 seconds per frame (0.6s total cycle)
- **Loop**: Yes
- **Description**: Wall climbing motion

**Frame Details:**
- Frame 12: Right hand up, left hand down
- Frame 13: Both hands at middle position
- Frame 14: Left hand up, right hand down
- Frame 15: Both hands at middle position

### 6. Damage Animation (Frames 16-17)
- **Purpose**: ROBO-9 taking damage
- **Frame Count**: 2 frames
- **Duration**: 0.1 seconds per frame
- **Loop**: No (plays once)
- **Description**: Impact reaction and recovery

**Frame Details:**
- Frame 16: Recoil position (knocked back)
- Frame 17: Recovery position (regaining balance)

## Art Style Guidelines

### Visual Design
- **Style**: Pixel art with clean, defined edges
- **Robot Aesthetic**: Mechanical appearance with visible joints, panels, or details
- **Color Palette**: 
  - Primary: Cool blues and metallic grays
  - Accent: Bright energy colors (blue, cyan) for power indicators
  - Damage: Red highlights for damage states

### Technical Requirements
- **Transparency**: Use alpha channel for non-character areas
- **Pixel Perfect**: Align to pixel grid for crisp rendering
- **Consistent Lighting**: Maintain consistent light source across all frames
- **Edge Definition**: Clear contrast between character and background

### Character Features
- **Eyes**: Visible LED-style eyes that can show expression
- **Body Panels**: Segmented robot body with visible joints
- **Power Indicators**: Small lights or energy effects
- **Proportions**: Compact, sturdy robot design suitable for platforming

## Implementation Notes

### Sprite Sheet Loading
- File location: `assets/player.png`
- Fallback: If not found, system generates a test sprite sheet
- Format support: PNG, with other formats possible via asset manager

### Animation System Integration
- Frame extraction: Automatic based on 32×32 grid
- State management: Handled by `AnimationController`
- Timing: Uses delta time for smooth, framerate-independent animation
- Mirroring: Automatic horizontal flip for left-facing direction

### Performance Considerations
- **Memory**: Single sprite sheet minimizes texture switches
- **GPU**: Efficient sub-image rendering using Ebitengine
- **Animation**: Lightweight frame-based system with minimal CPU overhead

## Future Expansion

### Additional Animations (Planned)
- **Double Jump**: Extended air maneuvers
- **Wall Slide**: Sliding down walls
- **Interact**: Giving hearts to cats
- **Death**: Destruction/shutdown sequence
- **Special Abilities**: Scanner, shield activation

### Sprite Sheet Evolution
- Current version supports core movement
- Future versions may expand to larger grids (8×4, 10×5)
- Modular design allows easy addition of new animation states

## Example Usage

```go
// Create player with sprite sheet
player := entities.NewPlayer(x, y, spriteSheetImage)

// Animation states are automatically configured
// No additional setup required for basic animations

// Custom animation timing (if needed)
player.AnimationController.AddAnimation(
    AnimationCustom, 
    startFrame,    // Starting frame index
    frameCount,    // Number of frames
    frameTime,     // Time per frame in seconds
    loop           // Whether to loop
)
```

## Testing

The current implementation includes a test sprite sheet generator (`entities.CreateTestSpriteSheet()`) that creates colored rectangles for each animation state:

- **Blue**: Idle frames
- **Turquoise**: Walk frames
- **Gold**: Jump frames
- **Orange**: Fall frames
- **Green**: Climb frames
- **Red**: Damage frames

This allows for immediate testing of the animation system without requiring final art assets.
