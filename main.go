// Simple static file server.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage: %v [dir]\n", os.Args[0])
	fmt.Fprint(out, "\n  [dir] is optional; if not passed, '.' is used\n\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	hostFlag := flag.String("host", "localhost", "specific host to listen on")
	portFlag := flag.String("port", "8080", "port to listen on; if 0, a random available port will be used")
	flag.Parse()

	if len(flag.Args()) > 1 {
		log.Println("Error: too many command-line arguments")
		usage()
		os.Exit(1)
	}

	rootDir := "."
	if len(flag.Args()) == 1 {
		rootDir = flag.Args()[0]
	}

	addr := *hostFlag + ":" + *portFlag
	handler := http.FileServer(http.Dir(rootDir))

	// Use an explicit listener to access .Addr() when serving on port :0
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Serving directory %q on http://%v", rootDir, listener.Addr())
	log.Fatal(http.Serve(listener, handler))
}
