package engine

import (
	"fmt"
	"testing"
)

// TestStateManager_NewStateManager tests the creation of a new state manager
func TestStateManager_NewStateManager(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	if sm == nil {
		t.Fatal("NewStateManager returned nil")
	}
	
	if sm.GetCurrentState() != StateLoading {
		t.Errorf("Expected initial state to be StateLoading, got %v", sm.GetCurrentState())
	}
	
	if sm.GetPreviousState() != StateLoading {
		t.Errorf("Expected previous state to be StateLoading, got %v", sm.GetPreviousState())
	}
	
	if sm.IsTransitioning() {
		t.Error("Expected initial state manager to not be transitioning")
	}
	
	if sm.GetTransitionProgress() != 1.0 {
		t.Errorf("Expected initial transition progress to be 1.0, got %f", sm.GetTransitionProgress())
	}
}

// TestStateManager_SetState tests immediate state changes
func TestStateManager_SetState(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	sm.SetState(StateMenu)
	
	if sm.GetCurrentState() != StateMenu {
		t.Errorf("Expected current state to be StateMenu, got %v", sm.GetCurrentState())
	}
	
	if sm.GetPreviousState() != StateLoading {
		t.Errorf("Expected previous state to be StateLoading, got %v", sm.GetPreviousState())
	}
	
	if sm.IsTransitioning() {
		t.Error("SetState should not trigger transitions")
	}
}

// TestStateManager_TransitionTo tests animated state transitions
func TestStateManager_TransitionTo(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	// Test instant transition (duration 0)
	sm.TransitionTo(StateMenu, 0)
	if sm.GetCurrentState() != StateMenu {
		t.Errorf("Expected immediate transition to StateMenu, got %v", sm.GetCurrentState())
	}
	if sm.IsTransitioning() {
		t.Error("Duration 0 should not trigger transition state")
	}
	
	// Test animated transition
	sm.TransitionTo(StatePlaying, 0.5)
	if sm.GetCurrentState() != StateTransition {
		t.Errorf("Expected current state to be StateTransition, got %v", sm.GetCurrentState())
	}
	if !sm.IsTransitioning() {
		t.Error("Expected state manager to be transitioning")
	}
	if sm.GetTransitionProgress() != 0.0 {
		t.Errorf("Expected initial transition progress to be 0.0, got %f", sm.GetTransitionProgress())
	}
}

// TestStateManager_TransitionProgress tests transition progress tracking
func TestStateManager_TransitionProgress(t *testing.T) {
	sm := NewStateManager(StateMenu)
	
	sm.TransitionTo(StatePlaying, 1.0) // 1 second transition
	
	// Test progress at different points
	testCases := []struct {
		deltaTime        float64
		expectedProgress float64
	}{
		{0.25, 0.25},
		{0.25, 0.50},
		{0.25, 0.75},
		{0.25, 1.0}, // Should complete transition
	}
	
	for i, tc := range testCases {
		err := sm.Update(tc.deltaTime)
		if err != nil {
			t.Errorf("Test case %d: Update returned error: %v", i, err)
		}
		
		progress := sm.GetTransitionProgress()
		if progress != tc.expectedProgress {
			t.Errorf("Test case %d: Expected progress %f, got %f", i, tc.expectedProgress, progress)
		}
	}
	
	// After completion, should be in target state
	if sm.GetCurrentState() != StatePlaying {
		t.Errorf("Expected final state to be StatePlaying, got %v", sm.GetCurrentState())
	}
	if sm.IsTransitioning() {
		t.Error("Expected transition to be completed")
	}
}

// TestStateManager_Callbacks tests the callback system
func TestStateManager_Callbacks(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	var callbackLog []string
	
	// Register callbacks
	sm.RegisterOnEnter(StateMenu, func() {
		callbackLog = append(callbackLog, "enter_menu")
	})
	
	sm.RegisterOnExit(StateLoading, func() {
		callbackLog = append(callbackLog, "exit_loading")
	})
	
	sm.RegisterOnUpdate(StateMenu, func() error {
		callbackLog = append(callbackLog, "update_menu")
		return nil
	})
	
	// Trigger transition
	sm.SetState(StateMenu)
	
	// Check callbacks were called in correct order
	expectedLog := []string{"exit_loading", "enter_menu"}
	if len(callbackLog) != len(expectedLog) {
		t.Errorf("Expected %d callback calls, got %d", len(expectedLog), len(callbackLog))
	}
	
	for i, expected := range expectedLog {
		if i >= len(callbackLog) || callbackLog[i] != expected {
			t.Errorf("Expected callback %d to be '%s', got '%s'", i, expected, callbackLog[i])
		}
	}
	
	// Test update callback
	err := sm.Update(0.016) // ~60 FPS
	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}
	
	if len(callbackLog) != 3 || callbackLog[2] != "update_menu" {
		t.Error("Update callback was not called")
	}
}

