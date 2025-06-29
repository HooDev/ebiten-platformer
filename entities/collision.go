package entities

// CollisionResult represents the result of a collision check
// This mirrors the structure from the level package to avoid circular imports
type CollisionResult struct {
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

// CollisionChecker interface for objects that can check collisions
type CollisionChecker interface {
	CheckCollision(entityX, entityY, entityWidth, entityHeight float64) *CollisionResult
}
