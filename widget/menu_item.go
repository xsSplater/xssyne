// xssyne/widget/menu.go

package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	fynecolor "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var (
	_ fyne.Widget       = (*menuItem)(nil)
	_ desktop.Hoverable = (*menuItem)(nil)
	_ fyne.Tappable     = (*menuItem)(nil)
)

// menuItem is a widget for displaying a fyne.menuItem.
type menuItem struct {
	widget.Base
	Item *fyne.MenuItem

	alignment     fyne.TextAlign
	child, parent *Menu
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *fyne.MenuItem, parent *Menu) *menuItem {
	i := &menuItem{Item: item, parent: parent}
	i.alignment = parent.alignment
	i.ExtendBaseWidget(i)
	return i
}

func (i *menuItem) Child() *Menu {
	if i.Item.ChildMenu != nil && i.child == nil {
		child := NewMenu(i.Item.ChildMenu)
		child.Hide()
		child.OnDismiss = i.parent.Dismiss
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu item.
func (i *menuItem) CreateRenderer() fyne.WidgetRenderer {
	th := i.parent.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	background := canvas.NewRectangle(th.Color(theme.ColorNameHover, v))
	background.CornerRadius = th.Size(theme.SizeNameMenuRadius)
	background.StrokeWidth = th.Size(theme.SizeNameMenuBorderWidth)
	background.StrokeColor = th.Color(theme.ColorNameMenuItemActiveBorder, v)
	background.Hide()
	text := canvas.NewText(i.Item.Label, th.Color(theme.ColorNameForeground, v))
	text.Alignment = i.alignment
	// Тень — тот же текст, что и основной, но смещённый на 1 px вниз и вправо
	// и с пониженной альфой. Рисуется ПОД основным текстом (порядок
	// в objects: background, shadow, text).
	shadow := canvas.NewText(i.Item.Label, color.NRGBA{R: 0, G: 0, B: 0, A: 200})
	shadow.Alignment = i.alignment
	shadow.TextStyle = text.TextStyle
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		shadow.Hide()
	}
	objects := []fyne.CanvasObject{background, shadow, text}

	var subtext *canvas.Text
	if i.Item.Subtitle != "" {
		subtext = canvas.NewText(i.Item.Subtitle, th.Color(theme.ColorNamePlaceHolder, v))
		subtext.TextSize = th.Size(theme.SizeNameCaptionText)
		subtext.Alignment = i.alignment
		objects = append(objects, subtext)
	}

	var expandIcon *canvas.Image
	if i.Item.ChildMenu != nil {
		expandIcon = canvas.NewImageFromResource(th.Icon(theme.IconNameMenuExpand))
		objects = append(objects, expandIcon)
	}
	checkIcon := canvas.NewImageFromResource(th.Icon(theme.IconNameConfirm))
	if !i.Item.Checked {
		checkIcon.Hide()
	}
	var icon *canvas.Image
	if i.Item.Icon != nil {
		icon = canvas.NewImageFromResource(i.Item.Icon)
		objects = append(objects, icon)
	}
	var shortcutTexts []*canvas.Text
	if s, ok := i.Item.Shortcut.(fyne.KeyboardShortcut); ok {
		shortcutTexts = textsForShortcut(s, th)
		for _, t := range shortcutTexts {
			objects = append(objects, t)
		}
	}

	objects = append(objects, checkIcon)
	r := &menuItemRenderer{
		BaseRenderer:  widget.NewBaseRenderer(objects),
		i:             i,
		expandIcon:    expandIcon,
		checkIcon:     checkIcon,
		icon:          icon,
		shortcutTexts: shortcutTexts,
		text:          text,
		subtext:       subtext,
		shadow:        shadow,
		background:    background,
	}
	r.updateVisuals()
	return r
}

// MouseIn activates the item which shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
func (i *menuItem) MouseIn(e *desktop.MouseEvent) {
	if i.Item.Header {
		return
	}
	i.activate()
}

// MouseMoved does nothing.
func (i *menuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut deactivates the item unless it has an open submenu.
func (i *menuItem) MouseOut() {
	if !i.isSubmenuOpen() {
		i.deactivate()
	}
}

// Tapped performs the action of the item and dismisses the menu.
// It does nothing if the item doesn't have an action.
func (i *menuItem) Tapped(*fyne.PointEvent) {
	if i.Item.Disabled || i.Item.Header {
		return
	}
	if i.Item.Action == nil {
		if fyne.CurrentDevice().IsMobile() {
			i.activate()
		}

		return
	} else if i.Item.ChildMenu != nil {
		if fyne.CurrentDevice().IsMobile() {
			i.activate()
		}

		return
	}

	i.trigger()
}

func (i *menuItem) activate() {
	if i.Item.Disabled || i.Item.Header {
		return
	}
	if i.Child() != nil {
		i.Child().Show()
	}
	i.parent.activateItem(i)
}

func (i *menuItem) activateLastSubmenu() bool {
	if i.Child() == nil {
		return false
	}
	if i.isSubmenuOpen() {
		return i.Child().ActivateLastSubmenu()
	}
	i.Child().Show()
	i.Child().ActivateNext()
	return true
}

func (i *menuItem) deactivate() {
	if i.Child() != nil {
		i.Child().Hide()
	}
	i.parent.DeactivateChild()
}

func (i *menuItem) deactivateLastSubmenu() bool {
	if !i.isSubmenuOpen() {
		return false
	}
	if !i.Child().DeactivateLastSubmenu() {
		i.Child().DeactivateChild()
		i.Child().Hide()
	}
	return true
}

func (i *menuItem) isActive() bool {
	return i.parent.activeItem == i
}

func (i *menuItem) isSubmenuOpen() bool {
	return i.Child() != nil && i.Child().Visible()
}

func (i *menuItem) trigger() {
	i.parent.Dismiss()
	if i.Item.Action != nil {
		i.Item.Action()
	}
}

func (i *menuItem) triggerLast() {
	if i.isSubmenuOpen() {
		i.Child().TriggerLast()
		return
	}
	i.trigger()
}

type menuItemRenderer struct {
	widget.BaseRenderer
	i                *menuItem
	background       *canvas.Rectangle
	checkIcon        *canvas.Image
	expandIcon       *canvas.Image
	icon             *canvas.Image
	lastThemePadding float32
	minSize          fyne.Size
	shortcutTexts    []*canvas.Text
	text             *canvas.Text
	subtext          *canvas.Text
	shadow           *canvas.Text
}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	th := r.i.parent.Theme()
	innerPad := th.Size(theme.SizeNameInnerPadding)
	pad := th.Size(theme.SizeNamePadding)
	inlineIcon := th.Size(theme.SizeNameInlineIcon)

	leftOffset := innerPad + r.checkSpace()
	rightOffset := size.Width
	iconSize := fyne.NewSquareSize(inlineIcon)
	iconTopOffset := (size.Height - inlineIcon) / 2

	if r.expandIcon != nil {
		rightOffset -= inlineIcon
		r.expandIcon.Resize(iconSize)
		r.expandIcon.Move(fyne.NewPos(rightOffset, iconTopOffset))
	}

	rightOffset -= innerPad
	for i := len(r.shortcutTexts) - 1; i >= 0; i-- {
		text := r.shortcutTexts[i]
		text.Resize(text.MinSize())
		rightOffset -= text.MinSize().Width
		text.Move(fyne.NewPos(rightOffset, innerPad))
		if i == 0 {
			rightOffset -= innerPad
		}
	}

	r.checkIcon.Resize(iconSize)
	r.checkIcon.Move(fyne.NewPos(innerPad, iconTopOffset))

	if r.icon != nil {
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(leftOffset, iconTopOffset))
		leftOffset += inlineIcon + innerPad
	}

	textWidth := rightOffset - leftOffset

	if r.subtext != nil {
		// Две строки: основная сверху, подзаголовок снизу.
		mainHeight := r.text.MinSize().Height
		subHeight := r.subtext.MinSize().Height
		totalHeight := mainHeight + subHeight
		top := (size.Height - totalHeight) / 2

		r.text.Resize(fyne.NewSize(textWidth, mainHeight))
		r.text.Move(fyne.NewPos(leftOffset, top))

		r.subtext.Resize(fyne.NewSize(textWidth, subHeight))
		r.subtext.Move(fyne.NewPos(leftOffset, top+mainHeight))
	} else {
		textHeight := r.text.MinSize().Height
		textY := (size.Height - textHeight) / 2
		r.text.Resize(fyne.NewSize(textWidth, textHeight))
		r.text.Move(fyne.NewPos(leftOffset, textY))
	}

	r.background.Resize(size.Subtract(fyne.NewSquareSize(pad)))
	r.background.Move(fyne.NewPos(pad/2, pad/2))

	if r.shadow != nil {
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			r.shadow.Hide()
		} else {
			// Синхронизируем метрику тени с основным текстом — иначе
			// Bold-тень не совпадёт с regular-текстом и «двоит» буквы.
			r.shadow.TextSize = r.text.TextSize
			r.shadow.TextStyle = r.text.TextStyle
			r.shadow.Resize(r.text.Size())
			pos := r.text.Position()
			r.shadow.Move(fyne.NewPos(pos.X+1, pos.Y+1))

			r.shadow.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 0x90}
			r.shadow.Show()
		}
		r.shadow.Refresh()
	}
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	th := r.i.parent.Theme()
	innerPad := th.Size(theme.SizeNameInnerPadding)
	inlineIcon := th.Size(theme.SizeNameInlineIcon)
	innerPad2 := innerPad * 2

	minSize := r.text.MinSize().AddWidthHeight(innerPad2+r.checkSpace(), innerPad2)
	if r.subtext != nil {
		minSize = minSize.AddWidthHeight(0, r.subtext.MinSize().Height)
	}
	if r.expandIcon != nil {
		minSize = minSize.AddWidthHeight(inlineIcon, 0)
	}
	if r.icon != nil {
		minSize = minSize.AddWidthHeight(inlineIcon+innerPad, 0)
	}
	if r.shortcutTexts != nil {
		var textWidth float32
		for _, text := range r.shortcutTexts {
			textWidth += text.MinSize().Width
		}
		minSize = minSize.AddWidthHeight(textWidth+innerPad, 0)
	}
	r.minSize = minSize
	r.lastThemePadding = innerPad
	return r.minSize
}

