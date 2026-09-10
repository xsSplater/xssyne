package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) cachedIconPath() string {
	if a.Icon() == nil {
		return ""
	}

	a.iconOnce.Do(func() {
		dirPath := rootCacheDir(a)
		filePath := filepath.Join(dirPath, "icon.png")
		if err := a.saveIconToCache(dirPath, filePath); err != nil {
			return
		}
		a.iconPath = filePath
	})

	return a.iconPath
}

func (a *fyneApp) saveIconToCache(dirPath, filePath string) error {
	err := os.MkdirAll(dirPath, 0o700)
	if err != nil {
		fyne.LogError("Unable to create application cache directory", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fyne.LogError("Unable to create icon file", err)
		return err
	}

	defer file.Close()

	if icon := a.Icon(); icon != nil {
		_, err = file.Write(icon.Content())
		if err != nil {
			fyne.LogError("Unable to write icon contents", err)
			return err
		}
	}

	return nil
}
