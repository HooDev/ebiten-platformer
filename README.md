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

- 2D platformer mechanics
- Sprite-based graphics
- Simple physics implementation
- Game state management

## Prerequisites

- Go 1.16 or higher
- Basic understanding of Go programming

## Getting Started

To run the game:

```bash
go run main.go
```

## Learning Resources

This project demonstrates various game development concepts including:

- Game loops and rendering
- Input handling
- Collision detection
- Animation and sprite management
- Game state management

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
