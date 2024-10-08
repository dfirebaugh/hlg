
# Sprites

The `Sprite` struct in the `hlg` package provides a way to handle sprite sheets for rendering animations or a series of images in a grid-like structure. Each `Sprite` is a subset of a larger texture, representing a single frame of an animation or a single image within a sprite sheet.

## Constructor Function

### NewSprite

Creates a new `Sprite` instance from a given image, frame size, and sheet size.

```golang
func NewSprite(img image.Image, frameSize, sheetSize image.Point) *Sprite
```

- `img`: The source image, typically a sprite sheet.
- `frameSize`: The size of a single frame within the sprite sheet.
- `sheetSize`: The size of the sprite sheet in terms of the number of frames horizontally and vertically.

## Methods

### NextFrame

Advances the sprite to the next frame in the sprite sheet. The method updates the clipping region of the underlying texture to display the next frame.

```golang
func (s *Sprite) NextFrame()
```

- This method increments the current frame and wraps around if it reaches the end of the sprite sheet.
- Automatically updates the clip region to render the new current frame.
- Calls `Render()` to render the sprite with the updated frame.

## Usage Example

```golang
// Assuming you have a sprite sheet image and you know the frame and sheet sizes
spriteSheetImage := // Load your sprite sheet image
frameSize := image.Pt(32, 32) // Each frame is 32x32 pixels
sheetSize := image.Pt(4, 4)  // 4 columns and 4 rows in the sprite sheet

sprite := hlg.NewSprite(spriteSheetImage, frameSize, sheetSize)

hlg.Run(func() {
    sprite.NextFrame() // Update the sprite to the next frame each update call
}, func() {
    hlg.Clear(colornames.Black)
})
```

In this usage example, `sprite.NextFrame()` is called in the `hlg.Update` loop to animate the sprite by cycling through frames on each update.

## Full Example

```golang
package main

import (
	"bytes"
	"image"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/assets"
)

func main() {
	hlg.SetWindowSize(600, 600)
	hlg.SetTitle("hlg sprite example")

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}

	sprite := hlg.NewSprite(img, frameSize, sheetSize)
	sprite.Scale(.5, .5)

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 200

	hlg.Run(func() {
			sprite.NextFrame()
    }, func() {
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
		}
	})
}
```
