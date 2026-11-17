//go:build !ci && !wasm && !test_web_driver && !android && !ios && !mobile && (linux || openbsd || freebsd || netbsd) && !noos && !tinygo

// xssyne/internal/app/config_xdg.go

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "fyne")
	}

	// $XDG_CONFIG_HOME не определён, $HOME отсутствует или нечитаем —
	// filepath.Join с пустой строкой даёт относительный путь. Fallback
	// на temp, чтобы приложение не писало prefs в CWD.
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "fyne")
	}
	return filepath.Join(os.TempDir(), "fyne")
}
