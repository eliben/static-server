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
	"strings"
)

// TODO: add the opener feature (-o)

func Main() int {
	errorLog := log.New(os.Stderr, "", log.LstdFlags)
	serveLog := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		out := flags.Output()
		fmt.Fprintf(out, "Usage: %v [dir]\n", os.Args[0])
		fmt.Fprint(out, "\n  [dir] is optional; if not passed, '.' is used\n\n")
		flags.PrintDefaults()
	}

	hostFlag := flags.String("host", "localhost", "specific host to listen on")
	portFlag := flags.String("port", "8080", "port to listen on; if 0, a random available port will be used")
	addrFlag := flags.String("addr", "localhost:8080", "address to listen on; don't use this is 'port' or 'host' are set")
	flags.Parse(os.Args[1:])

	if len(flags.Args()) > 1 {
		errorLog.Println("Error: too many command-line arguments")
		flags.Usage()
		os.Exit(1)
	}

	rootDir := "."
	if len(flags.Args()) == 1 {
		rootDir = flags.Args()[0]
	}

	allSetFlags := flagsSet(flags)
	if allSetFlags["addr"] && (allSetFlags["host"] || allSetFlags["port"]) {
		errorLog.Println("Error: if -addr is set, -host and -port must remain unset")
		flags.Usage()
		os.Exit(1)
	}

	var addr string
	if allSetFlags["addr"] {
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
	fileHandler := serveLogger(serveLog, http.FileServer(http.Dir(rootDir)))
	mux.Handle("/", fileHandler)
	srv.Handler = mux

	// Use an explicit listener to access .Addr() when serving on port :0
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		errorLog.Println(err)
		return 1
	}
	log.Printf("Serving directory %q on http://%v", rootDir, listener.Addr())

	if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errorLog.Println("Error in Serve:", err)
		return 1
	} else {
		return 0
	}
}

// flagsSet returns a set of all the flags what were actually set on the
// command line.
func flagsSet(flags *flag.FlagSet) map[string]bool {
	s := make(map[string]bool)
	flags.Visit(func(f *flag.Flag) {
		s[f.Name] = true
	})
	return s
}

// serveLogger is a logging middleware for serving. It generates logs for
// requests sent to the server.
func serveLogger(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteHost, _, _ := strings.Cut(r.RemoteAddr, ":")
		logger.Printf("%v %v %v\n", remoteHost, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
