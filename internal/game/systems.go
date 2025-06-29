package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// System interface that all systems must implement
type System interface {
	Update(em *EntityManager, inputManager *InputManager) error
}

// MovementSystem handles entity movement
type MovementSystem struct{}

func (s *MovementSystem) Update(em *EntityManager, inputManager *InputManager) error {
	entities := em.GetEntitiesWithComponents(ComponentTypePosition, ComponentTypeVelocity)
	
	for _, entityID := range entities {
		posComp, _ := em.GetComponent(entityID, ComponentTypePosition)
		velComp, _ := em.GetComponent(entityID, ComponentTypeVelocity)
		
		pos := posComp.(*PositionComponent)
		vel := velComp.(*VelocityComponent)
		
		// Update position based on velocity
		pos.X += vel.VX
		pos.Y += vel.VY
	}
	
	return nil
}

// GravitySystem applies gravity to entities
type GravitySystem struct{}

func (s *GravitySystem) Update(em *EntityManager, inputManager *InputManager) error {
	entities := em.GetEntitiesWithComponents(ComponentTypeVelocity, ComponentTypeGravity)
	
	for _, entityID := range entities {
		velComp, _ := em.GetComponent(entityID, ComponentTypeVelocity)
		gravComp, _ := em.GetComponent(entityID, ComponentTypeGravity)
		
		vel := velComp.(*VelocityComponent)
		grav := gravComp.(*GravityComponent)
		
		if grav.IsAffected {
			// Apply gravity
			vel.VY += grav.Force
			
			// Apply terminal velocity
			if vel.VY > grav.Terminal {
				vel.VY = grav.Terminal
			}
		}
	}
	
	return nil
}

// InputSystem handles player input
type InputSystem struct{}

func (s *InputSystem) Update(em *EntityManager, inputManager *InputManager) error {
	entities := em.GetEntitiesWithComponents(ComponentTypePlayer, ComponentTypeVelocity, ComponentTypeInput)
	
	for _, entityID := range entities {
		playerComp, _ := em.GetComponent(entityID, ComponentTypePlayer)
		velComp, _ := em.GetComponent(entityID, ComponentTypeVelocity)
		inputComp, _ := em.GetComponent(entityID, ComponentTypeInput)
		
		player := playerComp.(*PlayerComponent)
		vel := velComp.(*VelocityComponent)
		input := inputComp.(*InputComponent)
		
		if !input.Enabled {
			continue
		}
		
		// Handle horizontal movement
		if inputManager.IsKeyPressed(ebiten.KeyA) || inputManager.IsKeyPressed(ebiten.KeyArrowLeft) {
			vel.VX = -player.MoveSpeed
		} else if inputManager.IsKeyPressed(ebiten.KeyD) || inputManager.IsKeyPressed(ebiten.KeyArrowRight) {
			vel.VX = player.MoveSpeed
		} else {
			vel.VX = 0
		}
		
		// Handle jumping
		if (inputManager.IsKeyPressed(ebiten.KeySpace) || inputManager.IsKeyPressed(ebiten.KeyW) || inputManager.IsKeyPressed(ebiten.KeyArrowUp)) && player.IsGrounded && player.JumpCooldown <= 0 {
			vel.VY = -player.JumpPower
			player.IsGrounded = false
			player.JumpCooldown = 10 // Prevent rapid jumping
		}
		
		// Update jump cooldown
		if player.JumpCooldown > 0 {
			player.JumpCooldown--
		}
	}
	
	return nil
}

// CollisionSystem handles collision detection and resolution
type CollisionSystem struct{}

func (s *CollisionSystem) Update(em *EntityManager, inputManager *InputManager) error {
	entities := em.GetEntitiesWithComponents(ComponentTypePosition, ComponentTypeCollider)
	
	// Check collisions between all entities with colliders
	for i, entityA := range entities {
		for j, entityB := range entities {
			if i >= j {
				continue // Skip self and duplicate pairs
			}
			
			if s.checkCollision(em, entityA, entityB) {
				s.resolveCollision(em, entityA, entityB)
			}
		}
	}
	
	return nil
}

func (s *CollisionSystem) checkCollision(em *EntityManager, entityA, entityB EntityID) bool {
	posA, _ := em.GetComponent(entityA, ComponentTypePosition)
	posB, _ := em.GetComponent(entityB, ComponentTypePosition)
	colliderA, _ := em.GetComponent(entityA, ComponentTypeCollider)
	colliderB, _ := em.GetComponent(entityB, ComponentTypeCollider)
	
	pA := posA.(*PositionComponent)
	pB := posB.(*PositionComponent)
	cA := colliderA.(*ColliderComponent)
	cB := colliderB.(*ColliderComponent)
	
	// Calculate actual collision bounds
	leftA := pA.X + cA.OffsetX
	rightA := leftA + cA.Width
	topA := pA.Y + cA.OffsetY
	bottomA := topA + cA.Height
	
	leftB := pB.X + cB.OffsetX
	rightB := leftB + cB.Width
	topB := pB.Y + cB.OffsetY
	bottomB := topB + cB.Height
	
	// AABB collision detection
	return leftA < rightB && rightA > leftB && topA < bottomB && bottomA > topB
}

