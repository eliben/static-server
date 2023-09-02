// Simple static file server.
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
