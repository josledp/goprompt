package plugins

import (
	"fmt"
	"os"
	"strings"

	"github.com/josledp/termcolor"
)

//Path is the plugin struct
type Path struct {
	pwd string
}

//Name returns the plugin name
func (Path) Name() string {
	return "path"
}

//Load is the load function of the plugin
func (p *Path) Load(pr Prompter) error {
	p.pwd = os.Getenv("PWD")
	if p.pwd == "" {
		return fmt.Errorf("Unable to get PWD")
	}

	home := os.Getenv("HOME")
	if home != "" {
		p.pwd = strings.Replace(p.pwd, home, "~", -1)
	}

	if pr != nil {
		if value, ok := pr.GetOption(p.Name() + ".fullpath"); ok {
			if !value.(bool) {
				tmp := strings.Split(p.pwd, "/")
				p.pwd = tmp[len(tmp)-1]
			}
		}
	}

	return nil
}

//Get returns the string to use in the prompt
func (p Path) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	return format(p.pwd, termcolor.Bold, termcolor.FgBlue), []termcolor.Mode{termcolor.FgBlue}
}
