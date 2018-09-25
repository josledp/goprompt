package prompt

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

//GetDefaultTemplates returns the default templates defined by the prompt package
func GetDefaultTemplates() []string {
	templates := make([]string, 0)
	for name := range defaultTemplates {
		templates = append(templates, name)
	}
	return templates
}

//GetTemplateOptions returns the default options for a default template
func GetTemplateOptions(template string) (map[string]interface{}, bool) {
	o, ok := defaultTemplatesOptions[template]
	return o, ok
}

//GetTemplate returns the a default template by its name
func GetTemplate(template string) (string, bool) {
	t, ok := defaultTemplates[template]
	return t, ok
}

//ShowHelpPlugin writes on w the plugin help
func ShowHelpPlugin(w io.Writer) {
	fmt.Fprintf(w, "Plugin help\n")
	fmt.Fprintf(w, "===============\n")
	for _, p := range availablePlugins {
		name := p.Name()
		desc, opt := p.Help()
		fmt.Fprintf(w, "Plugin: %s\n", name)
		fmt.Fprintf(w, "Description: %s\n", desc)
		if len(opt) > 0 {
			fmt.Fprintf(w, "Options:\n")
			for o, od := range opt {
				fmt.Fprintf(w, "  %s: %s\n", o, od)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

//ShowHelpTemplate writes on w the templating help
func ShowHelpTemplate(w io.Writer) {
	fmt.Fprintf(w, "Templating help\n")
	fmt.Fprintf(w, "===============\n")
	fmt.Fprintln(w,
		`This project uses gotemplate. There are 4 functions over what gotemplate can do:
		load "plugin": will load plugin
		prefix, suffix, wrap: will add text/symbols before, after or both to any plugin output if it has content`)
}

func detectShell() string {
	pid := os.Getppid()
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdline, err := ioutil.ReadFile(cmdlineFile)
	if err != nil {
		cmdline = []byte(os.Getenv("SHELL"))
		if len(cmdline) <= 2 {
			return "unknown"
		}
	}

	shells := []string{"bash", "zsh", "fish"}
	for _, shell := range shells {
		if matches, _ := regexp.Match(shell, cmdline); matches {
			return shell
		}
	}
	return "unknown"
}
