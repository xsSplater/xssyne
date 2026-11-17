// xssyne/canvas/gradient.go

package canvas

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

// LinearGradient defines a Gradient travelling straight at a given angle.
// The only supported values for the angle are `0.0` (vertical) and `90.0` (horizontal), currently.
type LinearGradient struct {
	baseObject

	StartColor color.Color // The beginning color of the gradient
	EndColor   color.Color // The end color of the gradient
	Angle      float64     // The angle of the gradient (0/180 for vertical; 90/270 for horizontal)
}

// Generate calculates an image of the gradient with the specified width and height.
func (g *LinearGradient) Generate(iw, ih int) image.Image {
	w, h := float64(iw), float64(ih)
	var generator func(x, y float64) float64
	switch g.Angle {
	case 90, -270: // horizontal flipped
		generator = func(x, _ float64) float64 {
			return (w - x) / w
		}
	case 270, -90: // horizontal
		generator = func(x, _ float64) float64 {
			return x / w
		}
	case 45, -315: // diagonal negative flipped
		generator = func(x, y float64) float64 {
			return math.Abs((w - x + y) / (w + h)) // ((w+h)-(x+h-y)) / (w+h)
		}
	case 225, -135: // diagonal negative
		generator = func(x, y float64) float64 {
			return math.Abs((x + h - y) / (w + h))
		}
	case 135, -225: // diagonal positive flipped
		generator = func(x, y float64) float64 {
			return math.Abs((w + h - (x + y)) / (w + h))
		}
	case 315, -45: // diagonal positive
		generator = func(x, y float64) float64 {
			return math.Abs((x + y) / (w + h))
		}
	case 180, -180: // vertical flipped
		generator = func(_, y float64) float64 {
			return (h - y) / h
		}
	default: // vertical
		generator = func(_, y float64) float64 {
			return y / h
		}
	}
	return computeGradient(generator, iw, ih, g.StartColor, g.EndColor)
}

// Hide will set this gradient to not be visible
func (g *LinearGradient) Hide() {
	g.baseObject.Hide()

	repaint(g)
}

// Move the gradient to a new position, relative to its parent / canvas
func (g *LinearGradient) Move(pos fyne.Position) {
	if g.Position() == pos {
		return
	}

	g.baseObject.Move(pos)

	repaint(g)
}

// Resize resizes the gradient to a new size.
func (g *LinearGradient) Resize(size fyne.Size) {
	if size == g.Size() {
		return
	}
	g.baseObject.Resize(size)

	// refresh needed to invalidate cached textures
	g.Refresh()
}

// Refresh causes this gradient to be redrawn with its configured state.
func (g *LinearGradient) Refresh() {
	Refresh(g)
}

// RadialGradient defines a Gradient travelling radially from a center point outward.
type RadialGradient struct {
	baseObject

	StartColor color.Color // The beginning color of the gradient
	EndColor   color.Color // The end color of the gradient
	// The offset of the center for generation of the gradient.
	// This is not a DP measure but relates to the width/height.
	// A value of 0.5 would move the center by the half width/height.
	CenterOffsetX, CenterOffsetY float64
}

// Generate calculates an image of the gradient with the specified width and height.
func (g *RadialGradient) Generate(iw, ih int) image.Image {
	w, h := float64(iw), float64(ih)
	// define center plus offset
	centerX := w/2 + w*g.CenterOffsetX
	centerY := h/2 + h*g.CenterOffsetY

	// handle negative offsets
	var a, b float64
	if g.CenterOffsetX < 0 {
		a = w - centerX
	} else {
		a = centerX
	}
	if g.CenterOffsetY < 0 {
		b = h - centerY
	} else {
		b = centerY
	}

	generator := func(x, y float64) float64 {
		// calculate distance from center for gradient multiplier
		dx, dy := centerX-x, centerY-y
		da := math.Sqrt(dx*dx + dy*dy*a*a/b/b)
		if da > a {
			return 1
		}
		return da / a
	}
	return computeGradient(generator, iw, ih, g.StartColor, g.EndColor)
}

// Hide will set this gradient to not be visible
func (g *RadialGradient) Hide() {
	g.baseObject.Hide()

	repaint(g)
}

