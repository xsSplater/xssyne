//go:build !wasm && !test_web_driver

// xssyne/internal/driver/glfw/menu_notweb.go

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

func addMissingQuitForMainMenu(menus *fyne.MainMenu, w *window) *fyne.MainMenu {
	// Ищем первый не-сепараторный пункт верхнего меню.
	var first *fyne.Menu
	for _, m := range menus.Items {
		if !m.IsSeparator {
			first = m
			break
		}
	}
	if first == nil {
		return menus
	}
	localQuit := lang.L("Quit")
	var lastItem *fyne.MenuItem
	if len(first.Items) > 0 {
		lastItem = first.Items[len(first.Items)-1]
		if lastItem.Label == localQuit {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit {
		quitItem := fyne.NewMenuItem(localQuit, nil)
		quitItem.IsQuit = true
		first.Items = append(first.Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range first.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = func() {
				for _, win := range w.driver.AllWindows() {
					if glWin, ok := win.(*window); ok {
						glWin.closed(glWin.view())
					} else {
						win.Close()
					}
				}
			}
		}
	}
	return menus
}
