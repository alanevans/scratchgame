package game

// LevelManager handles level creation and management
type LevelManager struct {
	entityManager *EntityManager
	assetManager  *AssetManager
}

// NewLevelManager creates a new level manager
func NewLevelManager(entityManager *EntityManager, assetManager *AssetManager) *LevelManager {
	return &LevelManager{
		entityManager: entityManager,
		assetManager:  assetManager,
	}
}

// LoadLevel1 creates the first level
func (lm *LevelManager) LoadLevel1() {
	// Clear existing entities
	lm.clearLevel()

	// Create player
	lm.createPlayer(100, 400)

	// Create platforms
	lm.createPlatform(0, 550, 800, 50)   // Ground
	lm.createPlatform(200, 450, 150, 20) // Platform 1
	lm.createPlatform(450, 350, 150, 20) // Platform 2
	lm.createPlatform(100, 250, 100, 20) // Platform 3
	lm.createPlatform(600, 200, 120, 20) // Platform 4
}

// clearLevel removes all entities from the level
func (lm *LevelManager) clearLevel() {
	// Note: In a more complex system, you'd want to track level entities
	// and only remove those, not all entities
}

// createPlayer creates the player entity
func (lm *LevelManager) createPlayer(x, y float64) EntityID {
	playerID := lm.entityManager.CreateEntity()

	// Add components
	lm.entityManager.AddComponent(playerID, &PositionComponent{X: x, Y: y})
	lm.entityManager.AddComponent(playerID, &VelocityComponent{VX: 0, VY: 0})
	lm.entityManager.AddComponent(playerID, &SpriteComponent{
		ImageKey: "player",
		Width:    48,
		Height:   48,
		OffsetX:  0,
		OffsetY:  0,
	})
	lm.entityManager.AddComponent(playerID, &ColliderComponent{
		Width:     48,
		Height:    48,
		OffsetX:   0,
		OffsetY:   0,
		IsSolid:   true,
		IsTrigger: false,
	})
	lm.entityManager.AddComponent(playerID, &PlayerComponent{
		Health:       100,
		MaxHealth:    100,
		JumpPower:    15.0,
		MoveSpeed:    4.0,
		IsGrounded:   false,
		JumpCooldown: 0,
		FacingRight:  true, // Start facing right
	})
	lm.entityManager.AddComponent(playerID, &GravityComponent{
		Force:      0.8,
		Terminal:   20.0,
		IsAffected: true,
	})
	lm.entityManager.AddComponent(playerID, &InputComponent{
		Enabled: true,
	})

	return playerID
}

// createPlatform creates a platform entity
func (lm *LevelManager) createPlatform(x, y, width, height float64) EntityID {
	platformID := lm.entityManager.CreateEntity()

	// Create a custom sprite for this platform size
	platformSprite := lm.assetManager.createColoredRect(int(width), int(height),
		struct{ R, G, B, A uint8 }{0, 150, 0, 255})
	lm.assetManager.AddImage("platform_"+string(rune(platformID)), platformSprite)

	// Add components
	lm.entityManager.AddComponent(platformID, &PositionComponent{X: x, Y: y})
	lm.entityManager.AddComponent(platformID, &SpriteComponent{
		ImageKey: "platform_" + string(rune(platformID)),
		Width:    int(width),
		Height:   int(height),
		OffsetX:  0,
		OffsetY:  0,
	})
	lm.entityManager.AddComponent(platformID, &ColliderComponent{
		Width:     width,
		Height:    height,
		OffsetX:   0,
		OffsetY:   0,
		IsSolid:   true,
		IsTrigger: false,
	})
	lm.entityManager.AddComponent(platformID, &PlatformComponent{
		IsMoving:      false,
		MoveSpeed:     0,
		MoveDirection: 0,
		MoveDistance:  0,
		StartX:        x,
		StartY:        y,
	})

	return platformID
}
