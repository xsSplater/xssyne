//go:build noos || tinygo

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "fyne")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fyne")
}
