// Simple static file server.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"log"
	"net/http"
)

func main() {
	// TODO: make all of these configurable
	domain := "foo.local"
	port := "8080"
	addr := domain + ":" + port
	dir := "."
	handler := http.FileServer(http.Dir("."))

	log.Printf("Serving directory %q on http://%v", dir, addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
