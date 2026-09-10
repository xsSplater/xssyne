# xssyne

**xssyne** is a lightweight fork of the [Fyne](https://fyne.io) GUI toolkit, tailored specifically for the **Darktide-Servo-ModQuisitor-2** project – a mod manager for *Warhammer 40,000: Darktide*.

The fork removes all mobile, web, and macOS-specific components, leaving only the code needed to run on **Windows** and **Linux**. This reduces repository size, simplifies maintenance, and speeds up compilation.

## Features

- Full set of widgets (buttons, entries, lists, tables, trees, tabs, forms, dialogs, etc.)
- Light and dark theme support
- OpenGL (GLFW) rendering
- Single codebase for Windows and Linux
- All non-desktop platform code stripped out

## Why this fork

This fork is used in [Darktide-Servo-ModQuisitor-2](https://github.com/xsSplater/Darktide-Servo-Modquisitor-2) – a convenient mod manager for Darktide that allows:
- Drag-and-drop mod installation
- One-click automatic problem detection and fixes (outdated mods, conflicts, dependencies)
- Manual or automatic mod list sorting

Using a trimmed-down version of Fyne makes the application compact and fast, without unnecessary dependencies.

## Requirements

- Go 1.22 or newer (the project uses Go 1.27.1)
- C compiler (gcc or clang)
- GLFW system headers

## Quick start (for developers)

If you want to use xssyne in your own project, replace the import in your `go.mod`:
```go
replace fyne.io/fyne/v2 => /path/to/your/local/xssyne
```

Example from the main project:
```go
// go.mod
module Servo-Modquisitor

go 1.27.1

require (
    fyne.io/fyne/v2 v2.8.1
    // ...
)

replace fyne.io/fyne/v2 => A:/GitHub/xssyne
```

Then build as usual:
go mod tidy
go build
To build for a specific OS, set the environment variables:
GOOS=windows GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go build

# License
This project is a derivative of Fyne and is distributed under the BSD 3-Clause License, retaining the original copyright notice.

Copyright © 2018 Fyne.io developers.
All rights reserved.
The full license text is available in the LICENSE file.

Links
Original Fyne: https://fyne.io
Original source code: https://github.com/fyne-io/fyne
Main project using xssyne: [Darktide-Servo-ModQuisitor-2](https://github.com/xsSplater/Darktide-Servo-Modquisitor-2)
