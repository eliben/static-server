// Simple static file server.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"fmt"
	"log"
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
	hostFlag := flag.String("host", "", "specific host to listen on")
	portFlag := flag.String("port", "8080", "port to listen on")
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

	log.Printf("Serving directory %q on http://%v", rootDir, addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
