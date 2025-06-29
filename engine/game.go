package engine

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameState represents different states the game can be in
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

// StateManager handles game state transitions and callbacks
type StateManager struct {
	currentState  GameState
	previousState GameState
	targetState   GameState
	isTransitioning bool
	transitionTime  float64
	maxTransitionTime float64
	onEnterCallbacks map[GameState]func()
	onExitCallbacks  map[GameState]func()
	onUpdateCallbacks map[GameState]func() error
}

// NewStateManager creates a new state manager
func NewStateManager(initialState GameState) *StateManager {
	return &StateManager{
		currentState:     initialState,
		previousState:    initialState,
		targetState:      initialState,
		onEnterCallbacks: make(map[GameState]func()),
		onExitCallbacks:  make(map[GameState]func()),
		onUpdateCallbacks: make(map[GameState]func() error),
	}
}

// GetCurrentState returns the current game state
func (sm *StateManager) GetCurrentState() GameState {
	return sm.currentState
}

// GetPreviousState returns the previous game state
func (sm *StateManager) GetPreviousState() GameState {
	return sm.previousState
}

// IsTransitioning returns true if a state transition is in progress
func (sm *StateManager) IsTransitioning() bool {
	return sm.isTransitioning
}

// GetTransitionProgress returns progress of current transition (0.0 to 1.0)
func (sm *StateManager) GetTransitionProgress() float64 {
	if !sm.isTransitioning || sm.maxTransitionTime == 0 {
		return 1.0
	}
	return sm.transitionTime / sm.maxTransitionTime
}

// TransitionTo initiates a transition to a new state
func (sm *StateManager) TransitionTo(newState GameState, duration float64) {
	if sm.currentState == newState {
		return
	}

	log.Printf("Transitioning from %s to %s", sm.StateToString(sm.currentState), sm.StateToString(newState))

	sm.previousState = sm.currentState
	sm.targetState = newState
	sm.maxTransitionTime = duration
	sm.transitionTime = 0

	if duration > 0 {
		sm.isTransitioning = true
		sm.currentState = StateTransition
	} else {
		sm.completeTransition()
	}
}

// SetState immediately changes to a new state without transition
func (sm *StateManager) SetState(newState GameState) {
	sm.TransitionTo(newState, 0)
}

// Update handles state transitions and calls state-specific update functions
func (sm *StateManager) Update(deltaTime float64) error {
	if sm.isTransitioning {
		sm.transitionTime += deltaTime
		if sm.transitionTime >= sm.maxTransitionTime {
			sm.completeTransition()
		}
	}

	// Call state-specific update callback
	if callback, exists := sm.onUpdateCallbacks[sm.currentState]; exists {
		return callback()
	}

	return nil
}

// completeTransition finishes a state transition
func (sm *StateManager) completeTransition() {
	// Call exit callback for previous state
	if callback, exists := sm.onExitCallbacks[sm.previousState]; exists {
		callback()
	}

	// Change to target state
	sm.currentState = sm.targetState
	sm.isTransitioning = false
	sm.transitionTime = 0

	// Call enter callback for new state
	if callback, exists := sm.onEnterCallbacks[sm.currentState]; exists {
		callback()
	}

	log.Printf("Completed transition to %s", sm.StateToString(sm.currentState))
}

// RegisterOnEnter registers a callback for when entering a specific state
func (sm *StateManager) RegisterOnEnter(state GameState, callback func()) {
	sm.onEnterCallbacks[state] = callback
}

// RegisterOnExit registers a callback for when exiting a specific state
func (sm *StateManager) RegisterOnExit(state GameState, callback func()) {
	sm.onExitCallbacks[state] = callback
}

// RegisterOnUpdate registers a callback for updating a specific state
func (sm *StateManager) RegisterOnUpdate(state GameState, callback func() error) {
	sm.onUpdateCallbacks[state] = callback
}

// StateToString converts a GameState to a readable string
func (sm *StateManager) StateToString(state GameState) string {
	switch state {
	case StateLoading:
		return "Loading"
	case StateMenu:
		return "Menu"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	case StateGameOver:
		return "GameOver"
	case StateSettings:
		return "Settings"
	case StateTransition:
		return "Transition"
	default:
		return fmt.Sprintf("Unknown(%d)", int(state))
	}
}

// Game represents the main game instance
type Game struct {
	assetManager *AssetManager
	stateManager *StateManager
	screenWidth  int
	screenHeight int
	lastFrameTime float64
}

// GameConfig holds configuration for creating a new game
type GameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	AssetConfig  AssetConfig
}

// NewGame creates a new game instance
func NewGame(config GameConfig) *Game {
	game := &Game{
		assetManager: NewAssetManager(config.AssetConfig),
		stateManager: NewStateManager(StateLoading),
		screenWidth:  config.ScreenWidth,
		screenHeight: config.ScreenHeight,
	}

	// Register default state callbacks
	game.setupDefaultStateCallbacks()

	return game
}

// setupDefaultStateCallbacks sets up default behavior for each state
func (g *Game) setupDefaultStateCallbacks() {
	// Loading state callbacks
	g.stateManager.RegisterOnEnter(StateLoading, func() {
		log.Println("Entered Loading state")
	})

	g.stateManager.RegisterOnExit(StateLoading, func() {
		log.Println("Exited Loading state")
	})

	// Menu state callbacks
	g.stateManager.RegisterOnEnter(StateMenu, func() {
		log.Println("Entered Menu state")
	})

	// Playing state callbacks
	g.stateManager.RegisterOnEnter(StatePlaying, func() {
		log.Println("Entered Playing state - Game started!")
	})

	// Paused state callbacks
	g.stateManager.RegisterOnEnter(StatePaused, func() {
		log.Println("Game paused")
	})

	g.stateManager.RegisterOnExit(StatePaused, func() {
		log.Println("Game resumed")
	})

	// Game Over state callbacks
	g.stateManager.RegisterOnEnter(StateGameOver, func() {
		log.Println("Game Over!")
	})
}

// GetAssetManager returns the game's asset manager
func (g *Game) GetAssetManager() *AssetManager {
	return g.assetManager
}

// GetState returns the current game state
func (g *Game) GetState() GameState {
	return g.stateManager.GetCurrentState()
}

// GetStateManager returns the state manager
func (g *Game) GetStateManager() *StateManager {
	return g.stateManager
}

// SetState changes the game state immediately
func (g *Game) SetState(state GameState) {
	g.stateManager.SetState(state)
}

// TransitionToState changes the game state with a transition
func (g *Game) TransitionToState(state GameState, duration float64) {
	g.stateManager.TransitionTo(state, duration)
}

// TogglePause toggles between playing and paused states
func (g *Game) TogglePause() {
	currentState := g.stateManager.GetCurrentState()
	if currentState == StatePlaying {
		g.TransitionToState(StatePaused, 0.2) // Quick fade transition
	} else if currentState == StatePaused {
		g.TransitionToState(StatePlaying, 0.2)
	}
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	// Calculate delta time (assuming 60 FPS for now)
	deltaTime := 1.0 / 60.0

	// Update state manager
	return g.stateManager.Update(deltaTime)
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	// Base drawing will be handled by the specific game implementation
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}
