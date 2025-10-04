//go:build prod

package main

import (
	"embed"
)

// Embed the built UI: make sure you've run `npm run build` in /web first.
//go:embed all:web/dist/*
var embeddedFS embed.FS
