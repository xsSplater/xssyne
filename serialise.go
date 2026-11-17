// xssyne/serialise.go

package fyne

import (
	"fmt"
	"strconv"
	"strings"
)

// GoString converts a Resource object to Go code.
// This is useful if serialising to a Go file for compilation into a binary.
// xssyne
func (r *StaticResource) GoString() string {
	var buffer strings.Builder

	// Имя пакета вынесите в константу, чтобы в одном месте поменять при ребрендинге.
	const pkgPrefix = "&xssyne.StaticResource{\n\tStaticName: "

	buffer.WriteString(pkgPrefix)
	buffer.WriteString(strconv.Quote(r.StaticName))
	buffer.WriteString(",\n\tStaticContent: []byte{\n\t\t")
	for i, v := range r.StaticContent {
		if i > 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprint(&buffer, v)
	}
	buffer.WriteString("}}")

	return buffer.String()
}