func (s *CollisionSystem) resolveCollision(em *EntityManager, entityA, entityB EntityID) {
	// Get components
	posA, _ := em.GetComponent(entityA, ComponentTypePosition)
	posB, _ := em.GetComponent(entityB, ComponentTypePosition)
	colliderA, _ := em.GetComponent(entityA, ComponentTypeCollider)
	colliderB, _ := em.GetComponent(entityB, ComponentTypeCollider)
	
	pA := posA.(*PositionComponent)
	pB := posB.(*PositionComponent)
	cA := colliderA.(*ColliderComponent)
	cB := colliderB.(*ColliderComponent)
	
	// Check if either entity is a player
	playerA := em.HasComponent(entityA, ComponentTypePlayer)
	playerB := em.HasComponent(entityB, ComponentTypePlayer)
	platformA := em.HasComponent(entityA, ComponentTypePlatform)
	platformB := em.HasComponent(entityB, ComponentTypePlatform)
	
	// Only resolve if one is a player and one is a platform, and both are solid
	if (!cA.IsSolid || !cB.IsSolid) {
		return
	}
	
	var playerEntity EntityID
	var playerPos *PositionComponent
	var playerCollider *ColliderComponent
	var platformPos *PositionComponent
	var platformCollider *ColliderComponent
	
	if playerA && platformB {
		playerEntity = entityA
		playerPos = pA
		playerCollider = cA
		platformPos = pB
		platformCollider = cB
	} else if playerB && platformA {
		playerEntity = entityB
		playerPos = pB
		playerCollider = cB
		platformPos = pA
		platformCollider = cA
	} else {
		return // No player-platform collision
	}
	
	// Calculate overlap
	playerLeft := playerPos.X + playerCollider.OffsetX
	playerRight := playerLeft + playerCollider.Width
	playerTop := playerPos.Y + playerCollider.OffsetY
	playerBottom := playerTop + playerCollider.Height
	
	platformLeft := platformPos.X + platformCollider.OffsetX
	platformRight := platformLeft + platformCollider.Width
	platformTop := platformPos.Y + platformCollider.OffsetY
	platformBottom := platformTop + platformCollider.Height
	
	// Calculate overlap amounts
	overlapX := math.Min(playerRight-platformLeft, platformRight-playerLeft)
	overlapY := math.Min(playerBottom-platformTop, platformBottom-playerTop)
	
	// Get player velocity and player component
	velComp, hasVel := em.GetComponent(playerEntity, ComponentTypeVelocity)
	playerComp, hasPlayer := em.GetComponent(playerEntity, ComponentTypePlayer)
	
	if !hasVel || !hasPlayer {
		return
	}
	
	vel := velComp.(*VelocityComponent)
	player := playerComp.(*PlayerComponent)
	
	// Resolve collision based on smallest overlap
	if overlapX < overlapY {
		// Horizontal collision
		if playerPos.X < platformPos.X {
			playerPos.X = platformLeft - playerCollider.Width - playerCollider.OffsetX
		} else {
			playerPos.X = platformRight - playerCollider.OffsetX
		}
		vel.VX = 0
	} else {
		// Vertical collision
		if playerPos.Y < platformPos.Y {
			// Player is above platform
			playerPos.Y = platformTop - playerCollider.Height - playerCollider.OffsetY
			vel.VY = 0
			player.IsGrounded = true
		} else {
			// Player is below platform
			playerPos.Y = platformBottom - playerCollider.OffsetY
			vel.VY = 0
		}
	}
}

// RenderSystem handles drawing entities
type RenderSystem struct{}

func (s *RenderSystem) Draw(screen *ebiten.Image, em *EntityManager, assetManager *AssetManager) {
	entities := em.GetEntitiesWithComponents(ComponentTypePosition, ComponentTypeSprite)
	
	for _, entityID := range entities {
		posComp, _ := em.GetComponent(entityID, ComponentTypePosition)
		spriteComp, _ := em.GetComponent(entityID, ComponentTypeSprite)
		
		pos := posComp.(*PositionComponent)
		sprite := spriteComp.(*SpriteComponent)
		
		image := assetManager.GetImage(sprite.ImageKey)
		if image != nil {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(pos.X+sprite.OffsetX, pos.Y+sprite.OffsetY)
			screen.DrawImage(image, options)
		}
	}
}
