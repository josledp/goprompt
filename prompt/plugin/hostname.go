package plugin

import (
	"fmt"
	"os"
	"strings"

	"github.com/josledp/termcolor"
)

// Hostname is the plugin struct
type Hostname struct {
	hostname string
	user     string
}

// Name returns the plugin name
func (Hostname) Name() string {
	return "hostname"
}

// Help returns help information about this plugin
func (Hostname) Help() (description string, options map[string]string) {
	description = "This plugins shows the current hostname (red if you are root, green otherwise)"
	return
}

// Load is the load function of the plugin
func (h *Hostname) Load(Prompter) error {
	var err error
	h.user = os.Getenv("USER")

	h.hostname, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get Hostname: %v", err)
	}

	h.hostname = strings.Split(h.hostname, ".")[0]
	return nil
}

// Get returns the string to use in the prompt
func (h Hostname) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if h.user == "root" {
		return format(h.hostname, termcolor.Bold, termcolor.FgRed), []termcolor.Mode{termcolor.FgRed}
	}
	return format(h.hostname, termcolor.Bold, termcolor.FgGreen), []termcolor.Mode{termcolor.FgGreen}
}
