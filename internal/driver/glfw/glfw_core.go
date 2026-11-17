//go:build (linux || windows) && !gles && !arm && !arm64

package glfw

import "github.com/go-gl/glfw/v3.4/glfw"

func initWindowHints() {
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.CocoaGraphicsSwitching, glfw.True)
}
