package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewAnimation(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	anim := NewAnimation(img, 32, 32, 2, 1.0, true)

	if anim == nil {
		t.Fatal("NewAnimation returned nil")
	}

	if anim.FrameCount != 2 {
		t.Errorf("Expected FrameCount 2, got %d", anim.FrameCount)
	}

	if anim.FrameTime != 1.0 {
		t.Errorf("Expected FrameTime 1.0, got %f", anim.FrameTime)
	}

	if anim.Loop != true {
		t.Errorf("Expected Loop true, got %t", anim.Loop)
	}

	if anim.CurrentTime != 0 {
		t.Errorf("Expected CurrentTime 0, got %f", anim.CurrentTime)
	}

	if anim.Finished != false {
		t.Errorf("Expected Finished false, got %t", anim.Finished)
	}

	if len(anim.Frames) != 2 {
		t.Errorf("Expected 2 frames, got %d", len(anim.Frames))
	}
}

func TestAnimation_Update(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	anim := NewAnimation(img, 32, 32, 2, 0.5, true) // 0.5 seconds per frame

	// Test normal update
	initialTime := anim.CurrentTime
	deltaTime := 0.1
	anim.Update(deltaTime)
	
	if anim.CurrentTime != initialTime+deltaTime {
		t.Errorf("Expected CurrentTime %f, got %f", initialTime+deltaTime, anim.CurrentTime)
	}

	// Test frame advancement
	anim.Update(0.5) // Should advance to next frame
	if anim.CurrentTime < 0.5 {
		t.Error("Animation should have advanced")
	}

	// Test looping
	anim.Update(0.5) // Should complete loop and reset
	if anim.CurrentTime != 0 {
		t.Errorf("Looping animation should reset CurrentTime to 0, got %f", anim.CurrentTime)
	}
}

func TestAnimation_Reset(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	anim := NewAnimation(img, 32, 32, 2, 0.5, false)

	// Advance animation
	anim.Update(1.0)
	anim.Finished = true

	// Reset animation
	anim.Reset()

	if anim.CurrentTime != 0 {
		t.Errorf("Expected CurrentTime 0 after reset, got %f", anim.CurrentTime)
	}

	if anim.Finished != false {
		t.Errorf("Expected Finished false after reset, got %t", anim.Finished)
	}
}

func TestAnimation_GetCurrentFrame(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	anim := NewAnimation(img, 32, 32, 2, 1.0, true)

	frame := anim.GetCurrentFrame()
	if frame == nil {
		t.Fatal("GetCurrentFrame returned nil")
	}

	bounds := frame.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 32 {
		t.Errorf("Expected frame size 32x32, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestAnimation_NonLooping(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	anim := NewAnimation(img, 32, 32, 2, 0.5, false) // Non-looping

	// Play through entire animation
	anim.Update(1.0) // Total time = 2 frames * 0.5 = 1.0 second

	if !anim.IsFinished() {
		t.Error("Non-looping animation should be finished")
	}

	// Further updates should not change state
	prevTime := anim.CurrentTime
	anim.Update(0.5)
	if anim.CurrentTime != prevTime {
		t.Error("Finished non-looping animation should not advance time")
	}
}

func TestAnimation_EdgeCases(t *testing.T) {
	// Test with zero frame count
	img := ebiten.NewImage(32, 32)
	anim := NewAnimation(img, 32, 32, 0, 1.0, true)
	
	// Should not panic
	anim.Update(1.0)
	frame := anim.GetCurrentFrame()
	if frame != nil {
		t.Error("GetCurrentFrame should return nil for 0 frameCount")
	}

	// Test with very small sprite sheet
	smallImg := ebiten.NewImage(1, 1)
	smallAnim := NewAnimation(smallImg, 32, 32, 1, 1.0, true)
	
	// Should not panic
	smallAnim.Update(1.0)
	smallFrame := smallAnim.GetCurrentFrame()
	if smallFrame == nil {
		t.Error("GetCurrentFrame should return a frame even with small sprite sheet")
	}
}

func TestAnimation_FrameIndexing(t *testing.T) {
	img := ebiten.NewImage(96, 32)
	anim := NewAnimation(img, 32, 32, 3, 0.5, true) // 3 frames, 0.5s each

	// Test frame 0
	frame0 := anim.GetCurrentFrame()
	if frame0 == nil {
		t.Fatal("Frame 0 should not be nil")
	}

	// Advance to frame 1
	anim.Update(0.5)
	frame1 := anim.GetCurrentFrame()
	if frame1 == nil {
		t.Fatal("Frame 1 should not be nil")
	}

	// Advance to frame 2
	anim.Update(0.5)
	frame2 := anim.GetCurrentFrame()
	if frame2 == nil {
		t.Fatal("Frame 2 should not be nil")
	}

	// Should loop back to frame 0
	anim.Update(0.5)
	frameLooped := anim.GetCurrentFrame()
	if frameLooped == nil {
		t.Fatal("Looped frame should not be nil")
	}
}
