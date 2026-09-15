//go:build !windows

// xssyne/internal/driver/glfw/driver_notwindows.go

package glfw

func isDark() bool {
	return true // this is really a no-op placeholder for a windows menu workaround
}
