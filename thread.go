// xssyne/thread.go

package fyne

// DoAndWait is used to execute a specified function in the main Fyne runtime context.
// This is required when a background process wishes to adjust graphical elements of a running app.
// Developers should use this only from within goroutines they have created.
//
// Since: 2.6
func DoAndWait(fn func()) {
	app := CurrentApp()
	if app == nil {
		return // CurrentApp уже залогировал "no app started"
	}
	driver := app.Driver()
	if driver == nil {
		LogError("xssyne.DoAndWait called when driver is not available", nil)
		return
	}
	driver.DoFromGoroutine(fn, true)
}

// Do is used to execute a specified function in the main Fyne runtime context without waiting.
// This is required when a background process wishes to adjust graphical elements of a running app.
// Developers should use this only from within goroutines they have created and when the result does not have to
// be waited for.
//
// Since: 2.6
func Do(fn func()) {
	app := CurrentApp()
	if app == nil {
		return
	}
	driver := app.Driver()
	if driver == nil {
		LogError("xssyne.Do called when driver is not available", nil)
		return
	}
	driver.DoFromGoroutine(fn, false)
}
