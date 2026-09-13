// Package build contains information about they type of build currently running.
package build

import (
	"sync"

	"fyne.io/fyne/v2"
)

var (
	migrateCheck sync.Once

	migratedFyneDo bool
)

func MigratedToFyneDo() bool {
	if DisableThreadChecks {
		return true
	}

	migrateCheck.Do(func() {
		app := fyne.CurrentApp()
		if app == nil {
			return // no app yet; assume not migrated
		}
		v, ok := app.Metadata().Migrations["fyneDo"]
		if ok {
			migratedFyneDo = v
		}
	})

	return migratedFyneDo
}
