package main

import (
	"github.com/dfirebaugh/ggez"
	"golang.org/x/image/colornames"
)

func main() {
	ggez.SetWindowSize(720, 480)
	ggez.SetScreenSize(240, 160)
	// construct a triangle (returns a Renderable that is added to the render queue)
	t := ggez.Triangle(0, 160, 120, 0, 240, 160, colornames.Orangered)

	// dispose of the triangle texture when we're done with it
	defer t.Dispose()

	// fill the screen with blue
	ggez.Clear(colornames.Skyblue)

	// draw the triangle -- we don't need to re-render it,
	// it will draw to the screen until we dispose it
	// so we don't need to call it in Update()
	// renderables are rendered in the order in which .Render() is called
	t.Render()

	// call Update() to run the application and draw the screen
	// calling update will keep the application from exiting
	ggez.Update(func() {
		// the callback function doesn't need to do anything in this case
		// typically, this is where you would do things like update the state of your game/application
		// in this case it just keeps the app running
	})
}
