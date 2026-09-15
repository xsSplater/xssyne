//go:build windows && mobile

// xssyne/app/app_windows_stub.go

package app

import (
	"log"

	"fyne.io/fyne/v2"
)

func NewWithID(_ string) fyne.App {
	log.Fatal("Cannot launch the mobile simulator mode on Windows")
	return nil
}
