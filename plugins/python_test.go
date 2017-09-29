package plugins

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestPython(t *testing.T) {
	os.Setenv("VIRTUAL_ENV", "./Envs/env")
	expected := "\\[\\033[0m\\]\\[\\033[34m\\]env\\[\\033[0m\\]"

	p := &Python{}
	p.Load(nil)

	if p.virtualEnv != "env" {
		t.Error("Invalid virtualenv")
	}

	output := p.Get(termcolor.EscapedFormat)
	if output != expected {
		t.Errorf("Expected %s\nGot      %s", expected, output)
	}

}
