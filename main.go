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
				X:        100,
				Y:        400,
				Width:    32,
				Height:   32,
				OnGround: false,
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
	store *Store
}

// gameReducer handles state changes based on actions
func gameReducer(state GameState, action Action) GameState {
	newState := state

	switch action.Type {
	case MOVE_LEFT:
		newState.Player.VelX = -moveSpeed
	case MOVE_RIGHT:
		newState.Player.VelX = moveSpeed
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

		// Simple ground collision (ground at Y=500)
		if newState.Player.Y >= 500 {
			newState.Player.Y = 500
			newState.Player.VelY = 0
			newState.Player.OnGround = true
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
	platforms := []struct{ x, y, w, h float64 }{
		{200, 400, 100, 20},
		{400, 300, 100, 20},
		{600, 350, 100, 20},
		{800, 250, 100, 20},
	}

	for _, platform := range platforms {
		ebitenutil.DrawRect(screen, platform.x-state.CameraX, platform.y, platform.w, platform.h, color.RGBA{139, 69, 19, 255})
	}

	// Draw player (simple red rectangle)
	playerScreenX := state.Player.X - state.CameraX
	ebitenutil.DrawRect(screen, playerScreenX, state.Player.Y, state.Player.Width, state.Player.Height, color.RGBA{255, 0, 0, 255})

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
	// Create store with reducer
	store := NewStore(gameReducer)

	// Create game instance
	game := &Game{
		store: store,
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
