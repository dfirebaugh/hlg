package hlg

import (
	"time"
)

var fpsEnabled = false

// EnableFPS enables the FPS counter in the window title.
func EnableFPS() {
	fpsEnabled = true
}

// DisableFPS disables the FPS counter in the window title.
func DisableFPS() {
	fpsEnabled = false
}

type fpsCounter struct {
	frameCount   int
	lastTime     time.Time
	accumTime    time.Duration
	lastFPS      float64
	updatePeriod time.Duration
}

func newFPSCounter() *fpsCounter {
	return &fpsCounter{
		lastTime:     time.Now(),
		updatePeriod: time.Second, // Update FPS every second
	}
}

func (f *fpsCounter) Frame() {
	f.frameCount++
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(f.lastTime)

	f.accumTime += elapsedTime
	f.lastTime = currentTime

	if f.accumTime >= f.updatePeriod {
		// Calculate FPS
		f.lastFPS = float64(f.frameCount) / f.accumTime.Seconds()

		// Reset the counters
		f.accumTime = 0
		f.frameCount = 0
	}
}

func (f *fpsCounter) GetFPS() float64 {
	return f.lastFPS
}
