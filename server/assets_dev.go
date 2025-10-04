//go:build !prod

package main

import "io/fs"

// No embedded files in dev; Vite serves the UI.
// This keeps `go run .` working even if web/dist doesn't exist.
var embeddedFS fs.FS = nil
