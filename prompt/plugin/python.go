package plugin

import (
	"os"
	"strings"

	"github.com/josledp/termcolor"
)

// Python is the plugin struct
type Python struct {
	virtualEnv string
}

// Name returns the plugin name
func (Python) Name() string {
	return "python"
}

// Help returns help information about this plugin
func (Python) Help() (description string, options map[string]string) {
	description = "This plugins show the current python virtual environment"
	return
}

// Load is the load function of the plugin
func (p *Python) Load(Prompter) error {
	virtualEnv, ve := os.LookupEnv("VIRTUAL_ENV")
	if ve {
		ave := strings.Split(virtualEnv, "/")
		p.virtualEnv = ave[len(ave)-1]
	}
	return nil
}

// Get returns the string to use in the prompt
func (p Python) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if p.virtualEnv != "" {
		return format(p.virtualEnv, termcolor.FgBlue), []termcolor.Mode{termcolor.FgBlue}
	}
	return "", nil
}
