//go:build (linux || freebsd || openbsd || netbsd) && !android && !wasm && !js

// xssyne/dialog/file_unix_dialog.go

package dialog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// getFavoriteLocation читает путь к стандартной папке (Downloads,
// Documents, ...) с учётом локализации. Сначала пробует XDG-переменную
// из user-dirs.dirs, потом — англоязычное имя в $HOME.
//
// Fyne 2.8 использует жёстко зашитые английские имена, из-за чего на
// локализованных системах (ru, de, fr, ...) сайдбар FileDialog теряет
// пункты Downloads/Documents/Music/Pictures/Videos/Desktop, а в лог
// сыплется "this computer does not define a X folder".
func getFavoriteLocation(homeURI fyne.URI, favName string) (fyne.URI, error) {
	if localName := xdgUserDir(favName); localName != "" {
		if uri, err := storage.Child(homeURI, localName); err == nil {
			if exists, _ := storage.Exists(uri); exists {
				return uri, nil
			}
		}
	}

	if uri, err := storage.Child(homeURI, favName); err == nil {
		if exists, _ := storage.Exists(uri); exists {
			return uri, nil
		}
	}

	return nil, fmt.Errorf("this computer does not define a %s folder", strings.ToUpper(favName))
}

// xdgUserDir читает ~/.config/user-dirs.dirs и возвращает локальное имя
// папки для запрошенного английского ключа ("Downloads" → "Загрузки").
// Возвращает "" на не-Linux, если файла нет, или если переменная не
// определена.
func xdgUserDir(favName string) string {
	if runtime.GOOS != "linux" && runtime.GOOS != "freebsd" &&
		runtime.GOOS != "openbsd" && runtime.GOOS != "netbsd" {
		return ""
	}

	key := ""
	switch favName {
	case folderDesktop:
		key = "XDG_DESKTOP_DIR"
	case folderDocuments:
		key = "XDG_DOCUMENTS_DIR"
	case folderDownloads:
		key = "XDG_DOWNLOAD_DIR"
	case folderMusic:
		key = "XDG_MUSIC_DIR"
	case folderPictures:
		key = "XDG_PICTURES_DIR"
	case folderVideos:
		key = "XDG_VIDEOS_DIR"
	default:
		return ""
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(home, ".config")
	}

	data, err := os.ReadFile(filepath.Join(configDir, "user-dirs.dirs"))
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, key+"=") {
			continue
		}
		val := strings.TrimPrefix(line, key+"=")
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"`)
		val = strings.ReplaceAll(val, "$HOME", home)
		rel, err := filepath.Rel(home, val)
		if err != nil || strings.HasPrefix(rel, "..") {
			return ""
		}
		return rel
	}
	return ""
}
