//go:build !linux && !windows && !darwin && !freebsd && !openbsd && !netbsd

// xssyne/app/theme_other.go

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/theme"
)

func DefaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}
