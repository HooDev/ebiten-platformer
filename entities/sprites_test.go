package entities

import (
	"testing"
)

func TestCreateTestSpriteSheet(t *testing.T) {
	spriteSheet := CreateTestSpriteSheet()

	if spriteSheet == nil {
		t.Fatal("CreateTestSpriteSheet returned nil")
	}

	bounds := spriteSheet.Bounds()
	expectedWidth := 192  // 6 frames * 32 pixels
	expectedHeight := 96  // 3 rows * 32 pixels

	if bounds.Dx() != expectedWidth {
		t.Errorf("Expected sprite sheet width %d, got %d", expectedWidth, bounds.Dx())
	}

	if bounds.Dy() != expectedHeight {
		t.Errorf("Expected sprite sheet height %d, got %d", expectedHeight, bounds.Dy())
	}
}

func TestCreateTestSpriteSheet_Consistency(t *testing.T) {
	// Generate multiple sprite sheets and verify they're identical
	sheet1 := CreateTestSpriteSheet()
	sheet2 := CreateTestSpriteSheet()

	if sheet1 == nil || sheet2 == nil {
		t.Fatal("CreateTestSpriteSheet returned nil")
	}

	bounds1 := sheet1.Bounds()
	bounds2 := sheet2.Bounds()

	if bounds1 != bounds2 {
		t.Error("Generated sprite sheets should have identical dimensions")
	}

	// Test that both sheets have the same size
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		t.Error("Generated sprite sheets should have identical dimensions")
	}
}

func TestCreateTestSpriteSheet_NonEmpty(t *testing.T) {
	spriteSheet := CreateTestSpriteSheet()

	if spriteSheet == nil {
		t.Fatal("CreateTestSpriteSheet returned nil")
	}

	bounds := spriteSheet.Bounds()

	// Verify the sprite sheet has reasonable dimensions
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Error("Sprite sheet should have positive dimensions")
	}

	// Verify it's large enough to contain the expected frames
	minExpectedWidth := 32   // At least one frame
	minExpectedHeight := 32  // At least one row

	if bounds.Dx() < minExpectedWidth || bounds.Dy() < minExpectedHeight {
		t.Errorf("Sprite sheet too small: %dx%d, expected at least %dx%d", 
			bounds.Dx(), bounds.Dy(), minExpectedWidth, minExpectedHeight)
	}
}

func TestCreateTestSpriteSheet_FrameCompatibility(t *testing.T) {
	spriteSheet := CreateTestSpriteSheet()
	
	if spriteSheet == nil {
		t.Fatal("CreateTestSpriteSheet returned nil")
	}

	bounds := spriteSheet.Bounds()
	frameSize := 32

	// Verify dimensions are multiples of frame size
	if bounds.Dx()%frameSize != 0 {
		t.Errorf("Sprite sheet width %d is not a multiple of frame size %d", bounds.Dx(), frameSize)
	}

	if bounds.Dy()%frameSize != 0 {
		t.Errorf("Sprite sheet height %d is not a multiple of frame size %d", bounds.Dy(), frameSize)
	}

	// Calculate expected frame counts
	framesPerRow := bounds.Dx() / frameSize
	totalRows := bounds.Dy() / frameSize

	if framesPerRow <= 0 {
		t.Error("Should have at least one frame per row")
	}

	if totalRows <= 0 {
		t.Error("Should have at least one row")
	}

	// Test that we can create animations with this sprite sheet
	// (This tests compatibility with the animation system)
	anim := NewAnimation(spriteSheet, frameSize, frameSize, framesPerRow, 1.0, true)
	if anim == nil {
		t.Error("Should be able to create animation from generated sprite sheet")
	}

	// Test getting a frame doesn't panic
	frame := anim.GetCurrentFrame()
	if frame == nil {
		t.Error("Should be able to get frame from animation using generated sprite sheet")
	}
}