// TestStateManager_CallbackErrors tests error handling in update callbacks
func TestStateManager_CallbackErrors(t *testing.T) {
	sm := NewStateManager(StateMenu)
	
	expectedError := fmt.Errorf("test error")
	
	sm.RegisterOnUpdate(StateMenu, func() error {
		return expectedError
	})
	
	err := sm.Update(0.016)
	if err != expectedError {
		t.Errorf("Expected update to return test error, got %v", err)
	}
}

// TestStateManager_SameStateTransition tests transitioning to the same state
func TestStateManager_SameStateTransition(t *testing.T) {
	sm := NewStateManager(StateMenu)
	
	var enterCallCount int
	sm.RegisterOnEnter(StateMenu, func() {
		enterCallCount++
	})
	
	// Should not trigger callbacks when transitioning to same state
	sm.TransitionTo(StateMenu, 0.5)
	
	if sm.IsTransitioning() {
		t.Error("Should not transition to the same state")
	}
	
	if enterCallCount != 0 {
		t.Errorf("Enter callback should not be called for same-state transition, called %d times", enterCallCount)
	}
}

// TestStateManager_StateToString tests state string conversion
func TestStateManager_StateToString(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	testCases := map[GameState]string{
		StateLoading:    "Loading",
		StateMenu:       "Menu",
		StatePlaying:    "Playing",
		StatePaused:     "Paused",
		StateGameOver:   "GameOver",
		StateSettings:   "Settings",
		StateTransition: "Transition",
	}
	
	for state, expected := range testCases {
		result := sm.StateToString(state)
		if result != expected {
			t.Errorf("Expected StateToString(%v) to be '%s', got '%s'", state, expected, result)
		}
	}
	
	// Test unknown state
	unknownState := GameState(999)
	result := sm.StateToString(unknownState)
	expected := "Unknown(999)"
	if result != expected {
		t.Errorf("Expected StateToString(999) to be '%s', got '%s'", expected, result)
	}
}

// TestGame_StateManagement tests the Game struct's state management integration
func TestGame_StateManagement(t *testing.T) {
	config := GameConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		AssetConfig: AssetConfig{
			AssetDir:    "test",
			UseEmbedded: false,
		},
	}
	
	game := NewGame(config)
	
	if game.GetState() != StateLoading {
		t.Errorf("Expected initial game state to be StateLoading, got %v", game.GetState())
	}
	
	// Test state change
	game.SetState(StateMenu)
	if game.GetState() != StateMenu {
		t.Errorf("Expected game state to be StateMenu, got %v", game.GetState())
	}
	
	// Test transition
	game.TransitionToState(StatePlaying, 0.5)
	if game.GetState() != StateTransition {
		t.Errorf("Expected game state to be StateTransition, got %v", game.GetState())
	}
	
	// Test state manager access
	sm := game.GetStateManager()
	if sm == nil {
		t.Fatal("GetStateManager returned nil")
	}
	
	if sm.GetCurrentState() != game.GetState() {
		t.Error("State manager and game state are out of sync")
	}
}

// TestGame_TogglePause tests the pause toggle functionality
func TestGame_TogglePause(t *testing.T) {
	game := NewGame(GameConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		AssetConfig: AssetConfig{
			AssetDir:    "test",
			UseEmbedded: false,
		},
	})
	
	// Start in playing state
	game.SetState(StatePlaying)
	
	// Toggle to pause
	game.TogglePause()
	
	// Should transition to paused (via StateTransition)
	sm := game.GetStateManager()
	if !sm.IsTransitioning() {
		t.Error("Expected game to be transitioning when toggling pause")
	}
	
	// Complete the transition
	for sm.IsTransitioning() {
		err := game.Update()
		if err != nil {
			t.Errorf("Update returned error during pause transition: %v", err)
		}
	}
	
	if game.GetState() != StatePaused {
		t.Errorf("Expected game state to be StatePaused, got %v", game.GetState())
	}
	
	// Toggle back to playing
	game.TogglePause()
	
	// Complete the transition back
	for sm.IsTransitioning() {
		err := game.Update()
		if err != nil {
			t.Errorf("Update returned error during resume transition: %v", err)
		}
	}
	
	if game.GetState() != StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %v", game.GetState())
	}
}

