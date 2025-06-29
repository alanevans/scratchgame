package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputManager handles all input processing
type InputManager struct {
	previousKeys map[ebiten.Key]bool
	currentKeys  map[ebiten.Key]bool
}

// NewInputManager creates a new input manager
func NewInputManager() *InputManager {
	return &InputManager{
		previousKeys: make(map[ebiten.Key]bool),
		currentKeys:  make(map[ebiten.Key]bool),
	}
}

// Update updates the input state
func (im *InputManager) Update() {
	// Copy current keys to previous
	for key, pressed := range im.currentKeys {
		im.previousKeys[key] = pressed
	}
	
	// Clear current keys
	for key := range im.currentKeys {
		delete(im.currentKeys, key)
	}
	
	// Get all pressed keys
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if ebiten.IsKeyPressed(key) {
			im.currentKeys[key] = true
		}
	}
}

// IsKeyPressed returns true if the key is currently being held down
func (im *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return im.currentKeys[key]
}

// IsKeyJustPressed returns true if the key was just pressed this frame
func (im *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return im.currentKeys[key] && !im.previousKeys[key]
}

// IsKeyJustReleased returns true if the key was just released this frame
func (im *InputManager) IsKeyJustReleased(key ebiten.Key) bool {
	return !im.currentKeys[key] && im.previousKeys[key]
}

// GetPressedKeys returns all currently pressed keys
func (im *InputManager) GetPressedKeys() []ebiten.Key {
	var keys []ebiten.Key
	for key, pressed := range im.currentKeys {
		if pressed {
			keys = append(keys, key)
		}
	}
	return keys
}

// Mouse input methods

// GetMousePosition returns the current mouse position
func (im *InputManager) GetMousePosition() (int, int) {
	return ebiten.CursorPosition()
}

// IsMouseButtonPressed returns true if the mouse button is currently being held down
func (im *InputManager) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

// IsMouseButtonJustPressed returns true if the mouse button was just pressed this frame
func (im *InputManager) IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

// IsMouseButtonJustReleased returns true if the mouse button was just released this frame
func (im *InputManager) IsMouseButtonJustReleased(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}
