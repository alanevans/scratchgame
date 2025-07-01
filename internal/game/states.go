package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// MenuState represents the main menu state
type MenuState struct{}

func NewMenuState() *MenuState {
	return &MenuState{}
}

func (s *MenuState) Enter() {
	// Initialize menu
}

func (s *MenuState) Exit() {
	// Cleanup menu
}

func (s *MenuState) Update(game *Game) error {
	// Check for input to start game
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeySpace) {
		game.stateManager.ChangeState(StateTypePlaying)
	}
	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image, game *Game) {
	screen.Fill(color.RGBA{30, 30, 60, 255})
	ebitenutil.DebugPrint(screen, "2D Platform Game\n\nPress SPACE to start\n\nControls:\nWASD or Arrow Keys to move\nSPACE to jump")
}

// PlayingState represents the main gameplay state
type PlayingState struct {
	systems      []System
	renderSystem *RenderSystem
	levelManager *LevelManager
	initialized  bool
}

func NewPlayingState() *PlayingState {
	return &PlayingState{
		systems: []System{
			&InputSystem{},
			&GravitySystem{},
			&MovementSystem{},
			&CollisionSystem{},
		},
		renderSystem: &RenderSystem{},
		initialized:  false,
	}
}

func (s *PlayingState) Enter() {
	// This will be called when entering the playing state
	s.initialized = false
}

func (s *PlayingState) Exit() {
	// Cleanup playing state
}

func (s *PlayingState) Update(game *Game) error {
	// Initialize level on first update
	if !s.initialized {
		s.levelManager = NewLevelManager(game.GetEntityManager(), game.GetAssetManager())
		s.levelManager.LoadLevel1()
		s.initialized = true
	}

	// Check for pause
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeyEscape) {
		game.stateManager.ChangeState(StateTypePaused)
		return nil
	}

	// Update all systems
	for _, system := range s.systems {
		if err := system.Update(game.GetEntityManager(), game.GetInputManager()); err != nil {
			return err
		}
	}

	return nil
}

func (s *PlayingState) Draw(screen *ebiten.Image, game *Game) {
	// Clear screen with background color
	screen.Fill(color.RGBA{135, 206, 235, 255}) // Sky blue

	// Draw all entities
	s.renderSystem.Draw(screen, game.GetEntityManager(), game.GetAssetManager())

	// Draw UI
	ebitenutil.DebugPrint(screen, "ESC to pause")
}

// PausedState represents the paused game state
type PausedState struct{}

func NewPausedState() *PausedState {
	return &PausedState{}
}

func (s *PausedState) Enter() {
	// Initialize pause menu
}

func (s *PausedState) Exit() {
	// Cleanup pause menu
}

func (s *PausedState) Update(game *Game) error {
	// Check for unpause
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeyEscape) {
		game.stateManager.ChangeState(StateTypePlaying)
	}

	// Check for return to menu
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeyQ) {
		game.stateManager.ChangeState(StateTypeMenu)
	}

	return nil
}

func (s *PausedState) Draw(screen *ebiten.Image, game *Game) {
	// Draw a semi-transparent overlay
	overlay := ebiten.NewImage(game.width, game.height)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen, "PAUSED\n\nESC to resume\nQ to quit to menu")
}

// GameOverState represents the game over state
type GameOverState struct{}

func NewGameOverState() *GameOverState {
	return &GameOverState{}
}

func (s *GameOverState) Enter() {
	// Initialize game over screen
}

func (s *GameOverState) Exit() {
	// Cleanup game over screen
}

func (s *GameOverState) Update(game *Game) error {
	// Check for restart
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeySpace) {
		game.stateManager.ChangeState(StateTypePlaying)
	}

	// Check for return to menu
	if game.GetInputManager().IsKeyJustPressed(ebiten.KeyQ) {
		game.stateManager.ChangeState(StateTypeMenu)
	}

	return nil
}

func (s *GameOverState) Draw(screen *ebiten.Image, game *Game) {
	screen.Fill(color.RGBA{60, 30, 30, 255})
	ebitenutil.DebugPrint(screen, "GAME OVER\n\nSPACE to restart\nQ to quit to menu")
}