// Move the gradient to a new position, relative to its parent / canvas
func (g *RadialGradient) Move(pos fyne.Position) {
	if g.Position() == pos {
		return
	}

	g.baseObject.Move(pos)

	repaint(g)
}

// Resize resizes the gradient to a new size.
func (g *RadialGradient) Resize(size fyne.Size) {
	if size == g.Size() {
		return
	}
	g.baseObject.Resize(size)

	// refresh needed to invalidate cached textures
	g.Refresh()
}

// Refresh causes this gradient to be redrawn with its configured state.
func (g *RadialGradient) Refresh() {
	Refresh(g)
}

func computeGradient(generator func(x, y float64) float64, w, h int, startColor, endColor color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	if startColor == nil && endColor == nil {
		return img
	} else if startColor == nil {
		startColor = color.Transparent
	} else if endColor == nil {
		endColor = color.Transparent
	}

	// Premultiplied 16-bit начальный/конечный цвета (как возвращает Color.RGBA()).
	// Интерполяция ведётся в premultiplied-пространстве — это сохраняет
	// композитно-корректное поведение (красный→синий через 50% прозрачности
	// даёт полупрозрачный красновато-синий, а не грязно-фиолетовый).
	aR, aG, aB, aA := startColor.RGBA()
	bR, bG, bB, bA := endColor.RGBA()

	afR := float64(aR)
	afG := float64(aG)
	afB := float64(aB)
	afA := float64(aA)
	dR := float64(bR) - afR
	dG := float64(bG) - afG
	dB := float64(bB) - afB
	dA := float64(bA) - afA

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d := generator(float64(x)+0.5, float64(y)+0.5)
			if d < 0 {
				d = 0
			} else if d > 1 {
				d = 1
			}

			// Premultiplied результат (16-бит).
			pr := afR + d*dR
			pg := afG + d*dG
			pb := afB + d*dB
			pa := afA + d*dA

			// Превращаем в NRGBA 8-bit — прямая запись в Pix, без
			// промежуточного image.Set + interface boxing.
			i := img.PixOffset(x, y)
			if pa <= 0 {
				img.Pix[i+0] = 0
				img.Pix[i+1] = 0
				img.Pix[i+2] = 0
				img.Pix[i+3] = 0
				continue
			}
			// Формула как в color.nrgbaModel: unpremultiply + >>8.
			img.Pix[i+3] = uint8(uint32(pa+0.5) >> 8)
			if pa >= 65535 {
				img.Pix[i+0] = uint8(uint32(pr+0.5) >> 8)
				img.Pix[i+1] = uint8(uint32(pg+0.5) >> 8)
				img.Pix[i+2] = uint8(uint32(pb+0.5) >> 8)
			} else {
				pa16 := uint32(pa + 0.5)
				img.Pix[i+0] = uint8((uint32(pr+0.5) * 0xffff / pa16) >> 8)
				img.Pix[i+1] = uint8((uint32(pg+0.5) * 0xffff / pa16) >> 8)
				img.Pix[i+2] = uint8((uint32(pb+0.5) * 0xffff / pa16) >> 8)
			}
		}
	}
	return img
}

// NewHorizontalGradient creates a new horizontally travelling linear gradient.
// The start color will be at the left of the gradient and the end color will be at the right.
func NewHorizontalGradient(start, end color.Color) *LinearGradient {
	g := &LinearGradient{StartColor: start, EndColor: end}
	g.Angle = 270
	return g
}

// NewLinearGradient creates a linear gradient at the specified angle.
// The angle parameter is the degree angle along which the gradient is calculated.
// A NewHorizontalGradient uses 270 degrees and NewVerticalGradient is 0 degrees.
func NewLinearGradient(start, end color.Color, angle float64) *LinearGradient {
	g := &LinearGradient{StartColor: start, EndColor: end}
	g.Angle = angle
	return g
}

// NewRadialGradient creates a new radial gradient.
func NewRadialGradient(start, end color.Color) *RadialGradient {
	return &RadialGradient{StartColor: start, EndColor: end}
}

// NewVerticalGradient creates a new vertically travelling linear gradient.
// The start color will be at the top of the gradient and the end color will be at the bottom.
func NewVerticalGradient(start color.Color, end color.Color) *LinearGradient {
	return &LinearGradient{StartColor: start, EndColor: end}
}
