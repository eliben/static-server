package main_test

import (
	"net"
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
	addr := randomLocalAddr(t)
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
		Setup: func(env *testscript.Env) error {
			env.Setenv("ADDR", addr)
			return nil
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
