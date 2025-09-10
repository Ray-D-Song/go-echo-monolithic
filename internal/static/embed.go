package static

import (
	"embed"
	"io/fs"
)

// Embed web directory into binary
//go:embed all:web
var WebFS embed.FS

// GetWebFS returns the embedded web filesystem
func GetWebFS() fs.FS {
	webFS, err := fs.Sub(WebFS, "web")
	if err != nil {
		panic(err)
	}
	return webFS
}