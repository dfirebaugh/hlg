package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/hlg"
)

// Test that MSDF text respects draw order when using DrawText through DrawContext

var (
	catppuccinBase  = color.RGBA{30, 30, 46, 255}
	catppuccinBlue  = color.RGBA{137, 180, 250, 255}
	catppuccinGreen = color.RGBA{166, 227, 161, 255}
	catppuccinPeach = color.RGBA{250, 179, 135, 255}
	catppuccinText  = color.RGBA{205, 214, 244, 255}
)

func main() {
	hlg.SetWindowSize(800, 600)
	hlg.SetScreenSize(800, 600)
	hlg.SetTitle("MSDF Text Layering Test")

	font, err := hlg.LoadDefaultFont()
	if err != nil {
		panic(err)
	}

	// Set the font atlas as active for primitive buffer rendering
	font.SetAsActiveAtlas()
	hlg.SetDefaultFont(font)

	hlg.Run(func() {}, func() {
		hlg.Clear(catppuccinBase)
		hlg.BeginDraw()

		// Test 1: Text BEFORE rectangle - rectangle should cover text
		hlg.Text("This text is BEHIND the blue rectangle", 50, 80, 20, catppuccinText)

		hlg.FilledRect(50, 70, 400, 50, catppuccinBlue)

		// Test 2: Rectangle first, then text - text should be on top
		hlg.FilledRect(50, 160, 450, 80, catppuccinGreen)

		hlg.Text("This text is ON TOP of the green rectangle", 70, 190, 20, catppuccinPeach)

		// Test 3: Interleaved drawing
		hlg.FilledRect(50, 280, 200, 40, catppuccinBlue)

		hlg.Text("A", 100, 285, 28, catppuccinText)

		hlg.FilledRect(260, 280, 200, 40, catppuccinGreen)

		hlg.Text("B", 310, 285, 28, catppuccinText)

		hlg.FilledRect(470, 280, 200, 40, catppuccinPeach)

		hlg.Text("C", 520, 285, 28, catppuccinText)

		// Explanation text
		hlg.Text("Text now respects draw order!", 50, 380, 24, catppuccinText)
		hlg.Text("First line: text drawn BEFORE rectangle (should be hidden)", 50, 420, 16, catppuccinText)
		hlg.Text("Second line: text drawn AFTER rectangle (should be visible)", 50, 445, 16, catppuccinText)
		hlg.Text("Third line: interleaved A, B, C on colored backgrounds", 50, 470, 16, catppuccinText)

		hlg.EndDraw()
	})

	font.Dispose()
	fmt.Println("Done")
}