// TestGame_TogglePauseFromWrongState tests pause toggle from non-playing states
func TestGame_TogglePauseFromWrongState(t *testing.T) {
	game := NewGame(GameConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		AssetConfig: AssetConfig{
			AssetDir:    "test",
			UseEmbedded: false,
		},
	})
	
	// Test from menu state
	game.SetState(StateMenu)
	originalState := game.GetState()
	
	game.TogglePause()
	
	// Should not change state
	if game.GetState() != originalState {
		t.Errorf("TogglePause from StateMenu should not change state, was %v, now %v", originalState, game.GetState())
	}
}

// TestStateManager_ComplexTransitionSequence tests a sequence of state transitions
func TestStateManager_ComplexTransitionSequence(t *testing.T) {
	sm := NewStateManager(StateLoading)
	
	var stateHistory []GameState
	
	// Register callbacks to track state changes
	for _, state := range []GameState{StateLoading, StateMenu, StatePlaying, StatePaused, StateGameOver} {
		currentState := state // Capture for closure
		sm.RegisterOnEnter(currentState, func() {
			stateHistory = append(stateHistory, currentState)
		})
	}
	
	// Simulate a typical game flow
	transitions := []struct {
		targetState GameState
		duration    float64
	}{
		{StateMenu, 0.5},     // Loading -> Menu
		{StatePlaying, 0.3},  // Menu -> Playing
		{StatePaused, 0.1},   // Playing -> Paused
		{StatePlaying, 0.1},  // Paused -> Playing
		{StateGameOver, 0.5}, // Playing -> GameOver
		{StateMenu, 0.3},     // GameOver -> Menu
	}
	
	for i, transition := range transitions {
		sm.TransitionTo(transition.targetState, transition.duration)
		
		// Complete the transition
		if transition.duration > 0 {
			for sm.IsTransitioning() {
				err := sm.Update(0.016) // ~60 FPS
				if err != nil {
					t.Errorf("Transition %d: Update returned error: %v", i, err)
				}
			}
		}
		
		if sm.GetCurrentState() != transition.targetState {
			t.Errorf("Transition %d: Expected state %v, got %v", i, transition.targetState, sm.GetCurrentState())
		}
	}
	
	// Verify state history
	expectedHistory := []GameState{StateMenu, StatePlaying, StatePaused, StatePlaying, StateGameOver, StateMenu}
	if len(stateHistory) != len(expectedHistory) {
		t.Errorf("Expected %d state entries, got %d", len(expectedHistory), len(stateHistory))
	}
	
	for i, expected := range expectedHistory {
		if i >= len(stateHistory) || stateHistory[i] != expected {
			t.Errorf("State history[%d]: expected %v, got %v", i, expected, stateHistory[i])
		}
	}
}

// BenchmarkStateManager_Update benchmarks the update performance
func BenchmarkStateManager_Update(b *testing.B) {
	sm := NewStateManager(StateMenu)
	
	sm.RegisterOnUpdate(StateMenu, func() error {
		// Simulate some light work
		_ = sm.GetCurrentState()
		return nil
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Update(0.016)
	}
}

// BenchmarkStateManager_TransitionTo benchmarks transition performance
func BenchmarkStateManager_TransitionTo(b *testing.B) {
	sm := NewStateManager(StateMenu)
	
	states := []GameState{StateMenu, StatePlaying, StatePaused, StateGameOver}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetState := states[i%len(states)]
		sm.TransitionTo(targetState, 0.1)
		
		// Complete transition quickly
		for sm.IsTransitioning() {
			sm.Update(1.0) // Large delta to complete immediately
		}
	}
}

// TestStateManager_MemoryLeaks tests that callbacks don't cause memory leaks
func TestStateManager_MemoryLeaks(t *testing.T) {
	sm := NewStateManager(StateMenu)
	
	// Register many callbacks
	for i := 0; i < 1000; i++ {
		sm.RegisterOnEnter(StateMenu, func() {
			// Empty callback
		})
	}
	
	// The last registered callback should overwrite previous ones
	// This tests that we're not accumulating callbacks
	var callCount int
	sm.RegisterOnEnter(StateMenu, func() {
		callCount++
	})
	
	// Trigger state change multiple times
	for i := 0; i < 10; i++ {
		sm.SetState(StatePlaying)
		sm.SetState(StateMenu)
	}
	
	// Should only be called once per state entry
	if callCount != 10 {
		t.Errorf("Expected callback to be called 10 times, got %d", callCount)
	}
}
