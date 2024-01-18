package ggez

import (
	"time"
)

var (
	fpsEnabled = false
)

func EnableFPS() {
	fpsEnabled = true
}

func DisableFPS() {
	fpsEnabled = false
}

type FPSCounter struct {
	frameCount int
	lastTime   time.Time
	lastFPS    float64
}

func NewFPSCounter() *FPSCounter {
	return &FPSCounter{
		lastTime: time.Now(),
	}
}

func (f *FPSCounter) Reset() {
	f.frameCount = 0
	f.lastTime = time.Now()
}

func (f *FPSCounter) Frame() {
	f.frameCount++
	elapsedTime := time.Since(f.lastTime)
	if elapsedTime >= time.Second {
		f.lastFPS = float64(f.frameCount)
		f.Reset()
	}
}

func (f *FPSCounter) GetFPS() float64 {
	return f.lastFPS
}
