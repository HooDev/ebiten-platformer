package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"ebiten-platformer/engine"
)

// RoboGame extends the base engine.Game with platformer-specific logic
type RoboGame struct {
	*engine.Game
	playerImage *ebiten.Image
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
	
	return &RoboGame{
		Game: baseGame,
	}
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

	// Set game state to playing after assets are loaded
	g.SetState(engine.StatePlaying)
	
	log.Println("All assets loaded successfully")
	return nil
}

// Update implements ebiten.Game interface
func (g *RoboGame) Update() error {
	// Call base game update
	if err := g.Game.Update(); err != nil {
		return err
	}

	// Handle different game states
	switch g.GetState() {
	case engine.StateLoading:
		// Assets loading is handled in main()
	case engine.StatePlaying:
		// Game logic will go here
	}

	return nil
}

// Draw implements ebiten.Game interface
func (g *RoboGame) Draw(screen *ebiten.Image) {
	// Call base game draw
	g.Game.Draw(screen)

	switch g.GetState() {
	case engine.StateLoading:
		screen.Fill(color.RGBA{0, 0, 0, 255})
		ebitenutil.DebugPrint(screen, "Loading assets...")
	case engine.StatePlaying:
		// Clear screen with sky blue
		screen.Fill(color.RGBA{135, 206, 235, 255})
		ebitenutil.DebugPrint(screen, "ROBO-9 Platformer\nAsset Loading System Active!")

		// Draw player if loaded
		if g.playerImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(100, 100)
			screen.DrawImage(g.playerImage, op)
		}

		// Display asset manager stats
		assetManager := g.GetAssetManager()
		stats := "Assets loaded:\n"
		stats += "Images: " + string(rune(assetManager.GetLoadedImageCount() + '0')) + "\n"
		stats += "Audio: " + string(rune(assetManager.GetLoadedAudioCount() + '0'))
		ebitenutil.DebugPrintAt(screen, stats, 200, 10)
	}
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
