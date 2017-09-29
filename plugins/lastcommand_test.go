package plugins

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestLastCommand(t *testing.T) {

	os.Setenv("LAST_COMMAND_RC", "10")
	expected := "\\[\\033[0m\\]\\[\\033[93m\\]10\\[\\033[0m\\]"

	lc := &LastCommand{}
	lc.Load(map[string]interface{}{})

	if lc.lastrc != "10" {
		t.Error("Invalid Last command rc")
	}

	output := lc.Get(termcolor.EscapedFormat)
	if output != expected {
		t.Errorf("Expected %s\nGot      %s", expected, output)
	}

}

func TestLastCommandError(t *testing.T) {

	os.Setenv("LAST_COMMAND_RC", "")

	lc := &LastCommand{}
	err := lc.Load(map[string]interface{}{})

	if err.Error() != "Unable to get LAST_COMMAND_RC" {
		t.Errorf("Invalid Last command Error: %v", err)
	}

}
