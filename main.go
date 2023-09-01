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
	flag.Parse()

	if len(flag.Args()) > 1 {
		// TODO: print usage
		log.Fatal("Error: too many command-line arguments")
	}

	rootDir := "."
	if len(flag.Args()) == 1 {
		rootDir = flag.Args()[0]
	}

	addr := *hostFlag + ":" + *portFlag
	handler := http.FileServer(http.Dir(rootDir))

	log.Printf("Serving directory %q on http://%v", rootDir, addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
