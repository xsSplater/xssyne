//go:build !windows && !linux

// xssyne/internal/app/theme_other.go

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/theme"
)

// DefaultVariant returns the systems default fyne.ThemeVariant.
// xssyne supports only Windows and Linux — on all other platforms
// (darwin, aix, bsd, android, …) we fall back to dark.
func DefaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}
