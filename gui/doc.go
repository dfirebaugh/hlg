// Package gui provides an immediate mode GUI API.
//
// In immediate mode, the caller owns all state and widget functions both
// handle input AND render in a single call. Widgets return true when their
// state changed or when they were activated.
//
// Example usage:
//
//	inputCtx := gui.NewDefaultInputContext()
//	ctx := gui.NewContext(inputCtx)
//
//	// In your render loop:
//	ctx.Begin()
//
//	if ctx.Button("Click me", 10, 10, 100, 40) {
//	    fmt.Println("Button clicked!")
//	}
//
//	if ctx.Slider("vol", &volume, 0, 1, 10, 60, 200, 20) {
//	    fmt.Println("Volume changed:", volume)
//	}
//
//	ctx.End()
package gui
