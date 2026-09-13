// xssyne/internal/driver/glfw/menu_bar.go

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*MenuBar)(nil)

// MenuBar is a widget for displaying a fyne.MainMenu in a bar.
type MenuBar struct {
	widget.Base
	Items []fyne.CanvasObject // *menuBarItem или *menuBarSeparator

	active     bool
	activeItem *menuBarItem
	canvas     fyne.Canvas
}

// NewMenuBar creates a menu bar populated with items from the passed main menu structure.
func NewMenuBar(mainMenu *fyne.MainMenu, canvas fyne.Canvas) *MenuBar {
	b := &MenuBar{canvas: canvas}
	b.ExtendBaseWidget(b)
	for _, menu := range mainMenu.Items {
		if menu.IsSeparator {
			sep := &menuBarSeparator{}
			sep.ExtendBaseWidget(sep)
			b.Items = append(b.Items, sep)
			continue
		}
		barItem := &menuBarItem{Menu: menu, Parent: b}
		barItem.ExtendBaseWidget(barItem)
		b.Items = append(b.Items, barItem)
	}
	return b
}

// CreateRenderer returns a new renderer for the menu bar.
func (b *MenuBar) CreateRenderer() fyne.WidgetRenderer {
	cont := container.NewWithoutLayout()
	for _, item := range b.Items {
		cont.Add(item)
	}
	background := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
	widget.ApplyShadowForLevel(&background.Shadow, widget.MenuBarLevel, theme.Color(theme.ColorNameShadow))

	// Тонкая линия под полосой меню — визуально отделяет bar от
	// рабочей области окна. Рисуется над background, но под cont,
	// чтобы не перекрывать содержимое.
	underline := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))

	underlay := &menuBarUnderlay{action: b.deactivate}
	underlay.ExtendBaseWidget(underlay)
	objects := []fyne.CanvasObject{underlay, background, underline, cont}
	for _, item := range b.Items {
		if barItem, ok := item.(*menuBarItem); ok {
			objects = append(objects, barItem.Child())
		}
	}
	return &menuBarRenderer{ // Отражает порядок слоёв
		widget.NewBaseRenderer(objects),
		b,
		background,
		underline,
		underlay,
		cont,
	}
}

// IsActive returns whether the menu bar is active or not.
// An active menu bar shows the current selected menu and should have the focus.
func (b *MenuBar) IsActive() bool {
	return b.active
}

// Toggle changes the activation state of the menu bar.
// On activation, the first item will become active.
func (b *MenuBar) Toggle() {
	first := b.firstMenuBarItem()
	if first != nil {
		b.toggle(first)
	}
}

// firstMenuBarItem возвращает первый не-сепараторный пункт.
func (b *MenuBar) firstMenuBarItem() *menuBarItem {
	for _, item := range b.Items {
		if barItem, ok := item.(*menuBarItem); ok {
			return barItem
		}
	}
	return nil
}

func (b *MenuBar) activateChild(item *menuBarItem) {
	b.active = true
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if b.activeItem == item {
		return
	}

	if b.activeItem != nil {
		if c := b.activeItem.Child(); c != nil {
			c.Hide()
		}
		b.activeItem.Refresh()
	}
	b.activeItem = item
	if item == nil {
		return
	}

	item.Refresh()
	item.Child().Show()
	b.Refresh()
}

func (b *MenuBar) deactivate() {
	if !b.active {
		return
	}

	b.active = false
	if b.activeItem != nil {
		if c := b.activeItem.Child(); c != nil {
			defer c.Dismiss()
			c.Hide()
		}
		b.activeItem.Refresh()
		b.activeItem = nil
	}
	b.Refresh()
}

func (b *MenuBar) toggle(item *menuBarItem) {
	if b.active {
		b.canvas.Unfocus()
		b.deactivate()
	} else {
		b.activateChild(item)
		b.canvas.Focus(item)
	}
}

type menuBarRenderer struct {
	widget.BaseRenderer
	b          *MenuBar
	background *canvas.Rectangle
	underline  *canvas.Rectangle
	underlay   *menuBarUnderlay
	cont       *fyne.Container
}

func (r *menuBarRenderer) Layout(size fyne.Size) {
	minSize := r.MinSize()
	if size.Height != minSize.Height || size.Width < minSize.Width {
		r.b.Resize(fyne.NewSize(fyne.Max(size.Width, minSize.Width), minSize.Height))
		return
	}

	if r.b.active {
		r.underlay.Resize(r.b.canvas.Size())
	} else {
		r.underlay.Resize(fyne.NewSize(0, 0))
	}
	innerPadding := theme.InnerPadding()
	r.cont.Resize(fyne.NewSize(size.Width-2*innerPadding, size.Height))
	r.cont.Move(fyne.NewPos(innerPadding, 0))

	x := float32(0)
	for _, item := range r.b.Items {
		min := item.MinSize()
		item.Resize(fyne.NewSize(min.Width, size.Height))
		item.Move(fyne.NewPos(x, 0))
		x += min.Width
	}

	if item := r.b.activeItem; item != nil {
		if item.Child().Size().IsZero() {
			item.Child().Resize(item.Child().MinSize())
		}
		item.Child().Move(fyne.NewPos(item.Position().X+innerPadding, item.Size().Height))
	}

	r.background.Move(fyne.NewPos(0, 0))
	r.background.Resize(size)

	const underlineHeight float32 = 1
	r.underline.Resize(fyne.NewSize(size.Width, underlineHeight))
	r.underline.Move(fyne.NewPos(0, size.Height-underlineHeight))
}

func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(fyne.NewSize(theme.InnerPadding()*2, 0))
}

func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	r.background.FillColor = theme.Color(theme.ColorNameBackground)
	r.background.Shadow.Color = theme.Color(theme.ColorNameShadow)
	r.background.Refresh()

	r.underline.FillColor = theme.Color(theme.ColorNameSeparator)
	r.underline.Refresh()

	canvas.Refresh(r.b)
}

// Transparent underlay shown as soon as menu is active.
// It catches mouse events outside the menu's objects.
type menuBarUnderlay struct {
	widget.Base
	action func()
}

var (
	_ fyne.Widget       = (*menuBarUnderlay)(nil)
	_ fyne.Tappable     = (*menuBarUnderlay)(nil) // deactivate menu on click outside
	_ desktop.Hoverable = (*menuBarUnderlay)(nil) // block hover events on main content
)

func (u *menuBarUnderlay) CreateRenderer() fyne.WidgetRenderer {
	return &menuUnderlayRenderer{}
}

func (u *menuBarUnderlay) MouseIn(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) MouseOut() {
}

func (u *menuBarUnderlay) MouseMoved(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) Tapped(*fyne.PointEvent) {
	u.action()
}

type menuUnderlayRenderer struct {
	widget.BaseRenderer
}

var _ fyne.WidgetRenderer = (*menuUnderlayRenderer)(nil)

func (r *menuUnderlayRenderer) Layout(fyne.Size) {
}

func (r *menuUnderlayRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *menuUnderlayRenderer) Refresh() {
}
