package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Circle)(nil)

// Circle describes a colored circle primitive in a Fyne canvas
type Circle struct {
	baseObject

	FillColor   color.Color
	StrokeColor color.Color
	StrokeWidth float32

	Shadow Shadow
}

// NewCircle returns a new Circle instance
func NewCircle(color color.Color) *Circle {
	return &Circle{FillColor: color}
}

// Hide will set this circle to not be visible
func (c *Circle) Hide() {
	c.baseObject.Hide()
	repaint(c)
}

// Move the circle object to a new position, relative to its parent / canvas
func (c *Circle) Move(pos fyne.Position) {
	if c.Position() == pos {
		return
	}
	c.baseObject.Move(pos)
	repaint(c)
}

// Refresh causes this object to be redrawn with its configured state.
func (c *Circle) Refresh() {
	Refresh(c)
}

// Resize on a circle updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
func (c *Circle) Resize(size fyne.Size) {
	if size == c.Size() {
		return
	}
	c.baseObject.Resize(size)
	Refresh(c)
}

// Show will set this circle to be visible
func (c *Circle) Show() {
	c.baseObject.Show()
	c.Refresh()
}
