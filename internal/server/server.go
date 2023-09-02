package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// TODO: use local, not global flag, to avoid interference with the test suite!

func Main() int {
	flag.Usage = usage
	hostFlag := flag.String("host", "localhost", "specific host to listen on")
	portFlag := flag.String("port", "8080", "port to listen on; if 0, a random available port will be used")
	addrFlag := flag.String("addr", "localhost:8080", "address to listen on; don't use this is 'port' or 'host' are set")

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

	setOfFlags := flagsSet()
	if setOfFlags["addr"] && (setOfFlags["host"] || setOfFlags["port"]) {
		log.Println("Error: if -addr is set, -host and -port must remain unset")
		usage()
		os.Exit(1)
	}

	var addr string
	if setOfFlags["addr"] {
		addr = *addrFlag
	} else {
		addr = *hostFlag + ":" + *portFlag
	}
	srv := &http.Server{
		Addr: addr,
	}

	// To shut the server down cleanly in tests, we register a special route
	// where we ask it to stop. A separate goroutine performs the shutdown so
	// that the server can properly answer the shutdown request without abruptly
	// closing the connection.
	shutdownCh := make(chan struct{})
	go func() {
		<-shutdownCh
		srv.Shutdown(context.Background())
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/__internal/__shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		defer close(shutdownCh)
	})
	fileHandler := http.FileServer(http.Dir(rootDir))
	mux.Handle("/", fileHandler)
	srv.Handler = mux

	// Use an explicit listener to access .Addr() when serving on port :0
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return 1
	}
	log.Printf("Serving directory %q on http://%v", rootDir, listener.Addr())

	if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error in Serve:", err)
		return 1
	} else {
		return 0
	}
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage: %v [dir]\n", os.Args[0])
	fmt.Fprint(out, "\n  [dir] is optional; if not passed, '.' is used\n\n")
	flag.PrintDefaults()
}

// flagsSet returns a set of all the flags what were actually set on the
// command line.
func flagsSet() map[string]bool {
	s := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		s[f.Name] = true
	})
	return s
}
