package static

import "embed"

//go:embed js css images favicon.ico
var StaticFS embed.FS
