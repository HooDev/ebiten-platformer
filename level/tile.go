package level

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TileType represents different types of tiles
type TileType int

const (
	TileEmpty TileType = iota
	TileSolid
	TileClimbable
	TileSpike
	TileOneWay // Platform you can jump through from below
)

// Tile represents a single tile in the level
type Tile struct {
	Type     TileType
	X, Y     int     // Grid coordinates
	Solid    bool    // Whether the tile blocks movement
	Climbable bool   // Whether the player can climb on this tile
	Sprite   *ebiten.Image // Visual representation (optional)
}

// NewTile creates a new tile with the given type and position
func NewTile(tileType TileType, x, y int) *Tile {
	tile := &Tile{
		Type: tileType,
		X:    x,
		Y:    y,
	}
	
	// Set properties based on tile type
	switch tileType {
	case TileEmpty:
		tile.Solid = false
		tile.Climbable = false
	case TileSolid:
		tile.Solid = true
		tile.Climbable = false
	case TileClimbable:
		tile.Solid = true
		tile.Climbable = true
	case TileSpike:
		tile.Solid = false
		tile.Climbable = false
	case TileOneWay:
		tile.Solid = true // Special handling in collision detection
		tile.Climbable = false
	}
	
	return tile
}

// GetBounds returns the world-space bounds of this tile
func (t *Tile) GetBounds(tileSize int) (x, y, width, height float64) {
	return float64(t.X * tileSize), float64(t.Y * tileSize), float64(tileSize), float64(tileSize)
}

// IsClimbable returns whether this tile can be climbed
func (t *Tile) IsClimbable() bool {
	return t.Climbable
}

// IsSolid returns whether this tile blocks movement
func (t *Tile) IsSolid() bool {
	return t.Solid
}

// IsDangerous returns whether this tile damages the player
func (t *Tile) IsDangerous() bool {
	return t.Type == TileSpike
}

// IsOneWay returns whether this is a one-way platform
func (t *Tile) IsOneWay() bool {
	return t.Type == TileOneWay
}
