// xssyne/widget/radio_item.go

package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var (
	_ fyne.Widget       = (*radioItem)(nil)
	_ desktop.Hoverable = (*radioItem)(nil)
	_ fyne.Tappable     = (*radioItem)(nil)
	_ fyne.Focusable    = (*radioItem)(nil)
)

func newRadioItem(label string, onTap func(*radioItem)) *radioItem {
	i := &radioItem{Label: label, onTap: onTap}
	i.ExtendBaseWidget(i)
	return i
}

// radioItem is a single radio item to be used by RadioGroup.
type radioItem struct {
	DisableableWidget

	Label    string
	Selected bool

	focused bool
	hovered bool
	onTap   func(item *radioItem)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (i *radioItem) CreateRenderer() fyne.WidgetRenderer {
	txt := canvas.Text{Alignment: fyne.TextAlignLeading}
	txt.TextSize = i.Theme().Size(theme.SizeNameText)
	r := &radioItemRenderer{item: i, label: &txt}
	r.SetObjects([]fyne.CanvasObject{&r.focusIndicator, &r.icon, &r.over, &txt})
	r.update()
	return r
}

// FocusGained is called when this item gained the focus.
func (i *radioItem) FocusGained() {
	i.focused = true
	i.Refresh()
}

// FocusLost is called when this item lost the focus.
func (i *radioItem) FocusLost() {
	i.focused = false
	i.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
func (i *radioItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Disabled() {
		return
	}

	i.hovered = true
	i.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (i *radioItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (i *radioItem) MouseOut() {
	if i.Disabled() {
		return
	}

	i.hovered = false
	i.Refresh()
}

// SetSelected sets whether this radio item is selected or not.
func (i *radioItem) SetSelected(selected bool) {
	if i.Disabled() || i.Selected == selected {
		return
	}

	i.Selected = selected
	i.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (i *radioItem) Tapped(_ *fyne.PointEvent) {
	if !i.focused {
		focusIfNotMobile(i.super())
	}
	i.toggle()
}

// TypedKey is called when this item receives a key event.
func (i *radioItem) TypedKey(_ *fyne.KeyEvent) {
}

// TypedRune is called when this item receives a char event.
func (i *radioItem) TypedRune(r rune) {
	if r == ' ' {
		i.toggle()
	}
}

func (i *radioItem) toggle() {
	if i.Disabled() || i.onTap == nil {
		return
	}

	i.onTap(i)
}

// xssyne
type radioItemRenderer struct {
	widget.BaseRenderer
	item *radioItem

	focusIndicator canvas.Circle
	icon, over     canvas.Image
	label          *canvas.Text

	// кеш ресурсов: (raw, colorName) — пересоздаём только при изменении
	lastInRaw, lastOutRaw     fyne.Resource
	lastInColor, lastOutColor fyne.ThemeColorName
}

func (r *radioItemRenderer) Layout(size fyne.Size) {
	th := r.item.Theme()
	innerPadding := th.Size(theme.SizeNameInnerPadding)
	borderSize := th.Size(theme.SizeNameInputBorder)
	iconInlineSize := th.Size(theme.SizeNameInlineIcon)

	focusIndicatorSize := fyne.NewSquareSize(iconInlineSize + innerPadding)
	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(fyne.NewPos(borderSize, (size.Height-focusIndicatorSize.Height)/2))

	labelSize := fyne.NewSize(size.Width, size.Height)
	r.label.Resize(labelSize)
	r.label.Move(fyne.NewPos(focusIndicatorSize.Width+th.Size(theme.SizeNamePadding), 0))

	iconPos := fyne.NewPos(innerPadding/2+borderSize, (size.Height-iconInlineSize)/2)
	iconSize := fyne.NewSquareSize(iconInlineSize)
	r.icon.Resize(iconSize)
	r.icon.Move(iconPos)
	r.over.Resize(iconSize)
	r.over.Move(iconPos)
}

func (r *radioItemRenderer) MinSize() fyne.Size {
	th := r.item.Theme()
	inPad := th.Size(theme.SizeNameInnerPadding) * 2

	return r.label.MinSize().
		AddWidthHeight(inPad+th.Size(theme.SizeNameInlineIcon)+th.Size(theme.SizeNamePadding), inPad)
}

func (r *radioItemRenderer) Refresh() {
	r.update()
	canvas.Refresh(r.item.super())
}

// xssyne
func (r *radioItemRenderer) update() {
	th := r.item.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.label.Text = r.item.Label
	r.label.TextSize = th.Size(theme.SizeNameText)
	if r.item.Disabled() {
		r.label.Color = th.Color(theme.ColorNameDisabled, v)
	} else {
		r.label.Color = th.Color(theme.ColorNameForeground, v)
	}

	outColor := theme.ColorNameInputBorder
	inColor := theme.ColorNameInputBackground
	if r.item.Selected {
		inColor = theme.ColorNamePrimary
		outColor = theme.ColorNameForeground
	}
	if r.item.Disabled() {
		if r.item.Selected {
			inColor = theme.ColorNameDisabled
		} else {
			inColor = theme.ColorNameBackground
		}
		outColor = theme.ColorNameDisabled
	}

	rawIn := th.Icon(theme.IconNameRadioButtonFill)
	if r.lastInRaw != rawIn || r.lastInColor != inColor {
		in := theme.NewThemedResource(rawIn)
		in.ColorName = inColor
		r.icon.Resource = in
		r.lastInRaw = rawIn
		r.lastInColor = inColor
		r.icon.Refresh()
	}
	rawOut := th.Icon(theme.IconNameRadioButton)
	if r.lastOutRaw != rawOut || r.lastOutColor != outColor {
		out := theme.NewThemedResource(rawOut)
		out.ColorName = outColor
		r.over.Resource = out
		r.lastOutRaw = rawOut
		r.lastOutColor = outColor
		r.over.Refresh()
	}

	if r.item.Disabled() {
		r.focusIndicator.FillColor = color.Transparent
	} else if r.item.focused {
		r.focusIndicator.FillColor = th.Color(theme.ColorNameFocus, v)
	} else if r.item.hovered {
		r.focusIndicator.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		r.focusIndicator.FillColor = color.Transparent
	}
}
