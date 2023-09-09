// Tests the server via its command-line interface using the
// testscript package.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main_test

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/eliben/static-server/internal/server"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"server": server.Main,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:      "testdata/scripts",
		TestWork: true,
		Setup: func(env *testscript.Env) error {
			// Make all the files from testdata/datafiles available for tests in
			// their datafiles/ directory.
			rootdir, err := os.Getwd()
			check(t, err)
			copyDataFiles(t,
				filepath.Join(rootdir, "testdata", "datafiles"),
				filepath.Join(env.WorkDir, "datafiles"))

			// Generate a fresh address for every test script, to avoid collisions
			// between multiple tests running in parallel.
			addr := randomLocalAddr(t)
			env.Setenv("ADDR", addr)
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"shutdown": func(ts *testscript.TestScript, neg bool, args []string) {
				// Custom command that connects to the "please shutdown"
				// endpoint on the server for graceful shutdown. After this command,
				// the server will exit.
				shutdownServer(ts.Getenv("ADDR"))
			},

			"shutdown_tls": func(ts *testscript.TestScript, neg bool, args []string) {
				certfile := filepath.Join(ts.Getenv("WORK"), "datafiles", "cert.pem")
				shutdownServerTLS(ts.Getenv("ADDR"), certfile)
			},
		},
	})
}

// randomLocalAddr finds a random free port
func randomLocalAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

// copyDataFiles copies all files from rootdir to targetdir, creating
// targetdir if needed
func copyDataFiles(t *testing.T, rootdir string, targetdir string) {
	check(t, os.MkdirAll(targetdir, 0777))

	entries, err := os.ReadDir(rootdir)
	check(t, err)
	for _, e := range entries {
		if !e.IsDir() {
			fullpath := filepath.Join(rootdir, e.Name())
			targetpath := filepath.Join(targetdir, e.Name())

			data, err := os.ReadFile(fullpath)
			check(t, err)
			err = os.WriteFile(targetpath, data, 0666)
			check(t, err)
		}
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func shutdownServer(addr string) {
	path := "http://" + addr + "/__internal/__shutdown"
	resp, err := http.Get(path)
	if err == nil {
		resp.Body.Close()
	}
}

func shutdownServerTLS(addr string, certpath string) {
	path := "https://" + addr + "/__internal/__shutdown"
	cert, err := os.ReadFile(certpath)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", certpath)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	resp, err := client.Get(path)
	if err == nil {
		resp.Body.Close()
	}
}
