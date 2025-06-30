package level

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Constants for collision detection tuning
const (
	// GroundTolerance defines how close an entity must be to a tile surface
	// to be considered "on ground". This accounts for floating-point precision
	// and provides stable ground detection.
	GroundTolerance = 2.0
)

// Level represents a game level with tile-based collision
type Level struct {
	Width      int             // Level width in tiles
	Height     int             // Level height in tiles
	TileSize   int             // Size of each tile in pixels
	Tiles      [][]*Tile       // 2D array of tiles [y][x]
	Background *ebiten.Image   // Background image (optional)
	Name       string          // Level name
}

// NewLevel creates a new empty level
func NewLevel(width, height, tileSize int, name string) *Level {
	// Initialize 2D tile array
	tiles := make([][]*Tile, height)
	for y := 0; y < height; y++ {
		tiles[y] = make([]*Tile, width)
		for x := 0; x < width; x++ {
			tiles[y][x] = NewTile(TileEmpty, x, y)
		}
	}

	return &Level{
		Width:    width,
		Height:   height,
		TileSize: tileSize,
		Tiles:    tiles,
		Name:     name,
	}
}

// SetTile sets a tile at the given grid coordinates
func (l *Level) SetTile(x, y int, tileType TileType) {
	if l.IsValidCoord(x, y) {
		l.Tiles[y][x] = NewTile(tileType, x, y)
	}
}

// GetTile returns the tile at the given grid coordinates
func (l *Level) GetTile(x, y int) *Tile {
	if l.IsValidCoord(x, y) {
		return l.Tiles[y][x]
	}
	return NewTile(TileEmpty, x, y) // Return empty tile for out-of-bounds
}

// GetTileAtWorldPos returns the tile at the given world position
func (l *Level) GetTileAtWorldPos(worldX, worldY float64) *Tile {
	tileX := int(math.Floor(worldX / float64(l.TileSize)))
	tileY := int(math.Floor(worldY / float64(l.TileSize)))
	return l.GetTile(tileX, tileY)
}

// IsValidCoord checks if the given coordinates are within the level bounds
func (l *Level) IsValidCoord(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height
}

// GetWorldBounds returns the world-space bounds of the level
func (l *Level) GetWorldBounds() (width, height float64) {
	return float64(l.Width * l.TileSize), float64(l.Height * l.TileSize)
}

// CollisionResult represents the result of a collision check
type CollisionResult struct {
	Collided         bool
	CollisionX       bool  // True if collision occurred on X axis
	CollisionY       bool  // True if collision occurred on Y axis
	PenetrationX     float64 // How much the entity penetrated on X axis
	PenetrationY     float64 // How much the entity penetrated on Y axis
	OnGround         bool  // True if the entity is standing on solid ground
	TouchingWall     bool  // True if the entity is touching a wall
	ClimbableSurface bool // True if touching a climbable surface
	DangerousTile    bool // True if touching a dangerous tile
	OneWayPlatform   bool // True if touching a one-way platform from above
}

