# ROBO-9 Platformer Development Plan

## Project Overview

This document outlines the development roadmap for the ROBO-9 platformer game built with Ebitengine and Go. The plan is structured in phases to ensure steady progress and iterative testing.

## Development Phases

### Phase 1: Core Foundation (Weeks 1-2)
**Goal**: Establish basic game loop and player movement

#### 1.1 Project Setup
- [x] Initialise Go module and Ebitengine dependency
- [x] Create basic project structure
- [x] Set up asset loading system
- [x] Implement basic game state management

#### 1.2 Basic Player Entity
- [x] Create ROBO-9 sprite and animation system
- [x] Implement basic movement (left, right, jump)
- [x] Add collision detection with ground (basic implementation)
- [x] Wall climbing system (with debug controls)
- [x] Damage system with immunity frames
- [x] Animation state machine (6 states: idle, walk, jump, fall, climb, damage)
- [x] Comprehensive input handling (WASD + arrow keys)
- [x] Physics system (gravity, friction, velocity-based movement)	- [x] **Improve collision detection with proper tile-based system**
	- [x] **Fix jittery collision/ground detection issues**
	- [x] **Implement robust swept collision detection (prevents tunneling)**
	- [x] **Fix player sinking/getting stuck on platform edges**
	- [x] **Implement precise binary search collision (eliminates tunneling and variable sinking)**
- [x] **Add coyote time for more forgiving jumps**
- [ ] **Add jump buffering for responsive controls**
- [ ] **Basic camera following player**

#### 1.3 Level Framework
- [x] **Create simple tile-based level system**
- [x] **Implement static platforms**
- [x] **Robust collision detection with swept movement**
- [x] **Comprehensive collision regression tests**
- [x] **Collision system with binary search precision**
- [ ] **Basic level loading from data files**
- [ ] **Screen boundaries and camera constraints**

**Deliverable**: Playable character that can move and jump on basic platforms

### Phase 2: Game Mechanics (Weeks 3-4)
**Goal**: Implement core gameplay systems

#### 2.1 Enhanced Movement
- [ ] **Double jump ability**
- [x] Wall climbing on metallic surfaces (basic implementation complete)
- [x] **Robust physics and collision detection** (integrated with tile system)
- [x] Animation state machine for player (6 states implemented)
- [ ] **Variable jump height (hold for higher jumps)**

#### 2.2 Collectibles System
- [ ] **Energy heart entities**
- [ ] **Heart collection mechanics**
- [ ] **Visual feedback for collection**
- [ ] **Inventory/counter system**

#### 2.3 Cat Interaction System
- [ ] **Sad cat entities with basic AI**
- [ ] **Heart giving interaction**
- [ ] **Cat happiness state changes**
- [ ] **Point scoring system**

**Deliverable**: Core gameplay loop functional - collect hearts and help cats

### Phase 3: Obstacles & Challenges (Weeks 5-6)
**Goal**: Add danger and difficulty progression

#### 3.1 Basic Hazards
- [ ] **Spike traps (static)**
- [ ] **Falling debris system**
- [x] Damage system and health (basic damage state implemented)
- [ ] **Respawn/checkpoint system**
- [ ] **Health system with multiple hit points**

#### 3.2 Enemy Systems
- [ ] Killer drone entities
- [ ] Patrol AI patterns
- [ ] Detection and chase behaviour
- [ ] Collision damage

#### 3.3 Environmental Obstacles
- [ ] Moving platforms
- [ ] Energy barriers and switches
- [ ] Acid pools and steam vents
- [ ] Timing-based challenges

**Deliverable**: Challenging levels with various obstacles and enemies

### Phase 4: Advanced Features (Weeks 7-8)
**Goal**: Polish and advanced mechanics

#### 4.1 Robot Abilities
- [ ] Heart scanner implementation
- [ ] Emergency shield system
- [ ] Ability cooldowns and limitations
- [ ] Visual indicators for abilities

#### 4.2 Level Design Tools
- [ ] Level editor or data format
- [ ] Multiple level creation
- [ ] Level progression system
- [ ] Save/load game state

#### 4.3 Audio & Visual Polish
- [ ] Sound effects for actions
- [ ] Background music
- [ ] Particle effects
- [ ] Improved animations

**Deliverable**: Feature-complete game with multiple levels

### Phase 5: Content & Polish (Weeks 9-10)
**Goal**: Content creation and final polish

