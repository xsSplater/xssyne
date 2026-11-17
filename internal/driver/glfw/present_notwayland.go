//go:build windows || (linux && x11)

package glfw

func newPresentGate(_ *window) presentGate { return noGate{} }
