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

func makeCache(a fyne.App) fyne.Cache {
	c := &cache{a: a}

	root := c.RootURI()
	exists, err := storage.Exists(root)
	if !exists || err != nil {
		parent, _ := storage.Parent(root)
		exists, err = storage.Exists(parent)
		if !exists || err != nil {
			err = storage.CreateListable(parent)
			if err != nil {
				fyne.LogError("Failed to create fyne cache space", err)
			}
		}

		err = storage.CreateListable(root)
		if err != nil {
			fyne.LogError("Failed to create app cache space", err)
		}
	}

	return c
}

func (c *cache) encodePath(badName string) fyne.URI {
	name := base64.RawURLEncoding.EncodeToString([]byte(badName))
	child, _ := storage.Child(c.a.Cache().RootURI(), name)
	return child
}

func (c *cache) RootURI() fyne.URI {
	return storage.NewFileURI(rootCacheDir(c.a))
}

func (c *cache) Exists(name string) bool {
	path := c.encodePath(name)

	ok, err := storage.Exists(path)
	return ok && err == nil
}

func (c *cache) Read(name string) (io.ReadCloser, error) {
	path := c.encodePath(name)

	return storage.Reader(path)
}

func (c *cache) Write(name string) (io.WriteCloser, error) {
	path := c.encodePath(name)

	return storage.Writer(path)
}

func (c *cache) Remove(name string) error {
	path := c.encodePath(name)

	return storage.Delete(path)
}
