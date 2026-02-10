//go:build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// CompareColors checks if two colors are within tolerance
// Uses perceptual difference in sRGB space
func CompareColors(expected, actual color.Color, tolerance float64) bool {
	er, eg, eb, ea := expected.RGBA()
	ar, ag, ab, aa := actual.RGBA()

	// Convert to 0-1 range
	diff := math.Abs(float64(er)-float64(ar))/65535.0 +
		math.Abs(float64(eg)-float64(ag))/65535.0 +
		math.Abs(float64(eb)-float64(ab))/65535.0 +
		math.Abs(float64(ea)-float64(aa))/65535.0

	return diff/4.0 < tolerance
}

// AnalyzeScreenshot reads a screenshot and reports color values at specific positions
func AnalyzeScreenshot(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	// Sample positions (center of each swatch in the grid)
	// Based on: startX=20, startY=20, swatchW=50, swatchH=50, padding=10, cols=6
	swatchW, swatchH := 50, 50
	padding := 10
	startX, startY := 20, 20

	type sample struct {
		name     string
		expected color.RGBA
		col, row int
	}

	samples := []sample{
		{"Red", color.RGBA{255, 0, 0, 255}, 0, 0},
		{"Green", color.RGBA{0, 255, 0, 255}, 1, 0},
		{"Blue", color.RGBA{0, 0, 255, 255}, 2, 0},
		{"Gray 50%", color.RGBA{128, 128, 128, 255}, 2, 1},
		{"Orange", color.RGBA{255, 128, 0, 255}, 5, 1},
	}

	fmt.Println("Color Analysis Results:")
	fmt.Println("========================")

	for _, s := range samples {
		// Calculate center of swatch
		x := startX + s.col*(swatchW+padding) + swatchW/2
		y := startY + s.row*(swatchH+padding) + swatchH/2

		actual := img.At(x, y)
		ar, ag, ab, aa := actual.RGBA()

		// Convert from 16-bit to 8-bit
		actualR := uint8(ar >> 8)
		actualG := uint8(ag >> 8)
		actualB := uint8(ab >> 8)
		actualA := uint8(aa >> 8)

		match := "PASS"
		if !CompareColors(s.expected, actual, 0.05) {
			match = "FAIL"
		}

		fmt.Printf("%s: expected RGB(%d,%d,%d) got RGB(%d,%d,%d,%d) [%s]\n",
			s.name,
			s.expected.R, s.expected.G, s.expected.B,
			actualR, actualG, actualB, actualA,
			match)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run screenshot_test.go <screenshot.png>")
		fmt.Println("")
		fmt.Println("Take a screenshot of the color_test window and pass it to this tool.")
		fmt.Println("On macOS: Cmd+Shift+4, then space, then click the window")
		os.Exit(1)
	}

	if err := AnalyzeScreenshot(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
