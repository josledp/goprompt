package plugins

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
func (u *User) Load(options map[string]interface{}) error {
	var err error
	u.user = os.Getenv("USER")
	if u.user == "" {
		return fmt.Errorf("Unable to get USER")
	}
	u.hostname, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("Unable to get Hostname: %v", err)
	}

	return nil

}

//Get returns the string to use in the prompt
func (u User) Get(format func(string, ...termcolor.Mode) string) string {
	var userPromptInfo string
	if u.user == "root" {
		userPromptInfo = format(u.hostname, termcolor.Bold, termcolor.FgRed)
	} else {
		userPromptInfo = format(u.user+"@"+u.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	return userPromptInfo
}
