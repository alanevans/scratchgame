package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gravity      = 0.5
	jumpForce    = -12
	moveSpeed    = 3
)

// Platform represents a platform in the game
type Platform struct {
	X, Y, W, H float64
}

// Game platforms
var gamePlatforms = []Platform{
	{200, 400, 100, 20},
	{400, 300, 100, 20},
	{600, 350, 100, 20},
	{800, 250, 100, 20},
}

// Action types for Redux pattern
type ActionType string

const (
	MOVE_LEFT       ActionType = "MOVE_LEFT"
	MOVE_RIGHT      ActionType = "MOVE_RIGHT"
	JUMP            ActionType = "JUMP"
	APPLY_GRAVITY   ActionType = "APPLY_GRAVITY"
	UPDATE_POSITION ActionType = "UPDATE_POSITION"
)

// Action represents a game action
type Action struct {
	Type    ActionType
	Payload interface{}
}

// Player represents the player character
type Player struct {
	X, Y          float64
	VelX, VelY    float64
	Width, Height float64
	OnGround      bool
	FacingRight   bool // New field to track facing direction
}

// GameState represents the complete game state
type GameState struct {
	Player  Player
	CameraX float64
}

// Store manages the game state using Redux pattern
type Store struct {
	state   GameState
	reducer func(GameState, Action) GameState
}

// NewStore creates a new store with initial state
func NewStore(reducer func(GameState, Action) GameState) *Store {
	return &Store{
		state: GameState{
			Player: Player{
				X:           100,
				Y:           400,
				Width:       32,
				Height:      32,
				OnGround:    false,
				FacingRight: true, // Start facing right
			},
			CameraX: 0,
		},
		reducer: reducer,
	}
}

// Dispatch sends an action to the store
func (s *Store) Dispatch(action Action) {
	s.state = s.reducer(s.state, action)
}

// GetState returns the current state
func (s *Store) GetState() GameState {
	return s.state
}

// Game implements ebiten.Game interface
type Game struct {
	store          *Store
	dogSpriteRight *ebiten.Image
	dogSpriteLeft  *ebiten.Image
}

// gameReducer handles state changes based on actions
func gameReducer(state GameState, action Action) GameState {
	newState := state

	switch action.Type {
	case MOVE_LEFT:
		newState.Player.VelX = -moveSpeed
		newState.Player.FacingRight = false
	case MOVE_RIGHT:
		newState.Player.VelX = moveSpeed
		newState.Player.FacingRight = true
	case JUMP:
		if newState.Player.OnGround {
			newState.Player.VelY = jumpForce
			newState.Player.OnGround = false
		}
	case APPLY_GRAVITY:
		if !newState.Player.OnGround {
			newState.Player.VelY += gravity
		}
	case UPDATE_POSITION:
		// Update player position
		newState.Player.X += newState.Player.VelX
		newState.Player.Y += newState.Player.VelY

		// Platform collision detection
		playerBottom := newState.Player.Y + newState.Player.Height
		playerLeft := newState.Player.X
		playerRight := newState.Player.X + newState.Player.Width
		
		landed := false
		
		// Check collision with each platform
		for _, platform := range gamePlatforms {
			// Check if player is horizontally within platform bounds
			if playerRight > platform.X && playerLeft < platform.X+platform.W {
				// Check if player is falling and would land on platform
				if newState.Player.VelY > 0 && playerBottom >= platform.Y && playerBottom <= platform.Y+platform.H+5 {
					newState.Player.Y = platform.Y - newState.Player.Height
					newState.Player.VelY = 0
					newState.Player.OnGround = true
					landed = true
					break
				}
			}
		}

		// Ground collision (ground at Y=500) - only if not already landed on platform
		if !landed && newState.Player.Y >= 500 {
			newState.Player.Y = 500
			newState.Player.VelY = 0
			newState.Player.OnGround = true
			landed = true
		}

		// If not on ground or platform, player is in air
		if !landed {
			newState.Player.OnGround = false
		}

		// Reset horizontal velocity (no friction for now)
		newState.Player.VelX = 0

		// Update camera to follow player
		targetCameraX := newState.Player.X - screenWidth/2
		if targetCameraX < 0 {
			targetCameraX = 0
		}
		newState.CameraX = targetCameraX
	}

	return newState
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	// Handle input
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.store.Dispatch(Action{Type: MOVE_LEFT})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.store.Dispatch(Action{Type: MOVE_RIGHT})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.store.Dispatch(Action{Type: JUMP})
	}

	// Apply physics
	g.store.Dispatch(Action{Type: APPLY_GRAVITY})
	g.store.Dispatch(Action{Type: UPDATE_POSITION})

	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	state := g.store.GetState()

	// Clear screen with sky blue
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Draw ground
	groundY := 500 + 32 // Ground position + player height
	for x := 0; x < screenWidth+64; x += 32 {
		ebitenutil.DrawRect(screen, float64(x-int(state.CameraX)%32), float64(groundY), 32, float64(screenHeight-groundY), color.RGBA{34, 139, 34, 255})
	}

	// Draw simple platforms
	for _, platform := range gamePlatforms {
		ebitenutil.DrawRect(screen, platform.X-state.CameraX, platform.Y, platform.W, platform.H, color.RGBA{139, 69, 19, 255})
	}

	// Draw player (small black and white dog)
	playerScreenX := state.Player.X - state.CameraX
	var dogSprite *ebiten.Image
	if state.Player.FacingRight {
		dogSprite = g.dogSpriteRight
	} else {
		dogSprite = g.dogSpriteLeft
	}

	if dogSprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(playerScreenX, state.Player.Y)
		screen.DrawImage(dogSprite, op)
	}

	// Draw UI
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player: (%.1f, %.1f)", state.Player.X, state.Player.Y))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nVelocity: (%.1f, %.1f)", state.Player.VelX, state.Player.VelY))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nOnGround: %v", state.Player.OnGround))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nCamera: %.1f", state.CameraX))
	ebitenutil.DebugPrint(screen, "\n\nControls:")
	ebitenutil.DebugPrint(screen, "\nArrows/WASD: Move")
	ebitenutil.DebugPrint(screen, "\nSpace/Up/W: Jump")
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Load dog sprites (note: we swap the files because the flipping logic is reversed)
	dogSpriteRight, _, err := ebitenutil.NewImageFromFile("dog_sprite_left.png") // Use left file for right facing
	if err != nil {
		log.Fatal("Failed to load right dog sprite:", err)
	}

	dogSpriteLeft, _, err := ebitenutil.NewImageFromFile("dog_sprite_right.png") // Use right file for left facing
	if err != nil {
		log.Fatal("Failed to load left dog sprite:", err)
	}

	// Create store with reducer
	store := NewStore(gameReducer)

	// Create game instance
	game := &Game{
		store:          store,
		dogSpriteRight: dogSpriteRight,
		dogSpriteLeft:  dogSpriteLeft,
	}

	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Retro Platform Game - Redux Pattern")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
