package engine

import (
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
)

// Game represents the main game instance
type Game struct {
	assetManager *AssetManager
	state        GameState
	screenWidth  int
	screenHeight int
}

// GameConfig holds configuration for creating a new game
type GameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	AssetConfig  AssetConfig
}

// NewGame creates a new game instance
func NewGame(config GameConfig) *Game {
	return &Game{
		assetManager: NewAssetManager(config.AssetConfig),
		state:        StateLoading,
		screenWidth:  config.ScreenWidth,
		screenHeight: config.ScreenHeight,
	}
}

// GetAssetManager returns the game's asset manager
func (g *Game) GetAssetManager() *AssetManager {
	return g.assetManager
}

// GetState returns the current game state
func (g *Game) GetState() GameState {
	return g.state
}

// SetState changes the game state
func (g *Game) SetState(state GameState) {
	g.state = state
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	// This will be expanded as we add more game states
	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	// This will be expanded as we add more game states
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}
