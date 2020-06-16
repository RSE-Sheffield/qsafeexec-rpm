package gowdlerize

import "os"
import "strings"
import "testing"

func TestAbs(t *testing.T) {
	os.Setenv("LD_PRELOAD", "/usr/lib/somelib.so")
	os.Setenv("PERFECTLY_SAFE", "some_value")
	goodVarFound := false

	env := CleanEnv(os.Environ())
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		switch varName := pair[0]; varName {
		case "LD_PRELOAD",
			"MALLOC_TRACE":
			t.Errorf("Bad env var %s should have been stripped from env but wasn't", varName)
		case "PERFECTLY_SAFE":
			goodVarFound = true
		default:
			continue
		}
	}
	if !goodVarFound {
		t.Errorf("Good env var PERFECTLY_SAFE erroneously stripped from env")
	}
	os.Unsetenv("LD_PRELOAD")
	os.Unsetenv("PERFECTLY_SAFE")
}

// TODO: Check that strips bad locales but allows good ones
// TODO: also test cmd
