// xssyne/app/meta.go

package app

import (
	"sync"

	"fyne.io/fyne/v2"
)

var meta = fyne.AppMetadata{
	ID:         "",
	Name:       "",
	Version:    "0.0.1",
	Build:      1,
	Release:    false,
	Custom:     map[string]string{},
	Migrations: map[string]bool{},
}

// SetMetadata overrides the packaged application metadata.
// This data can be used in many places like notifications and about screens.
func SetMetadata(m fyne.AppMetadata) {
	meta = m

	if meta.Custom == nil {
		meta.Custom = map[string]string{}
	}
	if meta.Migrations == nil {
		meta.Migrations = map[string]bool{}
	}
}

// metaCheckOnce гарантирует, что поиск FyneApp.toml выполняется максимум
// один раз за время жизни процесса. В dev-режиме без файла Metadata()
// вызывается часто (Quit, Preferences, notificationAppID, ...) и каждый
// раз ходил на диск.
var metaCheckOnce sync.Once

func (a *fyneApp) Metadata() fyne.AppMetadata {
	if meta.ID == "" && meta.Name == "" {
		metaCheckOnce.Do(checkLocalMetadata)
	}

	return meta
}