// CheckCollision checks collision between a rectangular entity and the level tiles
func (l *Level) CheckCollision(entityX, entityY, entityWidth, entityHeight float64) *CollisionResult {
	result := &CollisionResult{}

	// Calculate which tiles the entity overlaps
	leftTile := int(math.Floor(entityX / float64(l.TileSize)))
	rightTile := int(math.Floor((entityX + entityWidth - 1) / float64(l.TileSize)))
	topTile := int(math.Floor(entityY / float64(l.TileSize)))
	bottomTile := int(math.Floor((entityY + entityHeight - 1) / float64(l.TileSize)))
	
	// For ground detection, also check the tile directly below the entity
	belowTile := int(math.Floor((entityY + entityHeight) / float64(l.TileSize)))

	// Check all overlapping tiles
	for tileY := topTile; tileY <= bottomTile; tileY++ {
		for tileX := leftTile; tileX <= rightTile; tileX++ {
			tile := l.GetTile(tileX, tileY)
			
			if tile.Type == TileEmpty {
				continue
			}

			// Get tile bounds
			tileBounds := l.getTileBounds(tileX, tileY)
			
			// Check if entity actually overlaps with this tile
			if l.rectanglesOverlap(entityX, entityY, entityWidth, entityHeight,
				tileBounds.X, tileBounds.Y, tileBounds.Width, tileBounds.Height) {
				
				l.processCollision(result, tile, entityX, entityY, entityWidth, entityHeight, tileBounds)
			}
		}
	}
	
	// Additionally check for ground contact with tiles directly below the entity
	if belowTile > bottomTile { // Only if we're not already checking this tile row
		for tileX := leftTile; tileX <= rightTile; tileX++ {
			tile := l.GetTile(tileX, belowTile)
			
			if tile.Type == TileEmpty {
				continue
			}

			// Get tile bounds
			tileBounds := l.getTileBounds(tileX, belowTile)
			
			// Check for ground contact (entity bottom touching tile top)
			entityBottom := entityY + entityHeight
			tileTop := tileBounds.Y
			
			if entityBottom >= tileTop && entityBottom <= tileTop + GroundTolerance && 
				entityX < tileBounds.X + tileBounds.Width && entityX + entityWidth > tileBounds.X {
				
				// Process ground contact
				result.Collided = true
				if tile.IsSolid() || tile.IsOneWay() {
					result.OnGround = true
				}
				if tile.IsOneWay() {
					result.OneWayPlatform = true
				}
				if tile.IsClimbable() {
					result.ClimbableSurface = true
				}
				if tile.IsDangerous() {
					result.DangerousTile = true
				}
			}
		}
	}

	return result
}

// processCollision handles collision logic for a single tile
func (l *Level) processCollision(result *CollisionResult, tile *Tile, entityX, entityY, entityWidth, entityHeight float64, tileBounds struct{ X, Y, Width, Height float64 }) {
	result.Collided = true
	
	// Calculate overlap amounts for all cases
	overlapX := l.getOverlapX(entityX, entityWidth, tileBounds.X, tileBounds.Width)
	overlapY := l.getOverlapY(entityY, entityHeight, tileBounds.Y, tileBounds.Height)
	
	// Set flags based on tile type and position
	if tile.IsSolid() {
		// Check if entity is on top of the tile (ground check)
		entityBottom := entityY + entityHeight
		tileTop := tileBounds.Y
		
		// More precise ground detection: entity is on ground if:
		// 1. Bottom of entity is very close to the top of the tile
		// 2. The overlap is minimal (tolerance for floating point precision)
		isOnGround := entityBottom >= tileTop && entityBottom <= tileTop + GroundTolerance
		
		if isOnGround {
			result.OnGround = true
		}
		
		// Improved collision detection logic:
		// - If entity is on ground, prioritize ground contact over wall collision
		// - Only trigger wall collision if it's truly a side impact, not ground contact
		
		if isOnGround {
			// When on ground, we should be very lenient about horizontal collision
			// Only consider it a wall hit if:
			// 1. The entity is significantly embedded horizontally (more than half its width)
			// 2. AND the entity is not just resting on the platform edge
			
			// Check if this is a true wall collision (entity hitting side of a block)
			// rather than just standing on a platform edge
			entityCenterX := entityX + entityWidth/2
			
			// If entity center is significantly outside the tile bounds horizontally,
			// and horizontal penetration is more than 75% of entity width, it's a wall hit
			if overlapX > entityWidth * 0.75 && 
			   (entityCenterX < tileBounds.X + 4 || entityCenterX > tileBounds.X + tileBounds.Width - 4) {
				result.TouchingWall = true
				result.CollisionX = true
				result.PenetrationX = overlapX
			}
			// Otherwise, allow horizontal movement even with overlap (standing on platform edge)
		} else {
			// When not on ground, use normal collision detection
			// Check for wall collision (horizontal overlap is significant)
			if overlapX > 2 && (overlapX < overlapY || overlapY < 2) {
				result.TouchingWall = true
				result.CollisionX = true
				result.PenetrationX = overlapX
			}
			
			// Check for vertical collision (vertical overlap is significant)
			if overlapY > 2 && (overlapY < overlapX || overlapX < 2) {
				result.CollisionY = true
				result.PenetrationY = overlapY
			}
		}
	}
	
	// Handle one-way platforms (can only collide from above)
	if tile.IsOneWay() {
		entityBottom := entityY + entityHeight
		tileTop := tileBounds.Y
		
		// Only collide if entity is coming from above and touching the surface
		if entityBottom >= tileTop && entityBottom <= tileTop + GroundTolerance {
			result.OnGround = true
			result.OneWayPlatform = true
			result.CollisionY = true
			result.PenetrationY = entityBottom - tileTop
		}
	}
	
	// Check for climbable surfaces
	if tile.IsClimbable() {
		result.ClimbableSurface = true
	}
	
	// Check for dangerous tiles
	if tile.IsDangerous() {
		result.DangerousTile = true
	}
}

