// xssyne/app/cache.go

package app

import (
	"encoding/base64"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type cache struct {
	a fyne.App
}

// xssyne
func makeCache(a fyne.App) fyne.Cache {
	c := &cache{a: a}

	root := c.RootURI()

	// Толерантная проверка: если root уже существует — ничего не делаем.
	if exists, err := storage.Exists(root); err == nil && exists {
		return c
	}

	// Пытаемся создать parent + root. Ошибки логируем, но cache
	// возвращаем рабочим: индивидуальные Read/Write/Remove дадут
	// понятную ошибку, а не молча сломаются при первом обращении.
	if parent, err := storage.Parent(root); err == nil {
		if pExists, _ := storage.Exists(parent); !pExists {
			if err := storage.CreateListable(parent); err != nil {
				fyne.LogError("Failed to create xssyne cache parent directory", err)
			}
		}
	}

	if err := storage.CreateListable(root); err != nil {
		fyne.LogError("Failed to create app cache directory", err)
	}

	return c
}

// xssyne
func (c *cache) encodePath(badName string) (fyne.URI, error) {
	name := base64.RawURLEncoding.EncodeToString([]byte(badName))
	child, err := storage.Child(c.RootURI(), name)
	if err != nil {
		return nil, err
	}
	return child, nil
}

func (c *cache) RootURI() fyne.URI {
	return storage.NewFileURI(rootCacheDir(c.a))
}

func (c *cache) Exists(name string) bool {
	path, err := c.encodePath(name)
	if err != nil {
		return false
	}
	ok, err := storage.Exists(path)
	return ok && err == nil
}

func (c *cache) Read(name string) (io.ReadCloser, error) {
	path, err := c.encodePath(name)
	if err != nil {
		return nil, err
	}
	return storage.Reader(path)
}

func (c *cache) Write(name string) (io.WriteCloser, error) {
	path, err := c.encodePath(name)
	if err != nil {
		return nil, err
	}

	// Defensive: makeCache мог провалиться (нет прав на старте, или
	// директория была удалена). Пытаемся создать parent на месте —
	// если не получится, вернём ошибку вызывающему вместо тихого nil.
	parent, err := storage.Parent(path)
	if err != nil {
		return nil, err
	}
	if exists, _ := storage.Exists(parent); !exists {
		if err := storage.CreateListable(parent); err != nil {
			return nil, err
		}
	}

	return storage.Writer(path)
}

func (c *cache) Remove(name string) error {
	path, err := c.encodePath(name)
	if err != nil {
		return err
	}
	return storage.Delete(path)
}
