package main

import (
	"math/rand"

	"github.com/dfirebaugh/hlg"
	"github.com/dfirebaugh/hlg/pkg/fb"
	"github.com/dfirebaugh/hlg/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	gridWidth  = 400
	gridHeight = 300
	cellSize   = 1
)

var (
	grid     [gridWidth][gridHeight]bool
	nextGrid [gridWidth][gridHeight]bool
)

func initGrid() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			grid[x][y] = rand.Float64() < 0.2
		}
	}
}

func countNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			neighborX := (x + i + gridWidth) % gridWidth
			neighborY := (y + j + gridHeight) % gridHeight
			if grid[neighborX][neighborY] {
				count++
			}
		}
	}
	return count
}

func updateGrid() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			neighbors := countNeighbors(x, y)
			if grid[x][y] {
				if neighbors < 2 || neighbors > 3 {
					nextGrid[x][y] = false
				} else {
					nextGrid[x][y] = true
				}
			} else {
				if neighbors == 3 {
					nextGrid[x][y] = true
				} else {
					nextGrid[x][y] = false
				}
			}
		}
	}

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			grid[x][y] = nextGrid[x][y]
		}
	}
}

func renderGridToTexture(fb *fb.ImageFB) {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			cellColor := colornames.Black
			if grid[x][y] {
				cellColor = colornames.White
			}
			for i := 0; i < cellSize; i++ {
				for j := 0; j < cellSize; j++ {
					fb.SetPixel(int16(x*cellSize+i), int16(y*cellSize+j), cellColor)
				}
			}
		}
	}
}

func main() {
	hlg.SetWindowSize(gridWidth, gridHeight)
	hlg.EnableFPS()
	initGrid()
	framebuffer := fb.New(gridWidth*cellSize, gridHeight*cellSize)

	renderGridToTexture(framebuffer)
	texture, _ := hlg.CreateTextureFromImage(framebuffer.ToImage())

	hlg.Run(func() {
		if hlg.IsKeyPressed(input.KeyR) {
			initGrid()
			renderGridToTexture(framebuffer)
			texture.UpdateImage(framebuffer.ToImage())
		}
		updateGrid()
		renderGridToTexture(framebuffer)
		texture.UpdateImage(framebuffer.ToImage())
	}, func() {
		hlg.Clear(colornames.Black)
		texture.Render()
	})
}
