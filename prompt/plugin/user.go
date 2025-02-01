package plugin

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

// User is the plugin struct
type User struct {
	user string
}

// Name returns the plugin name
func (User) Name() string {
	return "user"
}

// Help returns help information about this plugin
func (User) Help() (description string, options map[string]string) {
	description = "This plugins show the current user if its not root"
	return
}

// Load is the load function of the plugin
func (u *User) Load(Prompter) error {
	u.user = os.Getenv("USER")
	if u.user == "" {
		return fmt.Errorf("unable to get USER")
	}
	return nil
}

// Get returns the string to use in the prompt
func (u User) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if u.user == "root" {
		return "", nil
	}
	return format(u.user, termcolor.Bold, termcolor.FgGreen), []termcolor.Mode{termcolor.FgGreen}
}
