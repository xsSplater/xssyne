//go:build windows

// xssyne/internal/app/config_windows.go

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "fyne")
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "AppData", "Roaming", "fyne")
	}

	return filepath.Join(os.TempDir(), "fyne")
}
