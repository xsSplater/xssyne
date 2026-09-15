//go:build !ci && !wasm && test_web_driver

// xssyne/app/app_openurl_web.go

package app

import (
	"errors"
	"net/url"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("OpenURL is not supported with the test web driver.")
}
