package level

import (
	"testing"
)

func TestNewLevel(t *testing.T) {
	level := NewLevel(10, 8, 32, "Test")
	
	if level.Width != 10 {
		t.Errorf("Expected width 10, got %d", level.Width)
	}
	
	if level.Height != 8 {
		t.Errorf("Expected height 8, got %d", level.Height)
	}
	
	if level.TileSize != 32 {
		t.Errorf("Expected tile size 32, got %d", level.TileSize)
	}
	
	if level.Name != "Test" {
		t.Errorf("Expected name 'Test', got '%s'", level.Name)
	}
	
	// Check that all tiles are initialized as empty
	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			tile := level.GetTile(x, y)
			if tile.Type != TileEmpty {
				t.Errorf("Expected empty tile at (%d,%d), got %v", x, y, tile.Type)
			}
		}
	}
}

func TestSetAndGetTile(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	
	// Set a solid tile
	level.SetTile(2, 3, TileSolid)
	tile := level.GetTile(2, 3)
	
	if tile.Type != TileSolid {
		t.Errorf("Expected solid tile, got %v", tile.Type)
	}
	
	if !tile.IsSolid() {
		t.Error("Solid tile should report as solid")
	}
	
	// Test out of bounds (should return empty tile)
	outOfBounds := level.GetTile(-1, -1)
	if outOfBounds.Type != TileEmpty {
		t.Errorf("Out of bounds should return empty tile, got %v", outOfBounds.Type)
	}
	
	outOfBounds2 := level.GetTile(10, 10)
	if outOfBounds2.Type != TileEmpty {
		t.Errorf("Out of bounds should return empty tile, got %v", outOfBounds2.Type)
	}
}

func TestGetTileAtWorldPos(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	level.SetTile(1, 2, TileSolid)
	
	// World position (40, 70) should map to tile (1, 2)
	tile := level.GetTileAtWorldPos(40, 70)
	if tile.Type != TileSolid {
		t.Errorf("Expected solid tile at world pos (40, 70), got %v", tile.Type)
	}
	
	// Test edge of tile
	tile2 := level.GetTileAtWorldPos(32, 64)
	if tile2.Type != TileSolid {
		t.Errorf("Expected solid tile at world pos (32, 64), got %v", tile2.Type)
	}
}

func TestCollisionDetection(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	
	// Create a solid tile at (1, 1)
	level.SetTile(1, 1, TileSolid)
	
	// Test entity overlapping with solid tile
	result := level.CheckCollision(30, 30, 32, 32) // Overlaps tile (1,1)
	if !result.Collided {
		t.Error("Should detect collision with solid tile")
	}
	
	// Test entity not overlapping
	result2 := level.CheckCollision(100, 100, 32, 32) // No overlap
	if result2.Collided {
		t.Error("Should not detect collision when not overlapping")
	}
	
	// Test ground detection
	level.SetTile(2, 3, TileSolid)
	result3 := level.CheckCollision(64, 93, 32, 5) // Entity overlaps slightly with tile (2,3)
	
	// Debug: tile (2,3) is at world position (64, 96) to (96, 128)
	// Entity at (64, 93) with size (32, 5) has bottom at (64, 98)
	// So entity overlaps with tile by 2 pixels
	
	if !result3.Collided {
		t.Error("Should detect collision with tile when entities overlap")
	}
	
	if !result3.OnGround {
		t.Errorf("Should detect being on ground when standing on solid tile. Entity: (64,93) size (32,5), Tile (2,3) at world pos (64,96). Collided: %v", result3.Collided)
	}
}

func TestOneWayPlatform(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	level.SetTile(1, 1, TileOneWay)
	
	// Test entity coming from above (should collide)
	// Platform is at tile (1,1) = world pos (32,32)
	// Entity bottom at Y=30 with height=4 puts bottom at Y=34, which is 2 pixels below tile top (32)
	result := level.CheckCollision(32, 30, 32, 4) // Adjusted to work with GroundTolerance=2.0
	if !result.OneWayPlatform {
		t.Error("Should detect one-way platform collision from above")
	}
	
	// Test entity coming from below (should not collide significantly)
	result2 := level.CheckCollision(32, 40, 32, 32) // Below and overlapping
	if result2.OnGround {
		t.Error("Should not be on ground when coming from below one-way platform")
	}
}

func TestClimbableSurface(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	level.SetTile(1, 1, TileClimbable)
	
	result := level.CheckCollision(30, 30, 32, 32)
	if !result.ClimbableSurface {
		t.Error("Should detect climbable surface")
	}
	
	// Test that climbable tiles are also solid
	if !result.Collided {
		t.Error("Climbable tiles should also be solid")
	}
}

func TestDangerousTile(t *testing.T) {
	level := NewLevel(5, 5, 32, "Test")
	level.SetTile(1, 1, TileSpike)
	
	result := level.CheckCollision(30, 30, 32, 32)
	if !result.DangerousTile {
		t.Error("Should detect dangerous tile")
	}
	
	// Spikes should not be solid
	if result.Collided && result.CollisionX || result.CollisionY {
		t.Error("Spike tiles should not block movement")
	}
}

func TestTileProperties(t *testing.T) {
	// Test solid tile
	solidTile := NewTile(TileSolid, 0, 0)
	if !solidTile.IsSolid() || solidTile.IsClimbable() || solidTile.IsDangerous() {
		t.Error("Solid tile properties incorrect")
	}
	
	// Test climbable tile
	climbTile := NewTile(TileClimbable, 0, 0)
	if !climbTile.IsSolid() || !climbTile.IsClimbable() || climbTile.IsDangerous() {
		t.Error("Climbable tile properties incorrect")
	}
	
	// Test spike tile
	spikeTile := NewTile(TileSpike, 0, 0)
	if spikeTile.IsSolid() || spikeTile.IsClimbable() || !spikeTile.IsDangerous() {
		t.Error("Spike tile properties incorrect")
	}
	
	// Test one-way tile
	oneWayTile := NewTile(TileOneWay, 0, 0)
	if !oneWayTile.IsOneWay() || !oneWayTile.IsSolid() || oneWayTile.IsClimbable() {
		t.Error("One-way tile properties incorrect")
	}
}
