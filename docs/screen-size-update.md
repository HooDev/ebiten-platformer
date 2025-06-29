# Screen Size Update

## Changes Made

The game screen size has been increased to provide more room for UI elements and debug information.

### Previous Resolution
- **Game Resolution**: 320×240 pixels
- **Window Size**: 640×480 pixels (2x scaling)
- **Total Area**: 76,800 pixels

### New Resolution
- **Game Resolution**: 480×360 pixels  
- **Window Size**: 960×720 pixels (2x scaling)
- **Total Area**: 172,800 pixels

### Benefits

1. **More UI Space**: 
   - Debug information is no longer cramped
   - Menu text has better spacing
   - Room for future UI elements

2. **Better Gameplay Area**:
   - Ground level moved from Y=200 to Y=300
   - More vertical space for jumping and platforming
   - Room for taller level elements

3. **Improved Readability**:
   - Text is less crowded
   - Better visual hierarchy in menus
   - Clearer debug information layout

### Technical Details

#### Code Changes
- Updated `GameConfig` screen dimensions
- Adjusted overlay image size
- Moved ground collision from Y=200 to Y=300
- Repositioned all UI elements for better spacing

#### Aspect Ratio
- Maintains 4:3 aspect ratio (480:360 = 4:3)
- Consistent with retro gaming aesthetics
- Good balance between modern usability and classic feel

#### Performance Impact
- Minimal performance impact due to efficient rendering
- Still well within acceptable limits for 2D games
- Maintains smooth 60 FPS operation

### Updated Coordinates

| Element | Old Position | New Position |
|---------|-------------|-------------|
| Ground Level | Y=200 | Y=300 |
| Title (Menu) | (80,50) | (160,80) |
| Game Over Text | (120,100) | (190,150) |
| Settings Text | (120,50) | (200,80) |
| Debug Info | (10,70) | (10,90) |
| Asset Stats | (200,150) | (320,200) |

The larger screen provides a much better development and gameplay experience while maintaining the retro pixel-art aesthetic.
