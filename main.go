// Simple static file server.
//
// This is the Go command-line entry point; the actual implementation is
// split off to a separate package for testing purposes.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"os"

	"github.com/eliben/static-server/internal/server"
)

func main() {
	os.Exit(server.Main())
}
