package plugin

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

//ExitUserChar is the plugin struct
type ExitUserChar struct {
	user   string
	lastrc string
}

//Name returns the plugin name
func (ExitUserChar) Name() string {
	return "exituserchar"
}

//Help returns help information about this plugin
func (ExitUserChar) Help() (description string, options map[string]string) {
	description = "This plugins show the typical final char for the prompt (# is the user is root, $ otherwise) but it will be red if the last command exited with rc!=0"
	return
}

//Load is the load function of the plugin
func (euc *ExitUserChar) Load(Prompter) error {
	euc.user = os.Getenv("USER")
	if euc.user == "" {
		return fmt.Errorf("Unable to get USER")
	}
	euc.lastrc = os.Getenv("LAST_COMMAND_RC")
	if euc.lastrc == "" {
		return fmt.Errorf("Unable to get LAST_COMMAND_RC")
	}
	return nil

}

//Get returns the string to use in the prompt
func (euc *ExitUserChar) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	char := "$"
	if euc.user == "root" {
		char = "#"
	}
	if euc.lastrc == "0" {
		return char, nil
	}
	return format(char, termcolor.FgHiRed), []termcolor.Mode{termcolor.FgHiRed}
}
