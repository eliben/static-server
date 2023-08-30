// Simple static file server.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	hostFlag := flag.String("host", "", "specific host to listen on")
	portFlag := flag.String("port", "8080", "port to listen on")

	// TODO: dir should be argv[1] too?
	// TODO: sanitize dir for safety
	dirFlag := flag.String("dir", ".", "root directory to serve")
	flag.Parse()

	addr := *hostFlag + ":" + *portFlag
	handler := http.FileServer(http.Dir(*dirFlag))

	log.Printf("Serving directory %q on http://%v", *dirFlag, addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
