package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// StateType represents different game states
type StateType int

const (
	StateTypeMenu StateType = iota
	StateTypePlaying
	StateTypePaused
	StateTypeGameOver
)

// GameState interface that all states must implement
type GameState interface {
	Enter()
	Exit()
	Update(game *Game) error
	Draw(screen *ebiten.Image, game *Game)
}

// StateManager manages game state transitions
type StateManager struct {
	game         *Game
	currentState GameState
	states       map[StateType]GameState
}

// NewStateManager creates a new state manager
func NewStateManager(game *Game) *StateManager {
	sm := &StateManager{
		game:   game,
		states: make(map[StateType]GameState),
	}

	// Register all game states
	sm.states[StateTypeMenu] = NewMenuState()
	sm.states[StateTypePlaying] = NewPlayingState()
	sm.states[StateTypePaused] = NewPausedState()
	sm.states[StateTypeGameOver] = NewGameOverState()

	return sm
}

// ChangeState transitions to a new state
func (sm *StateManager) ChangeState(stateType StateType) {
	if sm.currentState != nil {
		sm.currentState.Exit()
	}

	sm.currentState = sm.states[stateType]
	sm.currentState.Enter()
}

// Update the current state
func (sm *StateManager) Update() error {
	if sm.currentState != nil {
		return sm.currentState.Update(sm.game)
	}
	return nil
}

// Draw the current state
func (sm *StateManager) Draw(screen *ebiten.Image) {
	if sm.currentState != nil {
		sm.currentState.Draw(screen, sm.game)
	}
}