func (r *menuItemRenderer) updateVisuals() {
	th := r.i.parent.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.background.CornerRadius = th.Size(theme.SizeNameMenuRadius)
	r.background.StrokeWidth = th.Size(theme.SizeNameMenuBorderWidth)
	r.background.StrokeColor = th.Color(theme.ColorNameMenuItemActiveBorder, v)

	switch {
	case r.i.Item.Header:
		// Header — постоянный фон, без hover/active.
		r.background.FillColor = th.Color(theme.ColorNameMenuItemHeaderBg, v)
		r.background.Show()
	case fyne.CurrentDevice().IsMobile():
		r.background.Hide()
	case r.i.isActive():
		r.background.FillColor = th.Color(theme.ColorNameFocus, v)
		r.background.Show()
	default:
		r.background.Hide()
	}

	r.background.Refresh()
	r.text.Alignment = r.i.alignment
	r.refreshText(r.text, false)
	if r.shadow != nil {
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			r.shadow.Hide()
		} else {
			r.shadow.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 0x90}
			r.shadow.TextStyle = r.text.TextStyle
			r.shadow.Show()
		}
		r.shadow.Refresh()
	}
	for _, text := range r.shortcutTexts {
		r.refreshText(text, true)
	}
	if r.subtext != nil {
		r.subtext.TextSize = th.Size(theme.SizeNameCaptionText)
		r.subtext.Color = th.Color(theme.ColorNamePlaceHolder, v)
		r.subtext.Refresh()
	}

	// Кастомная иконка выбора: приоритет у CheckedIcon/UncheckedIcon
	// из Item, если они заданы. Иначе — стандартная галочка Fyne.
	// updateIcon сюда НЕ вызываем: он бы перезаписал кастомный ресурс
	// значением из темы.
	switch {
	case r.i.Item.Checked && r.i.Item.CheckedIcon != nil:
		r.checkIcon.Resource = r.i.Item.CheckedIcon
		r.checkIcon.Show()
	case !r.i.Item.Checked && r.i.Item.UncheckedIcon != nil:
		r.checkIcon.Resource = r.i.Item.UncheckedIcon
		r.checkIcon.Show()
	case r.i.Item.Checked:
		r.checkIcon.Resource = th.Icon(theme.IconNameConfirm)
		r.checkIcon.Show()
	default:
		r.checkIcon.Hide()
	}
	r.updateIcon(r.expandIcon, th.Icon(theme.IconNameMenuExpand))
	r.updateIcon(r.icon, r.i.Item.Icon)
}

