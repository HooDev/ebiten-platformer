# Asset Loading System

## Overview

The ROBO-9 platformer uses a centralised asset loading system built on top of Ebitengine's asset management capabilities. This system provides efficient loading, caching, and management of game assets including images and audio files.

## Quick Start

For developers who want to get started immediately:

1. **Load an image**: `sprite, err := assetManager.LoadImage("player.png")`
2. **Batch load assets**: `err := assetManager.PreloadAssets([]string{"heart.png", "cat.png"}, nil)`
3. **Use in drawing**: `screen.DrawImage(sprite, &ebiten.DrawImageOptions{})`

## Architecture

### Core Components

#### AssetManager
The `AssetManager` struct handles all asset operations:

```go
type AssetManager struct {
    images      map[string]*ebiten.Image  // Cache for loaded images
    mu          sync.RWMutex              // Thread-safe access to cache
    assetDir    string                    // Base directory for assets
    useEmbedded bool                      // Whether to use embedded assets
    embeddedFS  embed.FS                  // Embedded filesystem (for production builds)
}
```

#### AssetConfig
Configuration structure for initialising the asset manager:

```go
type AssetConfig struct {
    AssetDir    string      // Path to assets directory (e.g., "assets")
    UseEmbedded bool        // Enable embedded assets for production
    EmbeddedFS  embed.FS    // Embedded filesystem instance
}
```

### Basic Setup

```go
// Create asset manager configuration
config := engine.AssetConfig{
    AssetDir:    "assets",
    UseEmbedded: false,  // Use filesystem assets during development
}

// Initialise asset manager
assetManager := engine.NewAssetManager(config)
```

## Features

### ✅ Implemented (Phase 1)

- **Lazy Loading**: Images loaded on first request and cached for subsequent use
- **Thread-Safe**: Concurrent access handled with read-write mutexes
- **Error Handling**: Comprehensive error reporting for failed loads
- **Format Support**: PNG images (primary format for sprites and textures)
- **Filesystem Flexibility**: Support for both regular filesystem and embedded assets
- **Automatic Caching**: All loaded assets automatically cached in memory
- **Cache Management**: Methods to query cache status and clear cache when needed
- **Preloading**: Batch loading of assets for performance optimisation
- **Asset Counting**: Track number of loaded assets for debugging
- **Asset Listing**: Enumerate all cached assets

### 🔄 Planned (Phase 4)

- **Audio Loading**: WAV file support with caching and streaming
- **Audio Player Management**: Background music and sound effects
- **Embedded Assets**: Production builds with embedded assets

## Developer Guide

### Setting Up Your Game

When creating your game instance, the asset manager is automatically initialised:

```go
func main() {
    game := NewRoboGame()
    
    // Asset manager is ready to use through game.GetAssetManager()
    if err := game.LoadAssets(); err != nil {
        log.Fatalf("Failed to load assets: %v", err)
    }
    
    ebiten.RunGame(game)
}
```

### Loading Assets

#### Individual Assets
```go
func (g *RoboGame) LoadAssets() error {
    assetManager := g.GetAssetManager()
    
    // Load individual assets
    playerSprite, err := assetManager.LoadImage("player.png")
    if err != nil {
        return fmt.Errorf("failed to load player: %w", err)
    }
    g.playerImage = playerSprite
    
    return nil
}
```

#### Batch Loading
```go
func (g *RoboGame) PreloadAssets() error {
    assetManager := g.GetAssetManager()
    
    imagePaths := []string{
        "heart.png",
        "cat_sad.png", 
        "cat_happy.png",
        "platform.png",
    }
    
    return assetManager.PreloadAssets(imagePaths, nil)
}
```

### Using Assets in Game Logic

```go
func (g *RoboGame) Draw(screen *ebiten.Image) {
    assetManager := g.GetAssetManager()
    
    // Get a cached asset (no file I/O)
    heartSprite, _ := assetManager.GetImage("heart.png")
    
    // Draw the sprite
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(100, 100)
    screen.DrawImage(heartSprite, op)
}
```

### Asset Organisation

Structure your assets directory logically:

```
assets/
├── sprites/
│   ├── characters/
│   │   ├── player.png
│   │   ├── cat_sad.png
│   │   └── cat_happy.png
│   ├── items/
│   │   ├── heart.png
│   │   └── coin.png
│   └── enemies/
│       └── drone.png
├── tiles/
│   ├── platform.png
│   ├── spike.png
│   └── background.png
└── ui/
    ├── button.png
    └── font.png
```

Then load with relative paths:

```go
playerImg, err := assetManager.LoadImage("sprites/characters/player.png")
heartImg, err := assetManager.LoadImage("sprites/items/heart.png")
platformImg, err := assetManager.LoadImage("tiles/platform.png")
```

### Error Handling

Always handle asset loading errors gracefully:

```go
func (g *Game) LoadCriticalAssets() error {
    assetManager := g.GetAssetManager()
    
    criticalAssets := []string{
        "sprites/characters/player.png",
        "tiles/platform.png", 
    }
    
    for _, assetPath := range criticalAssets {
        if _, err := assetManager.LoadImage(assetPath); err != nil {
            return fmt.Errorf("critical asset %s failed to load: %w", assetPath, err)
        }
    }
    
    return nil
}
```

### Memory Management

Monitor and manage memory usage during development:

