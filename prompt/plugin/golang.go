package plugin

import (
	"runtime"

	"github.com/josledp/termcolor"
)

//Golang is the plugin struct
type Golang struct {
	version string
}

//Name returns the plugin name
func (Golang) Name() string {
	return "golang"
}

//Help returns help information about this plugin
func (Golang) Help() (description string, options map[string]string) {
	description = "This plugins show current golang version"
	return
}

//Load is the load function of the plugin
func (g *Golang) Load(Prompter) error {
	g.version = runtime.Version()
	return nil
}

//Get returns the string to use in the prompt
func (g Golang) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	return format(g.version, termcolor.FgBlue), []termcolor.Mode{termcolor.FgBlue}
}