func (r *menuItemRenderer) Refresh() {
	r.updateVisuals()
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) checkSpace() float32 {
	if r.i.parent.containsCheck {
		return theme.IconInlineSize() + theme.InnerPadding()
	}
	return 0
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	th := r.i.parent.Theme()

	return !r.minSize.IsZero() &&
		r.text.TextSize == th.Size(theme.SizeNameText) &&
		(r.expandIcon == nil || r.expandIcon.Size().Width == th.Size(theme.SizeNameInlineIcon)) &&
		(r.subtext == nil || r.subtext.TextSize == th.Size(theme.SizeNameCaptionText)) &&
		r.lastThemePadding == th.Size(theme.SizeNameInnerPadding)
}

func (r *menuItemRenderer) updateIcon(img *canvas.Image, rsc fyne.Resource) {
	if img == nil {
		return
	}
	if r.i.Item.Disabled && svg.IsResourceSVG(rsc) {
		img.Resource = theme.NewDisabledResource(rsc)
	} else {
		img.Resource = rsc
	}
}

func (r *menuItemRenderer) refreshText(text *canvas.Text, shortcut bool) {
	th := r.i.parent.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	text.TextSize = th.Size(theme.SizeNameText)
	switch {
	case r.i.Item.Disabled:
		text.Color = th.Color(theme.ColorNameDisabled, v)
		text.TextStyle = fyne.TextStyle{}
	case r.i.Item.Danger:
		text.Color = th.Color(theme.ColorNameMenuItemDanger, v)
		text.TextStyle = fyne.TextStyle{Bold: true}
	case r.i.Item.Header:
		text.Color = th.Color(theme.ColorNameMenuItemHeader, v)
		text.TextStyle = fyne.TextStyle{Bold: true}
	case shortcut:
		text.Color = shortcutColor(th)
		text.TextStyle = fyne.TextStyle{}
	default:
		text.Color = th.Color(theme.ColorNameForeground, v)
		text.TextStyle = fyne.TextStyle{}
	}
	text.Refresh()
}

