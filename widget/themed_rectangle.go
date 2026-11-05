// xssyne/widget/themed_rectangle.go

package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// ThemedRectangle — это умный прямоугольник. Он автоматически
// обновляет свой цвет, когда в приложении меняется тема.
// Вам больше не нужно вручную перекрашивать его в refreshThemeColors.
type ThemedRectangle struct {
	BaseWidget
	Rect       *canvas.Rectangle
	colorName  fyne.ThemeColorName
	strokeName fyne.ThemeColorName
}

// NewThemedRectangle создаёт прямоугольник, который сам следит за темой.
func NewThemedRectangle(colorName fyne.ThemeColorName) *ThemedRectangle {
	tr := &ThemedRectangle{
		Rect:      canvas.NewRectangle(color.Transparent),
		colorName: colorName,
	}
	tr.ExtendBaseWidget(tr)
	return tr
}

// SetStroke задаёт цвет обводки, который тоже будет меняться с темой.
func (tr *ThemedRectangle) SetStroke(strokeName fyne.ThemeColorName, width float32) {
	tr.strokeName = strokeName
	tr.Rect.StrokeWidth = width
	tr.Refresh()
}

// SetMinSize задаёт минимальный размер прямоугольника.
func (tr *ThemedRectangle) SetMinSize(size fyne.Size) {
	tr.Rect.SetMinSize(size)
}

// SetCornerRadius задаёт радиус скругления углов.
func (tr *ThemedRectangle) SetCornerRadius(radius float32) {
	tr.Rect.CornerRadius = radius
}

// SetFillColor задаёт статичный цвет (если не нужен цвет из темы).
func (tr *ThemedRectangle) SetFillColor(c color.Color) {
	tr.Rect.FillColor = c
	tr.Refresh()
}

// CreateRenderer — стандартный рендерер Fyne.
func (tr *ThemedRectangle) CreateRenderer() fyne.WidgetRenderer {
	return NewSimpleRenderer(tr.Rect)
}

// Refresh вызывается Fyne автоматически при смене темы.
func (tr *ThemedRectangle) Refresh() {
	th := tr.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	if tr.colorName != "" {
		tr.Rect.FillColor = th.Color(tr.colorName, v)
	}
	if tr.strokeName != "" {
		tr.Rect.StrokeColor = th.Color(tr.strokeName, v)
	}
	tr.Rect.Refresh()
	tr.BaseWidget.Refresh()
}
