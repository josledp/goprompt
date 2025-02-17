package plugin

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

// LastCommand is the plugin struct
type LastCommand struct {
	lastrc string
}

// Name returns the plugin name
func (LastCommand) Name() string {
	return "lastcommand"
}

// Help returns help information about this plugin
func (LastCommand) Help() (description string, options map[string]string) {
	description = "This plugins show the last command return code"
	return
}

// Load is the load function of the plugin
func (lc *LastCommand) Load(Prompter) error {

	lc.lastrc = os.Getenv("LAST_COMMAND_RC")
	if lc.lastrc == "" {
		return fmt.Errorf("unable to get LAST_COMMAND_RC")
	}

	return nil

}

// Get returns the string to use in the prompt
func (lc LastCommand) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if lc.lastrc != "" {
		return format(lc.lastrc, termcolor.FgHiYellow), []termcolor.Mode{termcolor.FgHiYellow}
	}
	return "", nil
}
