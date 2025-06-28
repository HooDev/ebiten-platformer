# Displaying an Image in Ebitengine

This guide walks you through the steps required to load and display an image on the screen using [Ebitengine](https://ebitengine.org/).

* An image file (e.g. `hero.png`) placed in the same directory as your Go source code. This guide assumes a PNG file, but other formats are possible with the correct decoders.

## Basic Example

Below is a minimal Ebitengine program that loads and displays an image.

### 1. Create a file called `main.go`

```go
package main

import (
    "image/png"
    "log"
    "os"

    "github.com/hajimehoshi/ebiten/v2"
)

var myImage *ebiten.Image

// loadImage opens and decodes a PNG image from disk.
func loadImage(path string) *ebiten.Image {
    f, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    img, err := png.Decode(f)
    if err != nil {
        log.Fatal(err)
    }

    return ebiten.NewImageFromImage(img)
}

// Game implements ebiten.Game interface.
type Game struct{}

func (g *Game) Update() error {
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(100, 100) // Draw at position (100, 100)
    screen.DrawImage(myImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 640, 480
}

func main() {
    myImage = loadImage("hero.png")

    ebiten.SetWindowSize(640, 480)
    ebiten.SetWindowTitle("Display Image Example")

    if err := ebiten.RunGame(&Game{}); err != nil {
        log.Fatal(err)
    }
}
```

### 2. Run the game

```bash
go run main.go
```

A window should appear showing your image at coordinates `(100, 100)`.

## Notes

* To use JPEG, BMP, or other image formats, import and use the appropriate Go image decoder from the `image/*` packages.
* For performance, avoid repeatedly loading images inside the `Update` or `Draw` methods.
* You can use `ebiten.NewImageFromImage()` to convert a decoded `image.Image` into an Ebitengine-compatible image.

For more information, visit the [official Ebitengine documentation](https://ebitengine.org/documents/).

