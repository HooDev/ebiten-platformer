# Ebitengine Platformer

A 2D platformer game built with [Ebitengine](https://ebitengine.org/) and Go.

## About

This project serves as a learning opportunity to explore game development using Go and the Ebitengine game library. Ebitengine (formerly known as Ebiten) is a simple 2D game library that provides a straightforward way to develop games in Go.

## Game concept

### Story
You play as ROBO-9, a compassionate repair robot in a post-apocalyptic world where all the cats have become sad and dejected. Your mission is to restore joy to the world by collecting energy hearts scattered throughout the dangerous landscape and delivering them to the melancholic felines you encounter.

### Core Gameplay
This is a 2D side-scrolling platformer where exploration and kindness are rewarded over violence. As ROBO-9, you must:

- **Collect Energy Hearts**: These glowing power sources are hidden throughout each level, often in hard-to-reach places or behind environmental puzzles
- **Help Sad Cats**: When you find a dejected cat, you can give it an energy heart to restore its happiness, earning points and sometimes unlocking new areas
- **Navigate Hazardous Terrain**: Jump across platforms, climb walls with your magnetic grippers, and solve environmental puzzles to progress

### Obstacles & Challenges
The world is filled with dangers that will test your platforming skills:

- **Killer Drones**: Autonomous security bots that patrol set patterns - time your movements to avoid their detection
- **Falling Debris**: Unstable structures and falling objects that require quick reflexes to dodge
- **Spike Traps**: Hidden and visible spikes that can damage your circuits
- **Energy Barriers**: Electrified walls that require you to find power switches or alternative routes
- **Moving Platforms**: Timing-based challenges that test your jumping precision
- **Environmental Hazards**: Acid pools, steam vents, and other industrial dangers

### Robot Abilities
ROBO-9 comes equipped with several helpful capabilities:

- **Double Jump**: Use your built-in thrusters for enhanced mobility
- **Wall Climbing**: Magnetic feet allow you to scale certain metallic surfaces
- **Heart Scanner**: Detect nearby energy hearts and sad cats through walls
- **Emergency Shield**: Brief invincibility when taking damage (limited uses per level)

### Progression System
- **Score System**: Earn points by helping cats, collecting hearts, and completing levels quickly
- **Cat Collection**: Keep track of all the cats you've helped in your digital journal
- **Level Unlocks**: Some areas only open after helping a certain number of cats
- **Heart Efficiency**: Bonus points for completing levels while using the minimum number of hearts needed


## Features

- **Complete Animation System**: Multi-state sprite animations with smooth transitions
- **ROBO-9 Player Character**: Fully animated robot with movement, jumping, climbing, and damage states
- **Robust Collision System**: Tile-based collision detection with binary search for sub-pixel precision
- **Physics-Based Movement**: Gravity, friction, and swept collision-based player movement
- **Coyote Time**: Forgiving jump mechanics allowing players to jump briefly after leaving platforms
- **Flexible Input System**: Keyboard controls with multiple key bindings
- **Game State Management**: Menu, playing, paused, and game over states
- **Asset Management**: Efficient sprite loading and management
- **Test Framework**: Built-in test sprite generation for development
- **Delta-Time Animation**: Smooth, framerate-independent animations

## Prerequisites

- Go 1.16 or higher
- Basic understanding of Go programming
- Ebitengine v2 (automatically managed via Go modules)

## Getting Started

### Building and Running

To build the game:

```bash
go build -o robo-platformer
```

To run the game directly:

```bash
go run main.go
```

### WSL Environment

When running the game in WSL (Windows Subsystem for Linux), you need to set the target OS to Windows:

```bash
GOOS=windows go run main.go
```

Or use the provided convenience script:

```bash
./run.wsl.sh
```

### Controls

| Action | Primary Keys | Alternative Keys |
|--------|-------------|------------------|
| Move Left | Left Arrow | A |
| Move Right | Right Arrow | D |
| Jump | Space | Up Arrow, W |
| Climb Up | Up Arrow | W (when near climbable surface) |
| Climb Down | Down Arrow | S (when climbing) |
| Pause | Escape | - |
| Menu | M | - |

#### Debug Controls (Development)
- **C**: Toggle climbing mode
- **X**: Test damage state

## Project Structure

```
├── main.go                 # Game entry point and main game loop
├── engine/                 # Core game engine components
│   ├── game.go            # Base game state management
│   └── assets.go          # Asset loading and management
├── entities/              # Game entities and components
│   ├── player.go          # ROBO-9 player implementation
│   ├── collision.go       # Collision detection interfaces
│   ├── animation.go       # Animation system
│   ├── input.go           # Input handling
│   └── sprites.go         # Test sprite generation
├── level/                 # Level system and tile-based collision
│   ├── level.go           # Level implementation with tiles
│   ├── tile.go            # Tile definitions and properties
│   ├── adapter.go         # Collision adapter for entities
│   └── test_levels.go     # Test level generation
├── assets/                # Game assets (sprites, audio, etc.)
│   └── player.png         # Player sprite sheet (192x96px)
└── docs/                  # Documentation
    ├── development-plan.md         # Project roadmap
    ├── collision-system.md         # Collision system developer guide
    ├── animation-system.md         # Animation system docs
    ├── player-implementation.md    # Player entity docs
    ├── coyote-time.md             # Coyote time implementation guide
    └── robo9-sprite-specification.md # Sprite requirements
```

## Development Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Development Plan](docs/development-plan.md)**: Complete project roadmap and feature timeline
- **[Collision System](docs/collision-system.md)**: Tile-based collision detection developer guide
- **[Animation System](docs/animation-system.md)**: Technical details of the animation framework
- **[Player Implementation](docs/player-implementation.md)**: ROBO-9 character implementation guide
- **[Coyote Time](docs/coyote-time.md)**: Forgiving jump mechanics implementation guide
- **[Sprite Specification](docs/robo9-sprite-specification.md)**: Detailed sprite sheet requirements and layout


## Learning Resources

This project demonstrates various game development concepts including:

- **Game Loops and Rendering**: Main game loop, frame timing, and rendering pipeline
- **Input Handling**: Keyboard input processing and player control systems
- **Animation Systems**: Sprite-based animation with state management and timing
- **Physics Simulation**: Basic gravity, collision detection, and movement physics
- **Coyote Time Mechanics**: Forgiving jump timing for improved player experience
- **Game State Management**: Menu systems, pause functionality, and state transitions
- **Collision Detection**: Robust tile-based collision with binary search precision
- **Asset Management**: Loading and organizing game resources efficiently
- **Entity-Component Patterns**: Modular game object design and architecture

### Technical Highlights

- **Delta-Time Animation**: Framerate-independent animation timing
- **State Machine Design**: Clean separation of game states and transitions
- **Modular Architecture**: Easily extensible entity and component system
- **Performance Optimization**: Efficient sprite rendering and memory management

### Educational Value

The codebase serves as a practical example of:
- Go programming best practices for game development
- Ebitengine API usage and optimization techniques
- Game architecture patterns and design principles
- Testing strategies for interactive applications

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
