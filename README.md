# 2D Platform Game

A simple 2D platform game written in Go using the Ebitengine library.

## Features

- **Entity-Component-System (ECS) Architecture**: Clean separation of data and logic
- **Game State Manager**: Menu, Playing, Paused, and Game Over states
- **Physics System**: Gravity, collision detection, and platformer movement
- **Modular Design**: Easy to extend with new entities, components, and systems

## Game Design Patterns Used

1. **Entity-Component-System (ECS)**: Entities are composed of components, and systems operate on entities with specific component combinations
2. **State Pattern**: Game states (Menu, Playing, Paused, Game Over) with clean transitions
3. **Manager Pattern**: Separate managers for input, assets, entities, and levels
4. **Observer Pattern**: Systems respond to component changes

## Controls

- **WASD** or **Arrow Keys**: Move player
- **Space**: Jump
- **Escape**: Pause/unpause game
- **Q**: Quit to menu (from pause screen)

## Project Structure

```
scratchgame/
├── main.go                           # Entry point
├── internal/game/
│   ├── game.go                       # Main game loop
│   ├── state_manager.go              # Game state management
│   ├── entity_manager.go             # ECS entity management
│   ├── components.go                 # Component definitions
│   ├── systems.go                    # System implementations
│   ├── input_manager.go              # Input handling
│   ├── asset_manager.go              # Asset loading/management
│   ├── states.go                     # Game state implementations
│   ├── level_manager.go              # Level creation/management
│   └── assets/                       # Game assets directory
└── go.mod                            # Go module definition
```

## Architecture Overview

### Entity-Component-System (ECS)

- **Entities**: Unique IDs representing game objects
- **Components**: Data structures (Position, Velocity, Sprite, Collider, etc.)
- **Systems**: Logic that operates on entities with specific components

### Game States

- **MenuState**: Main menu screen
- **PlayingState**: Active gameplay
- **PausedState**: Game paused overlay
- **GameOverState**: Game over screen

### Systems

- **InputSystem**: Handles player input for movement and jumping
- **GravitySystem**: Applies gravity to entities
- **MovementSystem**: Updates positions based on velocity
- **CollisionSystem**: Detects and resolves collisions
- **RenderSystem**: Draws entities to the screen

## Building and Running

```bash
# Install dependencies
go mod tidy

# Build the game
go build -o scratchgame .

# Run the game
./scratchgame
```

## Extending the Game

The modular architecture makes it easy to add new features:

1. **New Components**: Add to `components.go`
2. **New Systems**: Add to `systems.go` and register in playing state
3. **New Entities**: Create in level manager or add entity factory methods
4. **New Game States**: Implement GameState interface and register in state manager
5. **New Assets**: Add to assets directory and load in asset manager

## Dependencies

- [Ebitengine](https://github.com/hajimehoshi/ebiten) - 2D game engine for Go

## License

See LICENSE file for details.
