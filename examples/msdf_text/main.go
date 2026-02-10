// Debug example for MSDF text rendering
package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/dfirebaugh/hlg"
)

// Catppuccin Mocha palette
var (
	catBase     = color.RGBA{30, 30, 46, 255}    // #1e1e2e
	catText     = color.RGBA{205, 214, 244, 255} // #cdd6f4
	catSubtext0 = color.RGBA{166, 173, 200, 255} // #a6adc8
	catOverlay2 = color.RGBA{147, 153, 178, 255} // #9399b2
	catFlamingo = color.RGBA{242, 205, 205, 255} // #f2cdcd
	catPink     = color.RGBA{245, 194, 231, 255} // #f5c2e7
	catMauve    = color.RGBA{203, 166, 247, 255} // #cba6f7
	catPeach    = color.RGBA{250, 179, 135, 255} // #fab387
	catYellow   = color.RGBA{249, 226, 175, 255} // #f9e2af
	catGreen    = color.RGBA{166, 227, 161, 255} // #a6e3a1
	catTeal     = color.RGBA{148, 226, 213, 255} // #94e2d5
	catSky      = color.RGBA{137, 220, 235, 255} // #89dceb
	catSapphire = color.RGBA{116, 199, 236, 255} // #74c7ec
	catBlue     = color.RGBA{137, 180, 250, 255} // #89b4fa
)

const (
	screenWidth  = 1024
	screenHeight = 700
)

var font *hlg.Font

func main() {
	hlg.SetWindowSize(screenWidth, screenHeight)
	hlg.SetTitle("MSDF Text Rendering Test")

	// Check for optional flags
	saveAtlas := false
	debugMode := -1 // -1 means not set (use default)
	useAtlas := ""  // Path prefix for pregenerated atlas (e.g., "assets/fonts/Noto/noto_atlas")
	for _, arg := range os.Args[1:] {
		if arg == "--save-atlas" {
			saveAtlas = true
		} else if strings.HasPrefix(arg, "--debug-mode=") {
			modeStr := strings.TrimPrefix(arg, "--debug-mode=")
			if m, err := strconv.Atoi(modeStr); err == nil && m >= 0 && m <= 3 {
				debugMode = m
			} else {
				fmt.Println("Invalid --debug-mode value. Use 0, 1, 2, or 3:")
				fmt.Println("  0: median(RGB) - MSDF reconstruction (default)")
				fmt.Println("  1: alpha only - true SDF fallback")
				fmt.Println("  2: visualize RGB - debug atlas colors")
				fmt.Println("  3: hard threshold - no AA, tests atlas data")
				os.Exit(1)
			}
		} else if arg == "--use-atlas" {
			// Default to the Noto atlas in assets/fonts/Noto/
			useAtlas = "assets/fonts/Noto/noto_atlas"
		} else if strings.HasPrefix(arg, "--use-atlas=") {
			useAtlas = strings.TrimPrefix(arg, "--use-atlas=")
		}
	}

	// Find font path (first non-flag argument)
	var fontPath string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "--") {
			fontPath = arg
			break
		}
	}

	var err error
	if useAtlas != "" {
		// Load from pregenerated atlas
		pngPath := useAtlas + ".png"
		jsonPath := useAtlas + ".json"
		font, err = hlg.LoadFontFromAtlas(pngPath, jsonPath)
		if err != nil {
			fmt.Printf("Failed to load font from atlas: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Font loaded from pregenerated atlas: %s\n", useAtlas)
	} else if fontPath != "" {
		// Load custom font from path
		font, err = hlg.LoadFont(fontPath)
		if err != nil {
			fmt.Printf("Failed to load font: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Custom font loaded successfully!")
	} else {
		// Use embedded default font (Noto Sans)
		font, err = hlg.LoadDefaultFont()
		if err != nil {
			fmt.Printf("Failed to load default font: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Default font (Noto Sans) loaded successfully!")
	}

	// Optionally save the atlas for debugging (pass --save-atlas flag)
	if saveAtlas {
		if err = font.SaveAtlasDebug("atlas_debug.png"); err != nil {
			fmt.Printf("Warning: could not save atlas debug image: %v\n", err)
		} else {
			fmt.Println("Atlas saved to atlas_debug.png")
		}

		if err = font.SaveAtlasJSON("atlas_debug.json"); err != nil {
			fmt.Printf("Warning: could not save atlas JSON: %v\n", err)
		} else {
			fmt.Println("Atlas JSON saved to atlas_debug.json")
		}
	}

	// Set the font atlas as active for primitive buffer rendering
	font.SetAsActiveAtlas()
	hlg.SetDefaultFont(font)

	// Apply debug mode if specified
	if debugMode >= 0 {
		hlg.SetMSDFMode(debugMode)
		modeNames := []string{"median(RGB)", "alpha only", "visualize RGB", "hard threshold"}
		fmt.Printf("MSDF debug mode set to %d: %s\n", debugMode, modeNames[debugMode])
	}

	// Optional: enable pixel snapping (diagnose subpixel fuzziness)
	for _, arg := range os.Args[1:] {
		if arg == "--snap" {
			hlg.EnableSnapMSDFToPixels(true)
			fmt.Println("MSDF pixel snapping enabled")
		}
	}

	hlg.Run(func() {}, func() {
		hlg.Clear(catBase)
		hlg.BeginDraw()

		x := 20
		y := 50
		hlg.Text("some small text", 0, 0, 10, catText)
		hlg.Text("s9egWQ", 0, 0, 256, catSapphire)
		hlg.Text("some small text", 0, 20, 12, catText)

		// Large title
		hlg.Text("MSDF Text Rendering", x, y, 16, catText)
		y += 60

		// Uppercase alphabet
		hlg.Text("ABCDEFGHIJKLMNOPQRSTUVWXYZ", x, y, 32, catPeach)
		y += 45

		// Lowercase alphabet
		hlg.Text("abcdefghijklmnopqrstuvwxyz", x, y, 32, catGreen)
		y += 45

		// Numbers
		hlg.Text("0123456789", x, y, 32, catSky)
		y += 45

		// Characters with counters (holes) - important MSDF test
		hlg.Text("Counters: eaobdpqg 0468 @#&%", x, y, 28, catYellow)
		y += 40

		// Descenders and ascenders
		hlg.Text("Descenders: gjpqy | Ascenders: bdfhklt", x, y, 28, catMauve)
		y += 40

		// Punctuation and symbols
		hlg.Text("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", x, y, 24, catTeal)
		y += 35

		// Mixed content - realistic text
		hlg.Text("The quick brown fox jumps over the lazy dog.", x, y, 24, catText)
		y += 35

		hlg.Text("Pack my box with five dozen liquor jugs!", x, y, 24, catFlamingo)
		y += 35

		// Different sizes comparison
		hlg.Text("Size 16px", x, y, 16, catOverlay2)
		hlg.Text("Size 24px", x+100, y, 24, catSubtext0)
		hlg.Text("Size 36px", x+250, y, 36, catText)
		y += 50

		// Edge cases - repeated characters
		hlg.Text("WWWWW MMMMM iiiii lllll", x, y, 28, catPink)
		y += 40

		// Numbers in context
		hlg.Text("Phone: (555) 123-4567 | Price: $99.99", x, y, 24, catGreen)
		y += 35

		// Email/URL style text
		hlg.Text("email@example.com | https://github.com", x, y, 22, catBlue)

		hlg.EndDraw()
	})

	font.Dispose()
}
