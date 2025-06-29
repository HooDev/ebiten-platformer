package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"ebiten-platformer/engine"
)

// RoboGame extends the base engine.Game with platformer-specific logic
type RoboGame struct {
	*engine.Game
	playerImage *ebiten.Image
	overlayImage *ebiten.Image
}

// NewRoboGame creates a new platformer game instance
func NewRoboGame() *RoboGame {
	config := engine.GameConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		AssetConfig: engine.AssetConfig{
			AssetDir:    "assets",
			UseEmbedded: false,
		},
	}

	baseGame := engine.NewGame(config)
	
	roboGame := &RoboGame{
		Game: baseGame,
		overlayImage: ebiten.NewImage(320, 240),
	}

	// Set up game-specific state callbacks
	roboGame.setupGameStateCallbacks()

	return roboGame
}

// setupGameStateCallbacks configures state-specific behavior for the platformer
func (g *RoboGame) setupGameStateCallbacks() {
	stateManager := g.GetStateManager()

	// Override loading state update to handle asset loading completion
	stateManager.RegisterOnUpdate(engine.StateLoading, func() error {
		// In a real implementation, you'd check if assets are loaded
		// For now, we'll transition after the assets are manually loaded
		return nil
	})

	// Playing state input handling
	stateManager.RegisterOnUpdate(engine.StatePlaying, func() error {
		// Handle pause input
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.TogglePause()
		}

		// Handle transition to menu (for testing)
		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			g.TransitionToState(engine.StateMenu, 0.5)
		}

		// Handle game over simulation (for testing)
		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			g.TransitionToState(engine.StateGameOver, 0.3)
		}

		return nil
	})

	// Paused state input handling
	stateManager.RegisterOnUpdate(engine.StatePaused, func() error {
		// Handle resume input
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.TogglePause()
		}

		// Handle transition to menu
		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			g.TransitionToState(engine.StateMenu, 0.5)
		}

		return nil
	})

	// Menu state input handling
	stateManager.RegisterOnUpdate(engine.StateMenu, func() error {
		// Handle start game
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.TransitionToState(engine.StatePlaying, 0.5)
		}

		// Handle settings
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.TransitionToState(engine.StateSettings, 0.3)
		}

		return nil
	})

	// Settings state input handling
	stateManager.RegisterOnUpdate(engine.StateSettings, func() error {
		// Handle back to menu
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			g.TransitionToState(engine.StateMenu, 0.3)
		}

		return nil
	})

	// Game Over state input handling
	stateManager.RegisterOnUpdate(engine.StateGameOver, func() error {
		// Handle restart
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.TransitionToState(engine.StatePlaying, 0.5)
		}

		// Handle back to menu
		if inpututil.IsKeyJustPressed(ebiten.KeyM) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.TransitionToState(engine.StateMenu, 0.5)
		}

		return nil
	})
}

// LoadAssets loads all game assets
func (g *RoboGame) LoadAssets() error {
	assetManager := g.GetAssetManager()
	
	// Load player image
	playerImg, err := assetManager.LoadImage("player.png")
	if err != nil {
		return err
	}
	g.playerImage = playerImg

	// Set state to menu after assets are loaded
	g.SetState(engine.StateMenu)
	
	log.Println("All assets loaded successfully")
	return nil
}

// Update implements ebiten.Game interface
func (g *RoboGame) Update() error {
	// Call base game update (handles state management)
	return g.Game.Update()
}

// Draw implements ebiten.Game interface
func (g *RoboGame) Draw(screen *ebiten.Image) {
	// Call base game draw
	g.Game.Draw(screen)

	stateManager := g.GetStateManager()
	currentState := stateManager.GetCurrentState()

	// Handle state-specific rendering
	switch currentState {
	case engine.StateLoading:
		g.drawLoadingScreen(screen)
	case engine.StateMenu:
		g.drawMenuScreen(screen)
	case engine.StatePlaying:
		g.drawGameScreen(screen)
	case engine.StatePaused:
		g.drawPausedScreen(screen)
	case engine.StateGameOver:
		g.drawGameOverScreen(screen)
	case engine.StateSettings:
		g.drawSettingsScreen(screen)
	case engine.StateTransition:
		g.drawTransitionScreen(screen)
	}
}

// drawLoadingScreen renders the loading screen
func (g *RoboGame) drawLoadingScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255}) // Dark blue background
	ebitenutil.DebugPrintAt(screen, "Loading ROBO-9...", 10, 10)
	
	// Simple loading animation
	dots := ""
	for i := 0; i < int(ebiten.TPS()/20)%4; i++ {
		dots += "."
	}
	ebitenutil.DebugPrintAt(screen, "Please wait"+dots, 10, 30)
}

