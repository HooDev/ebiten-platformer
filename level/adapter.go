package level

import (
	"ebiten-platformer/entities"
)

// CollisionAdapter adapts the Level to work with the entities.CollisionChecker interface
type CollisionAdapter struct {
	level *Level
}

// NewCollisionAdapter creates a new collision adapter for a level
func NewCollisionAdapter(level *Level) *CollisionAdapter {
	return &CollisionAdapter{level: level}
}

// CheckCollision implements the entities.CollisionChecker interface
func (ca *CollisionAdapter) CheckCollision(entityX, entityY, entityWidth, entityHeight float64) *entities.CollisionResult {
	result := ca.level.CheckCollision(entityX, entityY, entityWidth, entityHeight)
	
	return &entities.CollisionResult{
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
