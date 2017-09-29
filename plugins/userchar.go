package plugins

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

//UserChar is the plugin struct
type UserChar struct {
	user string
}

//Name returns the plugin name
func (UserChar) Name() string {
	return "userchar"
}

//Load is the load function of the plugin
func (uc *UserChar) Load(options map[string]interface{}) error {
	uc.user = os.Getenv("USER")
	if uc.user == "" {
		return fmt.Errorf("Unable to get USER")
	}
	return nil

}

//Get returns the string to use in the prompt
func (uc *UserChar) Get(format func(string, ...termcolor.Mode) string) string {
	if uc.user == "root" {
		return "#"
	}
	return "$"
}