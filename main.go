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

	// Draw player (small black and white dog)
	playerScreenX := state.Player.X - state.CameraX
	drawDog(screen, playerScreenX, state.Player.Y)

	// Draw UI
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player: (%.1f, %.1f)", state.Player.X, state.Player.Y))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nVelocity: (%.1f, %.1f)", state.Player.VelX, state.Player.VelY))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nOnGround: %v", state.Player.OnGround))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nCamera: %.1f", state.CameraX))
	ebitenutil.DebugPrint(screen, "\n\nControls:")
	ebitenutil.DebugPrint(screen, "\nArrows/WASD: Move")
	ebitenutil.DebugPrint(screen, "\nSpace/Up/W: Jump")
}

// drawDog draws a simple pixel art dog using rectangles
func drawDog(screen *ebiten.Image, x, y float64) {
	// Dog body (white)
	ebitenutil.DrawRect(screen, x+8, y+16, 16, 12, color.RGBA{255, 255, 255, 255})

	// Dog head (white)
	ebitenutil.DrawRect(screen, x+4, y+8, 12, 12, color.RGBA{255, 255, 255, 255})

	// Dog ears (black)
	ebitenutil.DrawRect(screen, x+2, y+6, 4, 6, color.RGBA{0, 0, 0, 255})  // Left ear
	ebitenutil.DrawRect(screen, x+14, y+6, 4, 6, color.RGBA{0, 0, 0, 255}) // Right ear

	// Dog nose (black)
	ebitenutil.DrawRect(screen, x+8, y+12, 2, 2, color.RGBA{0, 0, 0, 255})

	// Dog eyes (black)
	ebitenutil.DrawRect(screen, x+6, y+10, 1, 1, color.RGBA{0, 0, 0, 255})  // Left eye
	ebitenutil.DrawRect(screen, x+13, y+10, 1, 1, color.RGBA{0, 0, 0, 255}) // Right eye

	// Dog legs (white)
	ebitenutil.DrawRect(screen, x+10, y+28, 3, 4, color.RGBA{255, 255, 255, 255}) // Front left leg
	ebitenutil.DrawRect(screen, x+15, y+28, 3, 4, color.RGBA{255, 255, 255, 255}) // Front right leg
	ebitenutil.DrawRect(screen, x+19, y+28, 3, 4, color.RGBA{255, 255, 255, 255}) // Back left leg
	ebitenutil.DrawRect(screen, x+22, y+28, 3, 4, color.RGBA{255, 255, 255, 255}) // Back right leg

	// Dog paws (black)
	ebitenutil.DrawRect(screen, x+10, y+30, 3, 2, color.RGBA{0, 0, 0, 255}) // Front left paw
	ebitenutil.DrawRect(screen, x+15, y+30, 3, 2, color.RGBA{0, 0, 0, 255}) // Front right paw
	ebitenutil.DrawRect(screen, x+19, y+30, 3, 2, color.RGBA{0, 0, 0, 255}) // Back left paw
	ebitenutil.DrawRect(screen, x+22, y+30, 3, 2, color.RGBA{0, 0, 0, 255}) // Back right paw

	// Dog tail (black)
	ebitenutil.DrawRect(screen, x+24, y+18, 4, 2, color.RGBA{0, 0, 0, 255})

	// Dog spots (black patches for pattern)
	ebitenutil.DrawRect(screen, x+12, y+18, 4, 4, color.RGBA{0, 0, 0, 255}) // Body spot
	ebitenutil.DrawRect(screen, x+16, y+10, 2, 3, color.RGBA{0, 0, 0, 255}) // Head spot
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