// Helper function to get tile bounds in world coordinates
func (l *Level) getTileBounds(tileX, tileY int) struct{ X, Y, Width, Height float64 } {
	return struct{ X, Y, Width, Height float64 }{
		X:      float64(tileX * l.TileSize),
		Y:      float64(tileY * l.TileSize),
		Width:  float64(l.TileSize),
		Height: float64(l.TileSize),
	}
}

// Helper function to check if two rectangles overlap
func (l *Level) rectanglesOverlap(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

// Helper function to calculate X-axis overlap
func (l *Level) getOverlapX(x1, w1, x2, w2 float64) float64 {
	return math.Min(x1+w1, x2+w2) - math.Max(x1, x2)
}

// Helper function to calculate Y-axis overlap
func (l *Level) getOverlapY(y1, h1, y2, h2 float64) float64 {
	return math.Min(y1+h1, y2+h2) - math.Max(y1, y2)
}

// Draw renders the level (basic tile visualization)
func (l *Level) Draw(screen *ebiten.Image) {
	// Draw background if available
	if l.Background != nil {
		screen.DrawImage(l.Background, &ebiten.DrawImageOptions{})
	}
	
	// Draw tiles (simple colored rectangles for now)
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			tile := l.Tiles[y][x]
			if tile.Type != TileEmpty {
				l.drawTile(screen, tile, x, y)
			}
		}
	}
}

// drawTile draws a single tile with a color based on its type
func (l *Level) drawTile(screen *ebiten.Image, tile *Tile, x, y int) {
	tileImg := ebiten.NewImage(l.TileSize, l.TileSize)
	
	// Color based on tile type
	switch tile.Type {
	case TileSolid:
		tileImg.Fill(color.RGBA{128, 128, 128, 255}) // Gray
	case TileClimbable:
		tileImg.Fill(color.RGBA{139, 69, 19, 255})   // Brown (like metal/wood)
	case TileSpike:
		tileImg.Fill(color.RGBA{255, 0, 0, 255})     // Red
	case TileOneWay:
		tileImg.Fill(color.RGBA{0, 255, 0, 255})     // Green
	}
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*l.TileSize), float64(y*l.TileSize))
	screen.DrawImage(tileImg, op)
}

// PlayerCollisionResult represents collision data in a format compatible with the player
type PlayerCollisionResult struct {
	Collided         bool
	CollisionX       bool
	CollisionY       bool
	PenetrationX     float64
	PenetrationY     float64
	OnGround         bool
	TouchingWall     bool
	ClimbableSurface bool
	DangerousTile    bool
	OneWayPlatform   bool
}

// CheckPlayerCollision checks collision and returns a player-compatible result
func (l *Level) CheckPlayerCollision(entityX, entityY, entityWidth, entityHeight float64) *PlayerCollisionResult {
	result := l.CheckCollision(entityX, entityY, entityWidth, entityHeight)
	
	return &PlayerCollisionResult{
		Collided:         result.Collided,
		CollisionX:       result.CollisionX,
		CollisionY:       result.CollisionY,
		PenetrationX:     result.PenetrationX,
		PenetrationY:     result.PenetrationY,
		OnGround:         result.OnGround,
		TouchingWall:     result.TouchingWall,
		ClimbableSurface: result.ClimbableSurface,
		DangerousTile:    result.DangerousTile,
		OneWayPlatform:   result.OneWayPlatform,
	}
}
