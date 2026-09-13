package cache

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var renderers async.Map[fyne.Widget, *rendererInfo]

type isBaseWidget interface {
	ExtendBaseWidget(fyne.Widget)
	super() fyne.Widget
}

// Renderer looks up the render implementation for a widget
// If one does not exist, it creates and caches a renderer for the widget.
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	if wid == nil {
		return nil
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	if rinfo, ok := renderers.Load(wid); ok {
		rinfo.setAlive()
		return rinfo.renderer
	}

	renderer := wid.CreateRenderer()
	rinfo := &rendererInfo{renderer: renderer}
	rinfo.setAlive()

	actual, loaded := renderers.LoadOrStore(wid, rinfo)
	if loaded {
		// проиграли гонку: наш рендерер не пригодился, освобождаем ресурсы
		renderer.Destroy()
		actual.setAlive()
		return actual.renderer
	}
	return renderer
}

// CachedRenderer looks up the cached render implementation for a widget
// If a renderer does not exist in the cache, it returns nil, false.
func CachedRenderer(wid fyne.Widget) (fyne.WidgetRenderer, bool) {
	if wid == nil {
		return nil, false
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	rinfo, ok := renderers.Load(wid)
	if !ok {
		return nil, false
	}

	rinfo.setAlive()
	return rinfo.renderer, true
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	rinfo, ok := renderers.LoadAndDelete(wid)
	if !ok {
		return
	}
	if rinfo != nil {
		rinfo.renderer.Destroy()
	}
	overrides.Delete(wid)
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid fyne.Widget) bool {
	_, found := renderers.Load(wid)
	return found
}

type rendererInfo struct {
	expiringCache
	renderer fyne.WidgetRenderer
}