#### 5.1 Content Creation
- [ ] Design and implement 10-15 levels
- [ ] Progressive difficulty curve
- [ ] Hidden secrets and bonus areas
- [ ] Achievement/completion tracking

#### 5.2 UI/UX
- [ ] Main menu system
- [ ] In-game HUD (health, hearts, score)
- [ ] Pause menu
- [ ] Settings/options menu

#### 5.3 Final Polish
- [ ] Bug fixes and optimisation
- [ ] Balance tuning
- [ ] Performance optimisation
- [ ] Documentation completion

**Deliverable**: Polished, complete game ready for release

## Technical Architecture

### Core Systems

#### 1. Game Engine Structure
```
main.go                 // Entry point and game loop
├── engine/            // Core engine systems
│   ├── game.go        // Main game state manager
│   ├── scene.go       // Scene management
│   └── input.go       // Input handling
├── entities/          // Game objects
│   ├── player.go      // ROBO-9 implementation
│   ├── collision.go   // Collision interfaces
│   ├── cat.go         // Cat entities
│   ├── heart.go       // Energy hearts
│   └── enemies.go     // Drones and hazards
├── level/             // Level and collision system
│   ├── level.go       // Tile-based level implementation
│   ├── tile.go        // Tile definitions
│   ├── adapter.go     // Collision adapter
│   └── test_levels.go // Test level generation
├── systems/           // Game systems
│   ├── physics.go     // Collision and movement
│   ├── rendering.go   // Drawing and animation
│   └── audio.go       // Sound management
└── levels/            // Level data and loading
    ├── loader.go      // Level loading system
    └── data/          // Level definition files
```

#### 2. Asset Management
- Sprite sheets for animations
- Tileset for level construction
- Audio files for effects and music
- Level data in JSON or custom format

#### 3. Performance Considerations
- Efficient sprite batching
- Culling off-screen entities
- Optimized collision detection
- Memory management for large levels

## Asset Requirements

### Visual Assets
- [ ] ROBO-9 sprite sheet (idle, walk, jump, climb, damage states)
- [ ] Cat sprite sheet (sad, happy, idle animations)
- [ ] Energy heart sprite with glow animation
- [ ] Tileset for platforms and environment
- [ ] Drone enemies with patrol animations
- [ ] Environmental hazards (spikes, barriers, debris)
- [ ] UI elements and fonts
- [ ] Particle effect textures

### Audio Assets
- [ ] Background music tracks (2-3 ambient pieces)
- [ ] Sound effects:
  - Player movement (footsteps, jump, land)
  - Heart collection
  - Cat interaction (meow, purr)
  - Damage and death sounds
  - Environmental sounds (drone hums, electrical buzzes)

## Testing Strategy

### Unit Testing
- Physics calculations
- Collision detection algorithms
- Binary search collision precision
- Game state transitions
- Save/load functionality

### Integration Testing
- Complete gameplay loops
- Level progression
- Performance under load
- Cross-platform compatibility

### Playtesting
- Difficulty curve validation
- User experience feedback
- Bug identification
- Balance adjustments

## Risk Management

### Technical Risks
- **Ebitengine limitations**: Research library capabilities early
- **Performance issues**: Profile and optimise regularly
- **Asset pipeline**: Establish workflow early

### Design Risks
- **Scope creep**: Stick to core features for MVP
- **Difficulty balance**: Regular playtesting
- **Level design**: Create tools and templates

### Mitigation Strategies
- Regular milestone reviews
- Prototype risky features early
- Maintain backup plans for complex features
- Focus on core gameplay loop first

## Success Metrics

### Technical Goals
- Consistent 60 FPS on target hardware
- Sub-second level loading times
- Smooth, responsive controls
- Stable performance across platforms

### Gameplay Goals
- Engaging core loop (collect → help → progress)
- Clear difficulty progression
- Intuitive controls and mechanics
- Replayability through secrets and optimisation

## Timeline Summary

| Phase | Duration | Key Deliverable |
|-------|----------|----------------|
| 1 | Weeks 1-2 | Basic player movement |
| 2 | Weeks 3-4 | Core gameplay loop |
| 3 | Weeks 5-6 | Obstacles and enemies |
| 4 | Weeks 7-8 | Advanced features |
| 5 | Weeks 9-10 | Content and polish |

**Total Development Time**: 10 weeks

This plan provides a structured approach to developing ROBO-9 while maintaining flexibility for iteration and improvement based on testing and feedback.

## Additional Documentation
- [Collision System Developer Guide](collision-system.md) - How to work with and extend the collision system
