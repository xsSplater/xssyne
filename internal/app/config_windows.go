//go:build !ci && !android && !ios && !wasm && !test_web_driver && !noos && !tinygo

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "fyne")
	}

	// fallback only if UserConfigDir errored (extremely rare on Windows)
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "AppData", "Roaming", "fyne")
}
