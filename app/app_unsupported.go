//go:build !windows && !linux && !ci && !wasm && !test_web_driver && !tamago && !noos && !tinygo

// xssyne/app/app_unsupported.go

package app

import (
	"errors"
	"log"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
)

// xssyne ships drivers only for Windows and Linux. This file exists so
// that `go vet ./...` and gopls (which scan every GOOS by default) do
// not fail with cascading "undefined:" errors. Building for any other
// platform will hit log.Fatal in NewWithID at startup.

func NewWithID(_ string) fyne.App {
	log.Fatal("xssyne supports only windows and linux — no driver for this platform")
	return nil
}

func (a *fyneApp) OpenURL(_ *url.URL) error {
	return errors.New("xssyne supports only windows and linux")
}

func (a *fyneApp) SendNotification(_ *fyne.Notification) {
	fyne.LogError("xssyne does not support notifications on this platform", nil)
}

func (a *fyneApp) ScheduleNotification(_ *fyne.Notification, _ time.Time) (*fyne.ScheduledNotification, error) {
	return nil, errors.New("xssyne supports only windows and linux")
}

func (a *fyneApp) CancelScheduledNotification(_ string) error {
	return errors.New("xssyne supports only windows and linux")
}

func watchTheme(_ *settings) {
	// no-op
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
