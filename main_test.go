package main_test

import (
	"net"
	"net/http"
	"os"
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
				addr := ts.Getenv("ADDR")
				path := "http://" + addr + "/__internal/__shutdown"
				resp, err := http.Get(path)
				if err == nil {
					resp.Body.Close()
				}
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
