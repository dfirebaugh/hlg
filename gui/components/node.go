package components

import (
	"github.com/dfirebaugh/hlg/gui"
)

type Component interface {
	Update()
	Render(ctx gui.DrawContext)
}

type Node struct {
	offsetX  int
	offsetY  int
	children []Component
}

func (n *Node) GetOffset() (int, int) {
	return n.offsetX, n.offsetY
}

func (n *Node) SetOffset(x int, y int) {
	n.offsetX = x
	n.offsetY = y
}

func (n *Node) AppendChild(child Component) {
	n.children = append(n.children, child)
}

func (n *Node) Update() {
	for _, c := range n.children {
		c.Update()
	}
}

func (n *Node) Render(ctx gui.DrawContext) {
	for _, c := range n.children {
		c.Render(ctx)
	}
}
