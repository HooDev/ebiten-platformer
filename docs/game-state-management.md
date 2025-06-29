# Game State Management System

## Overview

The ROBO-9 platformer implements a robust game state management system that handles different game states, transitions between them, and provides a callback-based architecture for state-specific logic.

## Architecture

### Core Components

#### GameState Enum
```go
type GameState int

const (
    StateLoading GameState = iota
    StateMenu
    StatePlaying
    StatePaused
    StateGameOver
    StateSettings
    StateTransition
)
```

#### StateManager
The `StateManager` is the central component that handles:
- Current state tracking
- State transitions with optional animations
- Callback management for state events
- Transition progress tracking

#### Game Integration
The main `Game` struct integrates the state manager and provides convenience methods for common operations.

## States Description

### StateLoading
- **Purpose**: Display loading screen while assets are being loaded
- **Entry**: Game startup
- **Exit**: When all assets are loaded successfully
- **Duration**: Variable (depends on asset loading time)

### StateMenu
- **Purpose**: Main menu navigation
- **Entry**: After loading completes or from in-game menu
- **Exit**: When starting game or entering settings
- **Controls**: 
  - ENTER: Start game
  - S: Settings

### StatePlaying
- **Purpose**: Active gameplay
- **Entry**: From menu or resume from pause
- **Exit**: Pause, game over, or return to menu
- **Controls**:
  - ESC/SPACE: Pause
  - M: Return to menu
  - G: Simulate game over (debug)

### StatePaused
- **Purpose**: Game is paused, overlay on game screen
- **Entry**: From playing state
- **Exit**: Resume to playing or return to menu
- **Controls**:
  - ESC/SPACE: Resume
  - M: Return to menu

### StateGameOver
- **Purpose**: Display game over screen and options
- **Entry**: When player dies or fails level
- **Exit**: Restart game or return to menu
- **Controls**:
  - ENTER/R: Restart game
  - M/ESC: Return to menu

### StateSettings
- **Purpose**: Configuration and options menu
- **Entry**: From main menu
- **Exit**: Return to menu
- **Controls**:
  - ESC/BACKSPACE: Return to menu

### StateTransition
- **Purpose**: Special state for animated transitions
- **Entry**: Automatically when transitioning between states
- **Exit**: Automatically when transition completes
- **Features**: Progress tracking, fade effects

## API Reference

### StateManager Methods

#### State Query Methods
```go
GetCurrentState() GameState          // Returns current state
GetPreviousState() GameState         // Returns previous state
IsTransitioning() bool               // Check if transition in progress
GetTransitionProgress() float64      // Get transition progress (0.0-1.0)
```

#### State Change Methods
```go
SetState(newState GameState)                    // Immediate state change
TransitionTo(newState GameState, duration float64)  // Animated transition
```

#### Callback Registration
```go
RegisterOnEnter(state GameState, callback func())        // Called when entering state
RegisterOnExit(state GameState, callback func())         // Called when exiting state
RegisterOnUpdate(state GameState, callback func() error) // Called every frame in state
```

### Game Methods

#### High-Level State Control
```go
SetState(state GameState)                         // Immediate state change
TransitionToState(state GameState, duration float64)  // Animated transition
TogglePause()                                     // Toggle between playing/paused
```

#### State Information
```go
GetState() GameState                              // Get current state
GetStateManager() *StateManager                  // Access state manager directly
```

## Usage Examples

### Basic State Changes
```go
// Immediate state change
game.SetState(engine.StateMenu)

// Animated transition (0.5 second fade)
game.TransitionToState(engine.StatePlaying, 0.5)

// Toggle pause with quick transition
game.TogglePause()
```

### Registering State Callbacks
```go
stateManager := game.GetStateManager()

// Register callback for entering playing state
stateManager.RegisterOnEnter(engine.StatePlaying, func() {
    log.Println("Game started!")
    // Initialise player, load level, etc.
})

// Register update callback for playing state
stateManager.RegisterOnUpdate(engine.StatePlaying, func() error {
    // Handle input, update physics, etc.
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        game.TogglePause()
    }
    return nil
})

// Register exit callback
stateManager.RegisterOnExit(engine.StatePlaying, func() {
    log.Println("Leaving game")
    // Save progress, cleanup, etc.
})
```

### Custom Transition Effects
```go
// Quick transitions for responsive feel
game.TransitionToState(engine.StatePaused, 0.2)

// Slower transitions for dramatic effect
game.TransitionToState(engine.StateGameOver, 1.0)

// Instant transitions when needed
game.SetState(engine.StateMenu)
```

## Implementation Details

### Delta Time Calculation
The system currently uses a fixed 60 FPS assumption for delta time:
```go
deltaTime := 1.0 / 60.0
```

This can be enhanced with actual frame timing for more accurate transitions.

### Transition System
Transitions work by:
1. Setting `isTransitioning = true`
2. Changing current state to `StateTransition`
3. Accumulating transition time each frame
4. Calling `completeTransition()` when duration reached
5. Executing enter/exit callbacks appropriately

### Callback Execution Order
1. Exit callback for previous state
2. State change
3. Enter callback for new state
4. Update callbacks every frame

### Memory Management
- Callbacks are stored in maps indexed by GameState
- No automatic cleanup - callbacks persist for game lifetime
- Only one callback per state per event type

## Testing and Debugging

### Debug Controls
The current implementation includes debug controls for testing:
- **M**: Force transition to menu (from any state)
- **G**: Force game over (from playing state)
- **ESC/SPACE**: Standard pause toggle

### Logging
State transitions are automatically logged:
```
Transitioning from Playing to Paused
Game paused
Completed transition to Paused
```

### Transition Progress
You can monitor transition progress for custom effects:
```go
if stateManager.IsTransitioning() {
    progress := stateManager.GetTransitionProgress()
    // Use progress for fade effects, animations, etc.
}
```

## Best Practices

### State Design
1. **Keep states focused**: Each state should have a clear, single purpose
2. **Use callbacks wisely**: Register only necessary callbacks to avoid complexity
3. **Handle transitions gracefully**: Always provide visual feedback during transitions

### Performance Considerations
1. **Avoid heavy operations in callbacks**: Keep enter/exit callbacks lightweight
2. **Use appropriate transition durations**: Balance responsiveness with visual appeal
3. **Cache resources**: Don't reload assets on every state change

### Error Handling
1. **Update callbacks can return errors**: Use this for graceful error handling
2. **State validation**: Ensure valid state transitions
3. **Fallback states**: Always have a way to return to a safe state (menu)

### Future Enhancements

#### Potential Improvements
1. **State Stack**: Support for state stacking (e.g., settings over gameplay)
2. **Advanced Transitions**: Multiple transition types (slide, fade, zoom)
3. **State History**: Navigation history for back button functionality
4. **Conditional Transitions**: Rules-based state transitions
5. **Async State Loading**: Background loading for smoother transitions

#### Integration Points
1. **Save System**: State-aware saving and loading
2. **Audio Manager**: State-based music and sound management
3. **Input Manager**: State-specific input handling
4. **Render Pipeline**: State-aware rendering optimizations

## Conclusion

The game state management system provides a solid foundation for the ROBO-9 platformer. It's extensible, well-organized, and provides clear separation of concerns between different game modes. The callback system makes it easy to implement state-specific logic while maintaining clean code organization.
