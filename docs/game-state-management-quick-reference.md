# Game State Management - Quick Reference

## State Diagram

```
StateLoading
    ↓ (assets loaded)
StateMenu ←→ StateSettings
    ↓ (start game)
StatePlaying ←→ StatePaused
    ↓ (game over)
StateGameOver
    ↓ (restart/menu)
Back to StateMenu or StatePlaying
```

## Common Operations

### Immediate State Changes
```go
game.SetState(engine.StateMenu)
game.SetState(engine.StatePlaying)
```

### Animated Transitions
```go
// Quick pause (0.2s)
game.TogglePause()

// Menu transition (0.5s)
game.TransitionToState(engine.StateMenu, 0.5)

// Dramatic game over (1.0s)
game.TransitionToState(engine.StateGameOver, 1.0)
```

### State Callbacks
```go
stateManager := game.GetStateManager()

// One-time setup when entering state
stateManager.RegisterOnEnter(engine.StatePlaying, func() {
    // Initialize level, reset player position, etc.
})

// Continuous updates while in state
stateManager.RegisterOnUpdate(engine.StatePlaying, func() error {
    // Handle input, update physics, check win conditions
    return nil
})

// Cleanup when leaving state
stateManager.RegisterOnExit(engine.StatePlaying, func() {
    // Save progress, pause audio, etc.
})
```

## Default Controls

| State | Key | Action |
|-------|-----|--------|
| Menu | ENTER | Start Game |
| Menu | S | Settings |
| Playing | ESC/SPACE | Pause |
| Playing | M | Menu |
| Playing | G | Game Over (debug) |
| Paused | ESC/SPACE | Resume |
| Paused | M | Menu |
| GameOver | ENTER/R | Restart |
| GameOver | M/ESC | Menu |
| Settings | ESC/BACKSPACE | Back to Menu |

## State Properties

| State | Persistent | Overlay | Animated Entry | Typical Duration |
|-------|------------|---------|----------------|------------------|
| Loading | No | No | No | Variable |
| Menu | Yes | No | Yes | User-controlled |
| Playing | Yes | No | Yes | User-controlled |
| Paused | No | Yes | Yes | User-controlled |
| GameOver | No | No | Yes | User-controlled |
| Settings | Yes | No | Yes | User-controlled |
| Transition | No | Yes | N/A | 0.2-1.0s |

## Debugging

### Enable State Logging
State transitions are automatically logged to console:
```
Transitioning from Playing to Paused
Game paused
Completed transition to Paused
```

### Check Current State
```go
currentState := game.GetState()
isTransitioning := game.GetStateManager().IsTransitioning()
progress := game.GetStateManager().GetTransitionProgress()
```

### Force State Changes (Debug)
```go
// For testing - bypass normal game flow
game.SetState(engine.StateGameOver)  // Test game over screen
game.SetState(engine.StateMenu)      // Return to menu quickly
```

## Performance Tips

1. **Keep callbacks lightweight** - Heavy operations should be in separate systems
2. **Cache resources** - Don't reload assets on every state change
3. **Use appropriate transition times** - Balance feel with responsiveness
4. **Avoid nested state changes** - Don't change state within state callbacks

## Common Patterns

### Level Transition
```go
// Playing → Loading → Playing (new level)
stateManager.RegisterOnEnter(engine.StatePlaying, func() {
    if needsNewLevel {
        loadLevel(nextLevelID)
    }
})
```

### Save on State Exit
```go
stateManager.RegisterOnExit(engine.StatePlaying, func() {
    saveGame()
})
```

### Audio Management
```go
stateManager.RegisterOnEnter(engine.StateMenu, func() {
    audioManager.PlayMusic("menu_theme")
})

stateManager.RegisterOnEnter(engine.StatePlaying, func() {
    audioManager.PlayMusic("gameplay_theme")
})
```
