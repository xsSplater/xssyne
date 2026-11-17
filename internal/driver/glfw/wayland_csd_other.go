//go:build windows || (linux && ((x11 && !wayland) || (!x11 && wayland)))

package glfw

// forcePlatform returns platformAuto when either x11 or wayland was specified
// or if we're not on a suitable OS.
func forcePlatform() string {
	return platformAuto
}
