// xssyne/internal/build/build.go

// Package build contains information about they type of build currently running.
package build

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

var (
	migrateCheck sync.Once

	migratedFyneDo  bool
	migrateResolved atomic.Bool
)

func MigratedToFyneDo() bool {
	if DisableThreadChecks {
		return true
	}

	// Быстрый путь: уже определились с реальным app.
	if migrateResolved.Load() {
		return migratedFyneDo
	}

	// Не кэшируем "false", пока app не создан — иначе первый вызов
	// из раннего контекста навсегда зафиксирует "не мигрировали",
	// даже если в metadata потом окажется fyneDo=true.
	app := fyne.CurrentApp()
	if app == nil {
		return false
	}

	migrateCheck.Do(func() {
		v := app.Metadata().Migrations["fyneDo"]
		migratedFyneDo = v
		migrateResolved.Store(true)
	})

	return migratedFyneDo
}
