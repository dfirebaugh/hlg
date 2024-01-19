
# Renderables


`hlg` allows you to create create objects (e.g. textures, sprites, shapes) and render them to the screen.

These objects implement the `Renderable` interface.

```golang
type Renderable interface {
	Render() // add the Renderable to the render queue
	Dispose() // remove the Renderable from the render queue and destroy the Renderable
	Hide() // keep the Renderable in memory but don't render it to the screen
}
```

Objects are rendered in the order in which `Render()` was called.

