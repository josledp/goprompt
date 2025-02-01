package plugin

import (
	"fmt"
	"os"
	"strings"

	"github.com/josledp/termcolor"
)

const maxPathLength = 20

// Path is the plugin struct
type Path struct {
	pwd string
}

// Name returns the plugin name
func (Path) Name() string {
	return "path"
}

// Help returns help information about this plugin
func (Path) Help() (description string, options map[string]string) {
	description = "This plugins show the current path"
	options = map[string]string{
		"path.fullpath": "if 0 shows just the current dir, 1 for standard full path, 2 for fish path, 3 for variable path (it tries not to shrink it until it is >20)",
	}
	return
}

// Load is the load function of the plugin
func (p *Path) Load(pr Prompter) error {
	p.pwd = os.Getenv("PWD")
	if p.pwd == "" {
		return fmt.Errorf("unable to get PWD")
	}

	home := os.Getenv("HOME")
	if home != "" {
		p.pwd = strings.Replace(p.pwd, home, "~", -1)
	}

	if pr != nil {
		if value, ok := pr.GetOption(p.Name() + ".fullpath"); ok {
			if v, ok := value.(float64); ok {
				switch v {
				case 0:
					tmp := strings.Split(p.pwd, "/")
					p.pwd = tmp[len(tmp)-1]
				case 2:
					tmp := strings.Split(p.pwd, "/")
					for i, d := range tmp {
						if i == 0 {
							p.pwd = d
						} else if i < len(tmp)-1 {
							p.pwd += "/" + string(d[0])
						} else {
							p.pwd += "/" + d
						}
					}
				case 3:
					length := len(p.pwd)
					tmp := strings.Split(p.pwd, "/")
					var i int
					p.pwd = ""
					for i = 0; i < len(tmp) && length > maxPathLength; i++ {
						d := tmp[i]
						if i == 0 {
							p.pwd = d
						} else if i < len(tmp)-1 {
							p.pwd += "/" + string(d[0])
							length -= len(d) - 1
						} else {
							p.pwd += "/" + d
						}
					}
					for j := i; j < len(tmp); j++ {
						if j == 0 {
							p.pwd = tmp[0]
						} else {
							p.pwd += "/" + tmp[j]
						}
					}
				}
			} else if v, ok := value.(bool); ok {
				if !v {
					tmp := strings.Split(p.pwd, "/")
					p.pwd = tmp[len(tmp)-1]
				}
			} else {
				return fmt.Errorf("unable to parse path.fullpath option")
			}

		}
	}

	return nil
}

// Get returns the string to use in the prompt
func (p Path) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	return format(p.pwd, termcolor.Bold, termcolor.FgBlue), []termcolor.Mode{termcolor.FgBlue}
}
