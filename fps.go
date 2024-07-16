package hlg

import (
	"time"
)

var (
	fpsEnabled = false
)

// EnableFPS enables the FPS counter in the window title.
func EnableFPS() {
	fpsEnabled = true
}

// DisableFPS disables the FPS counter in the window title.
func DisableFPS() {
	fpsEnabled = false
}

type fpsCounter struct {
	frameCount int
	lastTime   time.Time
	lastFPS    float64
}

func newFPSCounter() *fpsCounter {
	return &fpsCounter{
		lastTime: time.Now(),
	}
}

func (f *fpsCounter) Reset() {
	f.frameCount = 0
	f.lastTime = time.Now()
}

func (f *fpsCounter) Frame() {
	f.frameCount++
	elapsedTime := time.Since(f.lastTime)
	if elapsedTime >= time.Second {
		f.lastFPS = float64(f.frameCount)
		f.Reset()
	}
}

func (f *fpsCounter) GetFPS() float64 {
	return f.lastFPS
}
