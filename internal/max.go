package internal

import "fyne.io/fyne/v2"

// MaxSizes returns the maximum of two fyne.Size values.
func MaxSizes(a, b fyne.Size) fyne.Size {
	return fyne.NewSize(
		fyne.Max(a.Width, b.Width),
		fyne.Max(a.Height, b.Height),
	)
}
