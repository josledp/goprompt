package plugin

import (
	"fmt"
	"os"

	"github.com/josledp/termcolor"
)

//User is the plugin struct
type User struct {
	user     string
	hostname string
}

//Name returns the plugin name
func (User) Name() string {
	return "user"
}

//Load is the load function of the plugin
func (u *User) Load(Prompter) error {
	u.user = os.Getenv("USER")
	if u.user == "" {
		return fmt.Errorf("Unable to get USER")
	}
	return nil
}

//Get returns the string to use in the prompt
func (u User) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if u.user == "root" {
		return "", nil
	}
	return format(u.user, termcolor.Bold, termcolor.FgGreen), []termcolor.Mode{termcolor.FgGreen}
}
