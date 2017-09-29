package plugins

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

//LastCommand is the plugin struct
type LastCommand struct {
	lastrc string
}

//Name returns the plugin name
func (LastCommand) Name() string {
	return "lastcommand"
}

//Load is the load function of the plugin
func (lc *LastCommand) Load(options map[string]interface{}) error {

	lc.lastrc = os.Getenv("LAST_COMMAND_RC")
	if lc.lastrc == "" {
		return fmt.Errorf("Unable to get LAST_COMMAND_RC")
	}

	return nil

}

//Get returns the string to use in the prompt
func (lc LastCommand) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if lc.lastrc != "" {
		return format(lc.lastrc, termcolor.FgHiYellow), []termcolor.Mode{termcolor.FgHiYellow}
	}
	return "", nil
}
