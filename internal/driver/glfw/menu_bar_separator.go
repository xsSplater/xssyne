// xssyne/internal/driver/glfw/menu_bar_separator.go

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*menuBarSeparator)(nil)

type menuBarSeparator struct {
	widget.Base
}

func (s *menuBarSeparator) CreateRenderer() fyne.WidgetRenderer {
	line := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	line.SetMinSize(fyne.NewSize(1, 1)) // 1 px шириной
	return &menuBarSeparatorRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{line}),
		line:         line,
	}
}

type menuBarSeparatorRenderer struct {
	widget.BaseRenderer
	line *canvas.Rectangle
}

func (r *menuBarSeparatorRenderer) Layout(size fyne.Size) {
	// Линия занимает всю выделенную высоту (меню), ширину — 1 px от MinSize.
	r.line.Resize(fyne.NewSize(1, size.Height))
	r.line.Move(fyne.NewPos(0, 0))
}

func (r *menuBarSeparatorRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (r *menuBarSeparatorRenderer) Refresh() {
	r.line.FillColor = theme.Color(theme.ColorNameSeparator)
	r.line.Refresh()
	canvas.Refresh(r.line)
}
