package level

// CreateTestLevel creates a simple test level for development
func CreateTestLevel() *Level {
	level := NewLevel(30, 20, 32, "Test Level")
	
	// Create a ground platform across the bottom
	for x := 0; x < level.Width; x++ {
		level.SetTile(x, level.Height-1, TileSolid) // Bottom row
		level.SetTile(x, level.Height-2, TileSolid) // Second bottom row for thickness
	}
	
	// Create some platforms
	// Left platform
	for x := 5; x <= 10; x++ {
		level.SetTile(x, 15, TileSolid)
	}
	
	// Middle floating platform
	for x := 15; x <= 20; x++ {
		level.SetTile(x, 12, TileSolid)
	}
	
	// High platform on the right
	for x := 25; x <= 28; x++ {
		level.SetTile(x, 8, TileSolid)
	}
	
	// Add some climbable walls
	for y := 10; y <= 17; y++ {
		level.SetTile(12, y, TileClimbable) // Climbable wall between platforms
	}
	
	for y := 5; y <= 12; y++ {
		level.SetTile(22, y, TileClimbable) // Climbable wall to high platform
	}
	
	// Add a one-way platform
	for x := 8; x <= 12; x++ {
		level.SetTile(x, 10, TileOneWay)
	}
	
	// Add some spike hazards
	level.SetTile(18, level.Height-3, TileSpike)
	level.SetTile(19, level.Height-3, TileSpike)
	
	return level
}

// CreateSimpleLevel creates a very basic level for initial testing
func CreateSimpleLevel() *Level {
	level := NewLevel(20, 15, 32, "Simple Level")
	
	// Just create a simple ground
	for x := 0; x < level.Width; x++ {
		level.SetTile(x, level.Height-1, TileSolid)
	}
	
	// Add a small platform to test jumping
	for x := 8; x <= 12; x++ {
		level.SetTile(x, 10, TileSolid)
	}
	
	return level
}
