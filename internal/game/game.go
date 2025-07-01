package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents the main game state and implements ebiten.Game interface
type Game struct {
	width  int
	height int

	// Game state management
	stateManager *StateManager

	// Input manager
	inputManager *InputManager

	// Asset manager
	assetManager *AssetManager

	// Entity manager (for ECS pattern)
	entityManager *EntityManager
}

// NewGame creates a new game instance
func NewGame(width, height int) *Game {
	g := &Game{
		width:  width,
		height: height,
	}

	// Initialize managers
	g.inputManager = NewInputManager()
	g.assetManager = NewAssetManager()
	g.entityManager = NewEntityManager()
	g.stateManager = NewStateManager(g)

	// Load initial assets
	g.assetManager.LoadAssets()

	// Set initial state to playing
	g.stateManager.ChangeState(StateTypePlaying)

	return g
}

// Update is called every tick (60 FPS by default)
func (g *Game) Update() error {
	// Update input
	g.inputManager.Update()

	// Update current state
	g.stateManager.Update()

	return nil
}

// Draw is called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Delegate drawing to current state
	g.stateManager.Draw(screen)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

// Getters for managers (used by states)
func (g *Game) GetInputManager() *InputManager {
	return g.inputManager
}

func (g *Game) GetAssetManager() *AssetManager {
	return g.assetManager
}

func (g *Game) GetEntityManager() *EntityManager {
	return g.entityManager
}