// drawMenuScreen renders the main menu
func (g *RoboGame) drawMenuScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 60, 100, 255}) // Blue-gray background
	
	// Title
	ebitenutil.DebugPrintAt(screen, "ROBO-9 PLATFORMER", 80, 50)
	ebitenutil.DebugPrintAt(screen, "================", 80, 65)
	
	// Menu options
	ebitenutil.DebugPrintAt(screen, "ENTER - Start Game", 90, 100)
	ebitenutil.DebugPrintAt(screen, "S - Settings", 90, 120)
	
	// Instructions
	ebitenutil.DebugPrintAt(screen, "Controls:", 10, 180)
	ebitenutil.DebugPrintAt(screen, "ESC/SPACE - Pause", 10, 195)
	ebitenutil.DebugPrintAt(screen, "M - Menu", 10, 210)
}

// drawGameScreen renders the main game
func (g *RoboGame) drawGameScreen(screen *ebiten.Image) {
	// Clear screen with sky blue
	screen.Fill(color.RGBA{135, 206, 235, 255})
	
	// Game title and info
	ebitenutil.DebugPrint(screen, "ROBO-9 Platformer - PLAYING")
	ebitenutil.DebugPrintAt(screen, "ESC/SPACE: Pause | M: Menu | G: Game Over", 10, 20)

	// Draw player if loaded
	if g.playerImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(100, 100)
		screen.DrawImage(g.playerImage, op)
	}

	// Display asset manager stats
	assetManager := g.GetAssetManager()
	stats := fmt.Sprintf("Assets loaded:\nImages: %d\nAudio: %d", 
		assetManager.GetLoadedImageCount(), 
		assetManager.GetLoadedAudioCount())
	ebitenutil.DebugPrintAt(screen, stats, 200, 50)
}

// drawPausedScreen renders the pause overlay
func (g *RoboGame) drawPausedScreen(screen *ebiten.Image) {
	// Draw the game screen first (as background)
	g.drawGameScreen(screen)
	
	// Add semi-transparent overlay using reusable image
	g.overlayImage.Fill(color.RGBA{0, 0, 0, 128}) // Semi-transparent black
	screen.DrawImage(g.overlayImage, nil)
	
	// Pause text
	ebitenutil.DebugPrintAt(screen, "GAME PAUSED", 120, 100)
	ebitenutil.DebugPrintAt(screen, "===========", 120, 115)
	ebitenutil.DebugPrintAt(screen, "ESC/SPACE - Resume", 100, 140)
	ebitenutil.DebugPrintAt(screen, "M - Main Menu", 100, 160)
}

// drawGameOverScreen renders the game over screen
func (g *RoboGame) drawGameOverScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{60, 20, 20, 255}) // Dark red background
	
	ebitenutil.DebugPrintAt(screen, "GAME OVER", 120, 100)
	ebitenutil.DebugPrintAt(screen, "=========", 120, 115)
	ebitenutil.DebugPrintAt(screen, "ENTER/R - Restart", 100, 140)
	ebitenutil.DebugPrintAt(screen, "M/ESC - Main Menu", 100, 160)
}

// drawSettingsScreen renders the settings screen
func (g *RoboGame) drawSettingsScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{60, 60, 60, 255}) // Gray background
	
	ebitenutil.DebugPrintAt(screen, "SETTINGS", 120, 50)
	ebitenutil.DebugPrintAt(screen, "========", 120, 65)
	ebitenutil.DebugPrintAt(screen, "Audio: ON", 100, 100)
	ebitenutil.DebugPrintAt(screen, "Resolution: 320x240", 100, 120)
	ebitenutil.DebugPrintAt(screen, "Difficulty: Normal", 100, 140)
	
	ebitenutil.DebugPrintAt(screen, "ESC/BACKSPACE - Back", 80, 180)
}

// drawTransitionScreen renders transition effects
func (g *RoboGame) drawTransitionScreen(screen *ebiten.Image) {
	stateManager := g.GetStateManager()
	progress := stateManager.GetTransitionProgress()
	
	// Draw the previous state as background
	switch stateManager.GetPreviousState() {
	case engine.StateLoading:
		g.drawLoadingScreen(screen)
	case engine.StateMenu:
		g.drawMenuScreen(screen)
	case engine.StatePlaying:
		g.drawGameScreen(screen)
	case engine.StatePaused:
		g.drawPausedScreen(screen)
	case engine.StateGameOver:
		g.drawGameOverScreen(screen)
	case engine.StateSettings:
		g.drawSettingsScreen(screen)
	}
	
	// Add transition effect (simple fade)
	alpha := uint8(progress * 255)
	g.overlayImage.Fill(color.RGBA{0, 0, 0, alpha})
	screen.DrawImage(g.overlayImage, nil)
	
	// Show transition progress for debugging
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Transitioning... %.1f%%", progress*100), 10, 220)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ROBO-9 Platformer")

	game := NewRoboGame()
	
	// Load assets before starting the game
	if err := game.LoadAssets(); err != nil {
		log.Fatalf("Failed to load assets: %v", err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
