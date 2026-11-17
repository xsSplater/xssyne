// xssyne/canvas/blur.go

package canvas

import "fyne.io/fyne/v2"

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Blur)(nil)

// Blur creates a rectangular blur region on the output.
// All objects drawn under this will be blurred, any above will not be affected.
//
// Since: 2.8
type Blur struct {
	baseObject

	// Radius refers to how far from a point should be used to calculate the blur.
	// It must be greater than 0 but no more than 50.
	Radius float32

	// Radius used to round the corners of the blur region.
	CornerRadius float32
}

// Hide will set this blur to not be visible
func (b *Blur) Hide() {
	b.baseObject.Hide()

	repaint(b)
}

// Move the blur to a new position, relative to its parent / canvas
func (b *Blur) Move(pos fyne.Position) {
	if b.Position() == pos {
		return
	}
	b.baseObject.Move(pos)

	repaint(b)
}

// Refresh causes this blur to be redrawn with its configured state.
func (b *Blur) Refresh() {
	Refresh(b)
}

// Resize on a blur updates the new size of this object.
func (b *Blur) Resize(s fyne.Size) {
	if s == b.Size() {
		return
	}

	b.baseObject.Resize(s)
	repaint(b)
}

// xssyne
// NewBlur returns a new Blur instance.
//
// Radius is clamped to the supported range (0.1 .. 50) because the shader
// implementation requires a strictly positive radius and a bounded kernel
// size. Раньше невалидный radius давал неопределённое поведение (пустой
// blur или артефакты), без какого-либо сообщения.
func NewBlur(radius float32) *Blur {
	if radius < 0.1 {
		radius = 0.1
	}
	if radius > 50 {
		radius = 50
	}
	return &Blur{
		Radius: radius,
	}
}
