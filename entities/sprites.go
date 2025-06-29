package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// CreateTestSpriteSheet creates a simple test sprite sheet for the player
// This is a temporary solution until proper art assets are available
func CreateTestSpriteSheet() *ebiten.Image {
	// Create a sprite sheet with 6 frames horizontally, 3 rows vertically (18 frames total)
	frameWidth := 32
	frameHeight := 32
	cols := 6
	rows := 3
	
	spriteSheet := ebiten.NewImage(frameWidth*cols, frameHeight*rows)
	
	// Define colors for different animation states
	colors := []color.RGBA{
		{100, 149, 237, 255}, // Cornflower blue - idle frames (0-3)
		{100, 149, 237, 255}, // Same blue
		{100, 149, 237, 255},
		{100, 149, 237, 255},
		{72, 209, 204, 255},  // Medium turquoise - walk frames (4-7)
		{72, 209, 204, 255},
		{255, 215, 0, 255},   // Gold - jump frames (8-9)
		{255, 215, 0, 255},
		{255, 140, 0, 255},   // Dark orange - fall frames (10-11)
		{255, 140, 0, 255},
		{144, 238, 144, 255}, // Light green - climb frames (12-15)
		{144, 238, 144, 255},
		{144, 238, 144, 255},
		{144, 238, 144, 255},
		{220, 20, 60, 255},   // Crimson - damage frames (16-17)
		{220, 20, 60, 255},
		{128, 128, 128, 255}, // Gray - unused frames
		{128, 128, 128, 255},
	}
	
	// Fill each frame with its color and add a simple robot-like shape
	for i := 0; i < len(colors) && i < cols*rows; i++ {
		x := (i % cols) * frameWidth
		y := (i / cols) * frameHeight
		
		// Create frame background
		frameImg := ebiten.NewImage(frameWidth, frameHeight)
		frameImg.Fill(colors[i])
		
		// Add simple robot details
		addRobotDetails(frameImg, frameWidth, frameHeight, i)
		
		// Draw frame to sprite sheet
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		spriteSheet.DrawImage(frameImg, op)
	}
	
	return spriteSheet
}

// addRobotDetails adds simple geometric shapes to make it look more robot-like
func addRobotDetails(frame *ebiten.Image, width, height, frameIndex int) {
	// Create simple robot features
	
	// Head (darker rectangle at top)
	head := ebiten.NewImage(width-8, 8)
	head.Fill(color.RGBA{0, 0, 0, 100})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(4, 2)
	frame.DrawImage(head, op)
	
	// Eyes (small white squares)
	eye1 := ebiten.NewImage(2, 2)
	eye1.Fill(color.RGBA{255, 255, 255, 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(8, 4)
	frame.DrawImage(eye1, op)
	
	eye2 := ebiten.NewImage(2, 2)
	eye2.Fill(color.RGBA{255, 255, 255, 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(22, 4)
	frame.DrawImage(eye2, op)
	
	// Body (darker rectangle in middle)
	body := ebiten.NewImage(width-6, height-16)
	body.Fill(color.RGBA{0, 0, 0, 80})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(3, 10)
	frame.DrawImage(body, op)
	
	// Add animation-specific details
	switch {
	case frameIndex >= 4 && frameIndex <= 7: // Walk animation
		// Add slight offset for walking motion
		offset := float64((frameIndex-4)%2) * 2 - 1
		op.GeoM.Translate(offset, 0)
		
	case frameIndex >= 8 && frameIndex <= 9: // Jump animation
		// Make slightly smaller (compressed for jump)
		op.GeoM.Scale(1.0, 0.9)
		
	case frameIndex >= 16 && frameIndex <= 17: // Damage animation
		// Add red tint effect by drawing a semi-transparent red overlay
		overlay := ebiten.NewImage(width, height)
		overlay.Fill(color.RGBA{255, 0, 0, 60})
		frame.DrawImage(overlay, &ebiten.DrawImageOptions{})
	}
}
