# Control Scheme Fix - Space Key Conflict Resolution

## Issue
The Space key was being used for both jumping (player movement) and pausing (menu navigation), creating a conflict where players couldn't jump without accidentally pausing the game.

## Solution
Removed Space key from menu/pause controls and kept it exclusively for jumping.

## Updated Control Scheme

### Game Controls (In-Game)
- **Movement**: WASD / Arrow Keys
- **Jump**: Space / W / Up Arrow
- **Climb Up**: Up Arrow / W (when near climbable surface)
- **Climb Down**: Down Arrow / S (when climbing)

### System Controls
- **Pause**: Escape (only)
- **Menu**: M
- **Resume**: Escape (from pause screen)

### Menu Navigation
- **Start Game**: Enter
- **Settings**: S
- **Back**: Escape / Backspace

## Changes Made

### Code Changes (`main.go`)
1. Removed `ebiten.KeySpace` from pause input handling in `StatePlaying`
2. Removed `ebiten.KeySpace` from resume input handling in `StatePaused`
3. Updated UI text to reflect correct controls

### Documentation Updates
1. Updated `README.md` control table
2. Updated in-game help text
3. Updated pause screen instructions

## Benefits
- **Clear Control Separation**: Jump and pause are now distinct actions
- **Better UX**: Players can jump freely without accidentally pausing
- **Intuitive Design**: Space for jumping is standard in platformers
- **Consistent Navigation**: Escape for system controls (pause/back) is standard

## Testing
After the fix:
- Space key only triggers jump action
- Escape key handles all pause/resume functionality
- No control conflicts during gameplay
- Menu navigation remains intuitive

The control scheme now follows standard platformer conventions where Space is dedicated to jumping and Escape handles system functions.
