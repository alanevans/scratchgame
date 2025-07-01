package game

// PositionComponent represents an entity's position in 2D space
type PositionComponent struct {
	X, Y float64
}

func (c *PositionComponent) GetType() ComponentType {
	return ComponentTypePosition
}

// VelocityComponent represents an entity's velocity
type VelocityComponent struct {
	VX, VY float64
}

func (c *VelocityComponent) GetType() ComponentType {
	return ComponentTypeVelocity
}

// SpriteComponent represents visual rendering information
type SpriteComponent struct {
	ImageKey string
	Width    int
	Height   int
	OffsetX  float64
	OffsetY  float64
}

func (c *SpriteComponent) GetType() ComponentType {
	return ComponentTypeSprite
}

// ColliderComponent represents collision boundaries
type ColliderComponent struct {
	Width     float64
	Height    float64
	OffsetX   float64
	OffsetY   float64
	IsSolid   bool
	IsTrigger bool
}

func (c *ColliderComponent) GetType() ComponentType {
	return ComponentTypeCollider
}

// PlayerComponent marks an entity as the player
type PlayerComponent struct {
	Health       int
	MaxHealth    int
	JumpPower    float64
	MoveSpeed    float64
	IsGrounded   bool
	JumpCooldown int
	FacingRight  bool // true = facing right, false = facing left
}

func (c *PlayerComponent) GetType() ComponentType {
	return ComponentTypePlayer
}

// PlatformComponent marks an entity as a platform
type PlatformComponent struct {
	IsMoving      bool
	MoveSpeed     float64
	MoveDirection float64
	MoveDistance  float64
	StartX        float64
	StartY        float64
}

func (c *PlatformComponent) GetType() ComponentType {
	return ComponentTypePlatform
}

// GravityComponent applies gravity to an entity
type GravityComponent struct {
	Force      float64
	Terminal   float64
	IsAffected bool
}

func (c *GravityComponent) GetType() ComponentType {
	return ComponentTypeGravity
}

// InputComponent marks an entity as controllable by input
type InputComponent struct {
	Enabled bool
}

func (c *InputComponent) GetType() ComponentType {
	return ComponentTypeInput
}
