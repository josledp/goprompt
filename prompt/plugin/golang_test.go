package plugin

import (
	"runtime"
	"testing"

	"github.com/josledp/termcolor"
)

func TestGolang(t *testing.T) {
	golang := runtime.Version()
	expected := "\\[\\033[0m\\]\\[\\033[34m\\]" + golang + "\\[\\033[0m\\]"

	g := &Golang{}
	g.Load(nil)

	if g.version != golang {
		t.Error("Invalid golang version")
	}

	output, _ := g.Get(termcolor.EscapedFormat)
	if output != expected {
		t.Errorf("Expected %s\nGot      %s", expected, output)
	}

}
