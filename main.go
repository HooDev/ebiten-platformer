package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"ebiten-platformer/engine"
	"ebiten-platformer/entities"
	"ebiten-platformer/level"
)

// RoboGame extends the base engine.Game with platformer-specific logic
type RoboGame struct {
	*engine.Game
	playerImage    *ebiten.Image
	overlayImage   *ebiten.Image
	player         *entities.Player
	inputHandler   *entities.InputHandler
	currentLevel   *level.Level
	levelAdapter   *level.CollisionAdapter
	deltaTime      float64
	lastUpdateTime float64
}

// NewRoboGame creates a new platformer game instance
func NewRoboGame() *RoboGame {
	config := engine.GameConfig{
		ScreenWidth:  480,
		ScreenHeight: 360,
		AssetConfig: engine.AssetConfig{
			AssetDir:    "assets",
			UseEmbedded: false,
		},
	}

	baseGame := engine.NewGame(config)
	
	roboGame := &RoboGame{
		Game:         baseGame,
		overlayImage: ebiten.NewImage(480, 360),
		deltaTime:    1.0 / 60.0, // Initialize with 60 FPS
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
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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
	
	// Try to load player image, fallback to test sprite sheet if not available
	playerImg, err := assetManager.LoadImage("player.png")
	if err != nil {
		log.Printf("Could not load player.png, using test sprite sheet: %v", err)
		// Create a test sprite sheet for development
		playerImg = entities.CreateTestSpriteSheet()
	}
	g.playerImage = playerImg

	// Create level
	g.currentLevel = level.CreateSimpleLevel()
	g.levelAdapter = level.NewCollisionAdapter(g.currentLevel)

	// Create player entity
	g.player = entities.NewPlayer(100, 200, playerImg) // Start higher up
	
	// Connect player with level for collision detection
	g.player.SetLevel(g.levelAdapter)
	
	// Create input handler
	g.inputHandler = entities.NewInputHandler(g.player)

	// Set state to menu after assets are loaded
	g.SetState(engine.StateMenu)
	
	log.Println("All assets loaded successfully")
	log.Printf("Level created: %dx%d tiles, tile size: %d", g.currentLevel.Width, g.currentLevel.Height, g.currentLevel.TileSize)
	return nil
}

// Update implements ebiten.Game interface
func (g *RoboGame) Update() error {
	// Calculate delta time
	currentTime := float64(ebiten.CurrentTPS())
	if currentTime > 0 {
		g.deltaTime = 1.0 / currentTime
	}
	
	// Update player and input when playing
	stateManager := g.GetStateManager()
	if stateManager.GetCurrentState() == engine.StatePlaying {
		if g.inputHandler != nil {
			g.inputHandler.Update()
		}
		if g.player != nil {
			g.player.Update(g.deltaTime)
		}
	}
	
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
	ebitenutil.DebugPrintAt(screen, "ROBO-9 PLATFORMER", 160, 80)
	ebitenutil.DebugPrintAt(screen, "================", 160, 95)
	
	// Menu options
	ebitenutil.DebugPrintAt(screen, "ENTER - Start Game", 170, 140)
	ebitenutil.DebugPrintAt(screen, "S - Settings", 170, 160)
	
	// Game Controls
	ebitenutil.DebugPrintAt(screen, "Game Controls:", 20, 220)
	ebitenutil.DebugPrintAt(screen, "WASD/Arrow Keys - Move", 20, 240)
	ebitenutil.DebugPrintAt(screen, "Space/W/Up - Jump", 20, 260)
	ebitenutil.DebugPrintAt(screen, "C - Toggle Climb (Debug)", 20, 280)
	ebitenutil.DebugPrintAt(screen, "X - Test Damage (Debug)", 20, 300)
	
	// System Controls
	ebitenutil.DebugPrintAt(screen, "ESC - Pause | M - Menu", 20, 330)
}

// drawGameScreen renders the main game
func (g *RoboGame) drawGameScreen(screen *ebiten.Image) {
	// Clear screen with sky blue
	screen.Fill(color.RGBA{135, 206, 235, 255})
	
	// Draw level first (background)
	if g.currentLevel != nil {
		g.currentLevel.Draw(screen)
	}
	
	// Draw player on top of level
	if g.player != nil {
		g.player.Draw(screen)
		
		// Debug info
		x, y := g.player.GetPosition()
		vx, vy := g.player.GetVelocity()
		animState := g.player.GetAnimationState()
		
		debugInfo := fmt.Sprintf("Player: (%.1f, %.1f)\nVelocity: (%.1f, %.1f)\nOn Ground: %v\nFacing Right: %v\nAnimation: %v\nLevel: %s", 
			x, y, vx, vy, g.player.IsOnGround(), g.player.IsFacingRight(), animState, g.currentLevel.Name)
		ebitenutil.DebugPrintAt(screen, debugInfo, 10, 90)
	}
	
	// Game title and info
	ebitenutil.DebugPrint(screen, "ROBO-9 Platformer - PLAYING (Tile-Based Collision)")
	ebitenutil.DebugPrintAt(screen, "ESC: Pause | M: Menu | G: Game Over", 10, 20)
	ebitenutil.DebugPrintAt(screen, "Controls: WASD/Arrows to move, Space/W/Up to jump", 10, 35)
	ebitenutil.DebugPrintAt(screen, "Debug: C to toggle climb, X to test damage", 10, 50)
	
	// Display asset manager stats
	assetManager := g.GetAssetManager()
	stats := fmt.Sprintf("Assets loaded:\nImages: %d\nAudio: %d", 
		assetManager.GetLoadedImageCount(), 
		assetManager.GetLoadedAudioCount())
	ebitenutil.DebugPrintAt(screen, stats, 320, 200)
}

// drawPausedScreen renders the pause overlay
func (g *RoboGame) drawPausedScreen(screen *ebiten.Image) {
	// Draw the game screen first (as background)
	g.drawGameScreen(screen)
	
	// Add semi-transparent overlay using reusable image
	g.overlayImage.Fill(color.RGBA{0, 0, 0, 128}) // Semi-transparent black
	screen.DrawImage(g.overlayImage, nil)
	
	// Pause text
	ebitenutil.DebugPrintAt(screen, "GAME PAUSED", 190, 150)
	ebitenutil.DebugPrintAt(screen, "===========", 190, 165)
	ebitenutil.DebugPrintAt(screen, "ESC - Resume", 180, 190)
	ebitenutil.DebugPrintAt(screen, "M - Main Menu", 180, 210)
}

// drawGameOverScreen renders the game over screen
func (g *RoboGame) drawGameOverScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{60, 20, 20, 255}) // Dark red background
	
	ebitenutil.DebugPrintAt(screen, "GAME OVER", 190, 150)
	ebitenutil.DebugPrintAt(screen, "=========", 190, 165)
	ebitenutil.DebugPrintAt(screen, "ENTER/R - Restart", 170, 190)
	ebitenutil.DebugPrintAt(screen, "M/ESC - Main Menu", 170, 210)
}

// drawSettingsScreen renders the settings screen
func (g *RoboGame) drawSettingsScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{60, 60, 60, 255}) // Gray background
	
	ebitenutil.DebugPrintAt(screen, "SETTINGS", 200, 80)
	ebitenutil.DebugPrintAt(screen, "========", 200, 95)
	ebitenutil.DebugPrintAt(screen, "Audio: ON", 180, 140)
	ebitenutil.DebugPrintAt(screen, "Resolution: 480x360", 180, 160)
	ebitenutil.DebugPrintAt(screen, "Difficulty: Normal", 180, 180)
	
	ebitenutil.DebugPrintAt(screen, "ESC/BACKSPACE - Back", 160, 240)
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Transitioning... %.1f%%", progress*100), 10, 340)
}

func main() {
	ebiten.SetWindowSize(960, 720)
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
