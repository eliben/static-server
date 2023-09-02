package main_test

import (
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

func TestScriptWithExtraEnvVars(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