```go
func (g *Game) Update() error {
    // Check memory usage periodically
    if g.frameCount%3600 == 0 { // Every 60 seconds at 60 FPS
        assetManager := g.GetAssetManager()
        log.Printf("Loaded assets: %d images", assetManager.GetLoadedImageCount())
    }
    
    // Clear cache when changing levels
    if g.levelChanged {
        assetManager := g.GetAssetManager()
        assetManager.ClearCache()
        g.loadLevelAssets()
    }
    
    return nil
}
```

### Cache Management

```go
// Check cache status
imageCount := assetManager.GetLoadedImageCount()
fmt.Printf("Loaded %d images\n", imageCount)

// List cached assets
images, audio := assetManager.ListCachedAssets()
fmt.Printf("Cached images: %v\n", images)

// Clear cache (useful for memory management)
assetManager.ClearCache()
```

## Common Development Patterns

### Lazy Loading Pattern
```go
type Player struct {
    sprite *ebiten.Image
    assetManager *engine.AssetManager
}

func (p *Player) GetSprite() *ebiten.Image {
    if p.sprite == nil {
        var err error
        p.sprite, err = p.assetManager.LoadImage("player.png")
        if err != nil {
            log.Printf("Failed to load player sprite: %v", err)
            return nil
        }
    }
    return p.sprite
}
```

### Level-Based Asset Loading
```go
func (g *Game) LoadLevelAssets(levelName string) error {
    assetManager := g.GetAssetManager()
    
    // Define assets per level
    levelAssets := map[string][]string{
        "level1": {"grass.png", "tree.png", "cloud.png"},
        "level2": {"metal.png", "pipe.png", "steam.png"},
        "boss": {"boss_arena.png", "boss_sprite.png"},
    }
    
    assets, exists := levelAssets[levelName]
    if !exists {
        return fmt.Errorf("unknown level: %s", levelName)
    }
    
    return assetManager.PreloadAssets(assets, nil)
}
```

### Asset Validation
```go
func (g *Game) ValidateAssets() error {
    assetManager := g.GetAssetManager()
    requiredAssets := []string{
        "player.png",
        "heart.png", 
        "cat_sad.png",
        "cat_happy.png",
    }
    
    var missingAssets []string
    for _, asset := range requiredAssets {
        if _, err := assetManager.LoadImage(asset); err != nil {
            missingAssets = append(missingAssets, asset)
        }
    }
    
    if len(missingAssets) > 0 {
        return fmt.Errorf("missing required assets: %v", missingAssets)
    }
    
    return nil
}
```

## Performance Considerations

### Memory Management
- **Lazy Loading**: Assets are only loaded when needed
- **Caching**: Prevents redundant file I/O operations
- **Manual Cache Control**: Developers can clear cache when memory is constrained

### Thread Safety
- **Read-Write Mutexes**: Allow concurrent read access while protecting writes
- **Atomic Operations**: Cache checks and updates are thread-safe

### Loading Optimisation
- **Batch Preloading**: Load multiple assets in a single operation
- **Error Aggregation**: Collect all loading errors before failing
- **Logging**: Track loading operations for debugging

## Development vs Production

### Development Mode
```go
config := AssetConfig{
    AssetDir:    "assets",
    UseEmbedded: false,  // Load from filesystem
}
```

- Assets loaded from filesystem
- Easy to modify and test assets
- Faster iteration during development

### Production Mode (Planned)
```go
//go:embed assets/*
var embeddedAssets embed.FS

config := AssetConfig{
    UseEmbedded: true,
    EmbeddedFS:  embeddedAssets,
}
```

- Assets embedded in executable
- Single file distribution
- No external asset dependencies

## Error Handling

The asset loading system provides detailed error messages:

```go
// Example error messages
"failed to open image assets/player.png: no such file or directory"
"failed to decode image assets/corrupted.png: png: invalid format"
"failed to preload assets:\nimage player.png: file not found\nimage heart.png: decode error"
```

## Testing

### Manual Testing
```bash
# Development (WSL)
GOOS=windows go run main.go

# Or use the convenience script
./run.wsl.sh
```

### Asset Verification
The game displays asset loading status:
- Number of loaded images
- Loading progress
- Error messages for failed loads

## Troubleshooting

### Common Issues

**File Not Found Errors**
- Verify asset paths are relative to the configured `AssetDir`
- Ensure files exist in the assets directory
- Check file permissions

**Memory Issues**
- Use `ClearCache()` to free memory when switching levels
- Monitor loaded asset count with `GetLoadedImageCount()`
- Consider lazy loading for large asset sets

**WSL Environment**
- Use `GOOS=windows go run main.go` for WSL development
- Ensure proper file path separators
- Use the provided `run.wsl.sh` script

## Implementation Status

| Feature | Status | Phase |
|---------|--------|-------|
| Image Loading | ✅ Complete | 1 |
| Image Caching | ✅ Complete | 1 |
| Preloading | ✅ Complete | 1 |
| Error Handling | ✅ Complete | 1 |
| Thread Safety | ✅ Complete | 1 |
| Cache Management | ✅ Complete | 1 |
| Audio Loading | 🔄 Planned | 4 |
| Embedded Assets | 🔄 Planned | 4 |
| Hot Reloading | 🔄 Planned | 5 |

## Future Enhancements

### Phase 4 - Audio System
- [ ] WAV file loading and caching
- [ ] Audio player management
- [ ] Background music support
- [ ] Sound effect triggering
- [ ] Audio format validation

### Phase 5 - Advanced Features
- [ ] Compressed texture support
- [ ] Sprite sheet management
- [ ] Animation frame loading
- [ ] Asset hot-reloading for development
- [ ] Asset dependency tracking

The asset loading system provides a solid foundation for the ROBO-9 platformer, with room for expansion as the game develops through subsequent phases.