func shortcutColor(th fyne.Theme) color.Color {
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r, g, b, a := fynecolor.ToNRGBA(th.Color(theme.ColorNameForeground, v))
	a = uint8(float32(a) * 0.198)
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func textsForShortcut(sc fyne.KeyboardShortcut, th fyne.Theme) (texts []*canvas.Text) {
	// add modifier
	b := strings.Builder{}
	mods := sc.Mod()
	if mods&fyne.KeyModifierControl != 0 {
		b.WriteString(textModifierControl)
	}
	if mods&fyne.KeyModifierAlt != 0 {
		b.WriteString(textModifierAlt)
	}
	if mods&fyne.KeyModifierShift != 0 {
		b.WriteString(textModifierShift)
	}
	if mods&fyne.KeyModifierSuper != 0 {
		b.WriteString(textModifierSuper)
	}
	shortColor := shortcutColor(th)
	if b.Len() > 0 {
		t := canvas.NewText(b.String(), shortColor)
		t.TextStyle = styleModifiers
		texts = append(texts, t)
	}
	// add key
	style := defaultStyleKeys
	s, ok := keyTexts[sc.Key()]
	if !ok {
		s = string(sc.Key())
	} else if len(s) == 1 {
		style = fyne.TextStyle{Symbol: true}
	}
	t := canvas.NewText(s, shortColor)
	t.TextStyle = style
	texts = append(texts, t)
	return texts
}
