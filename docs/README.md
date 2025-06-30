# ROBO-9 Platformer Documentation

## About This Project

ROBO-9 is a 2D platformer game built with Go and Ebitengine, featuring a robot character that collects energy hearts to help sad cats. The game emphasises precise movement, robust collision detection, and engaging platformer mechanics.

## About This Documentation

This directory contains comprehensive documentation for the ROBO-9 platformer game project. All documentation follows British English conventions and should be maintained to support both development and future contributors.

## Documentation Standards

### Language and Style
- **Language**: All documentation must be written in British English
- **Spelling**: Use British spellings (e.g., "colour" not "color", "realise" not "realize", "centre" not "center")
- **Terminology**: Prefer British terminology where applicable
- **Consistency**: Maintain consistent terminology throughout all documents

### Adding New Documentation
When adding new documentation:

1. **File Naming**: Use kebab-case for file names (e.g., `asset-loading-system.md`)
2. **Structure**: Follow the established structure with clear headings and sections
3. **Code Examples**: Include practical code examples where relevant
4. **Cross-References**: Link to related documentation and update this index
5. **Status Tracking**: Use checkboxes and status indicators for implementation progress

### Updating Documentation
- Keep documentation current with code changes
- Update implementation status as features are completed
- Review language for British English compliance
- Ensure all links and references remain valid

## Development

* [Development Plan](development-plan.md) - Complete roadmap and timeline for the ROBO-9 platformer

## Core Systems

### Asset Management
* [Asset Loading System](asset-loading-system.md) - Centralised asset management and loading

### Game State Management  
* [Game State Management](game-state-management.md) - State management, transitions, and callback system
* [Game State Management Quick Reference](game-state-management-quick-reference.md) - Common patterns and operations

### Collision System
* [Collision System Developer Guide](collision-system.md) - How to work with and extend the collision system

### Animation & Graphics
* [Animation System](animation-system.md) - Animation controller and state management
* [Drawing an Image](draw-image.md) - Ebitengine image rendering basics

## Entity Implementation

* [Player Implementation](player-implementation.md) - ROBO-9 player entity design and features

## Technical Specifications

* [ROBO-9 Sprite Specification](robo9-sprite-specification.md) - Sprite sheet layout and animation specifications

## Development Notes & Fixes

* [Control Scheme Fix](control-scheme-fix.md) - Input handling improvements
* [Screen Size Update](screen-size-update.md) - Screen resolution and scaling adjustments

---

## Quick Navigation

### For New Developers
1. Start with the [Development Plan](development-plan.md) for project overview
2. Review [Player Implementation](player-implementation.md) for core gameplay mechanics
3. Read [Collision System Developer Guide](collision-system.md) for physics integration

### For System Integration
- [Asset Loading System](asset-loading-system.md) - Loading and managing game assets
- [Game State Management](game-state-management.md) - Managing different game states
- [Animation System](animation-system.md) - Character and entity animations

### For Technical Reference
- [ROBO-9 Sprite Specification](robo9-sprite-specification.md) - Asset specifications
- [Game State Management Quick Reference](game-state-management-quick-reference.md) - Common operations
