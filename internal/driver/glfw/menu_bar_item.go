// xssyne/internal/driver/glfw/menu_bar_item.go

package glfw

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
	publicWidget "fyne.io/fyne/v2/widget"
)

var (
	_ fyne.Widget       = (*menuBarItem)(nil)
	_ desktop.Hoverable = (*menuBarItem)(nil)
	_ fyne.Focusable    = (*menuBarItem)(nil)
)

// menuBarItem is a widget for displaying an item for a fyne.Menu in a MenuBar.
type menuBarItem struct {
	widget.Base
	Menu   *fyne.Menu
	Parent *MenuBar

	active  bool
	child   *publicWidget.Menu
	hovered bool
}

func (i *menuBarItem) Child() *publicWidget.Menu {
	if i.child == nil {
		child := publicWidget.NewMenu(i.Menu)
		child.Hide()
		child.OnDismiss = i.Parent.deactivate
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu bar item.
func (i *menuBarItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.Color(theme.ColorNameHover))
	background.CornerRadius = theme.SelectionRadiusSize()
	background.Hide()

	indicator := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBarAccent))
	indicator.Hide()

	text := canvas.NewText(i.Menu.Label, theme.Color(theme.ColorNameForeground))
	shadow := canvas.NewText(i.Menu.Label, color.NRGBA{R: 0, G: 0, B: 0, A: 0x90})
	shadow.TextStyle = text.TextStyle
	objects := []fyne.CanvasObject{background, indicator, shadow, text}

	r := &menuBarItemRenderer{
		BaseRenderer: widget.NewBaseRenderer(objects),
		i:            i,
		text:         text,
		shadow:       shadow,
		background:   background,
		indicator:    indicator,
	}
	if i.Menu.Icon != nil {
		icon := canvas.NewImageFromResource(i.Menu.Icon)
		icon.FillMode = canvas.ImageFillContain
		r.icon = icon
		r.SetObjects(append(objects, icon))
	}

	return r
}

func (i *menuBarItem) FocusGained() {
	i.active = true
	if i.Parent.active {
		i.Parent.activateChild(i)
	}
	i.Refresh()
}

func (i *menuBarItem) FocusLost() {
	i.active = false
	i.Refresh()
}

func (i *menuBarItem) Focused() bool {
	return i.active
}

// MouseIn activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
//
// If the bar is not active, the item will be hovered.
func (i *menuBarItem) MouseIn(_ *desktop.MouseEvent) {
	i.hovered = true
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
	i.Refresh()
}

// MouseMoved activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
// This might have an effect when mouse and keyboard control are mixed.
// Changing the active menu with the keyboard will make the hovered menu bar item inactive.
// On the next mouse move the hovered item is activated again.
//
// If the bar is not active, this will do nothing.
func (i *menuBarItem) MouseMoved(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
}

// MouseOut does nothing if the bar is active.
//
// IF the bar is not active, it changes the item to not be hovered.
func (i *menuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Tapped toggles the activation state of the menu bar.
// It shows the item's menu if the bar is activated and hides it if the bar is deactivated.
func (i *menuBarItem) Tapped(*fyne.PointEvent) {
	i.Parent.toggle(i)
}

func (i *menuBarItem) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyLeft:
		if !i.Child().DeactivateLastSubmenu() {
			i.Parent.canvas.FocusPrevious()
		}
	case fyne.KeyRight:
		if !i.Child().ActivateLastSubmenu() {
			i.Parent.canvas.FocusNext()
		}
	case fyne.KeyDown:
		i.Child().ActivateNext()
	case fyne.KeyUp:
		i.Child().ActivatePrevious()
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		i.Child().TriggerLast()
	}
}

func (i *menuBarItem) TypedRune(_ rune) {
}

type menuBarItemRenderer struct {
	widget.BaseRenderer
	i          *menuBarItem
	text       *canvas.Text
	shadow     *canvas.Text
	background *canvas.Rectangle
	indicator  *canvas.Rectangle
	icon       *canvas.Image
}

func (r *menuBarItemRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	inlineIcon := theme.IconInlineSize()
	innerPad := theme.InnerPadding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.Color(theme.ColorNameForeground)
	r.text.Resize(r.text.MinSize())

	textX := padding.Width / 2
	textY := (size.Height - r.text.MinSize().Height) / 2

	if r.icon != nil {
		r.icon.Resize(fyne.NewSquareSize(inlineIcon))
		iconY := (size.Height - inlineIcon) / 2
		r.icon.Move(fyne.NewPos(textX, iconY))
		textX += inlineIcon + innerPad/2
	}

	// Тень и текст позиционируются один раз, после учёта иконки.
	// Синхронизируем метрику, чтобы при смене TextStyle (Bold в
	// активном пункте) MinSize совпадал с фактическим рендером.
	r.shadow.TextSize = r.text.TextSize
	r.shadow.TextStyle = r.text.TextStyle
	r.shadow.Resize(r.text.Size())
	r.shadow.Move(fyne.NewPos(textX+1, textY+1))
	r.text.Move(fyne.NewPos(textX, textY))

	r.background.Resize(size)

	const indicatorHeight float32 = 2
	r.indicator.Resize(fyne.NewSize(size.Width, indicatorHeight))
	r.indicator.Move(fyne.NewPos(0, size.Height-indicatorHeight))
}

func (r *menuBarItemRenderer) MinSize() fyne.Size {
	base := r.text.MinSize().Add(r.padding())
	if r.icon != nil {
		base = base.AddWidthHeight(theme.IconInlineSize()+theme.InnerPadding()/2, 0)
	}
	return base
}

func (r *menuBarItemRenderer) Refresh() {
	r.background.CornerRadius = theme.SelectionRadiusSize()

	// Тень скрыта в светлой теме целиком; в тёмной — показывается
	// всегда, независимо от состояния (isActive/hovered/default).
	// Вынесено из switch, потому что Show() вызывался только в isActive,
	// и при возврате из активного в обычное состояние тень могла остаться
	// скрытой.
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		r.shadow.Hide()
	} else {
		r.shadow.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 0x90}
		r.shadow.TextStyle = r.text.TextStyle
		r.shadow.Show()
	}

	accent := theme.Color(theme.ColorNameMenuBarAccent)
	isActive := r.i.active && r.i.Parent.active

	switch {
	case isActive:
		r.background.FillColor = theme.Color(theme.ColorNameMenuBarActiveBg)
		r.background.Show()
		r.indicator.FillColor = accent
		r.indicator.Show()
		r.text.Color = accent
		r.text.TextStyle = fyne.TextStyle{Bold: true}

	case r.i.hovered && !r.i.Parent.active:
		r.background.FillColor = theme.Color(theme.ColorNameMenuBarHoverBg)
		r.background.Show()
		r.indicator.FillColor = accent
		r.indicator.Show()
		r.text.Color = theme.Color(theme.ColorNameForeground)
		r.text.TextStyle = fyne.TextStyle{}

	default:
		r.background.Hide()
		r.indicator.Hide()
		r.text.Color = theme.Color(theme.ColorNameForeground)
		r.text.TextStyle = fyne.TextStyle{}
	}

	// Синхронизируем TextStyle тени с текстом. Ветки switch выше меняют
	// Bold у text, а тень всегда должна иметь ту же метрику, иначе
	// «двоит».
	r.shadow.TextStyle = r.text.TextStyle

	r.background.Refresh()
	r.indicator.Refresh()
	r.shadow.Refresh()
	r.text.Refresh()
	canvas.Refresh(r.i)
}

func (r *menuBarItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.InnerPadding()*2, theme.InnerPadding())
}
