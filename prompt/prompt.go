package prompt

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/josledp/goprompt/prompt/plugin"
	"github.com/josledp/termcolor"
)

var availablePlugins = []Plugin{
	&plugin.Aws{},
	&plugin.Git{},
	&plugin.LastCommand{},
	&plugin.Path{},
	&plugin.Python{},
	&plugin.User{},
	&plugin.Hostname{},
	&plugin.UserChar{},
	&plugin.Golang{},
	&plugin.Kubernetes{},
	&plugin.ExitUserChar{},
}

var defaultTemplates = map[string]string{
	"Evermeet": `{{load "python" |suffix " "}}{{load "aws"|suffix "|"}} {{load "user"|suffix "@"}}{{load "hostname"}} {{load "lastcommand"|suffix " "}}{{load "path"}}{{load "git"|prefix " "}}{{load "userchar"}}`,
	"Fedora":   `[ {{load "python"|wrap "(" ") "}}{{load "aws"|suffix "|"}}{{load "user"|suffix "@"}}{{load "hostname"}} {{load "lastcommand"|suffix " "}}{{load "path"}}{{load "git"|prefix " "}} ]{{load "userchar"}} `,
	"Prefered": `{{load "k8s"}}{{load "python"|wrap "("  ") "}}{{load "aws"|replace "(.*):.*-([^-]*)$" "$1:$2"|suffix "|"}}{{load "path"}}{{load "git"|prefix " "}}{{load "exituserchar"}} `,
}

var defaultTemplatesOptions = map[string]map[string]interface{}{
	"Evermeet": map[string]interface{}{
		"path.fullpath": float64(1),
	},
	"Fedora": map[string]interface{}{
		"path.fullpath": float64(0),
	},
	"Prefered": map[string]interface{}{
		"path.fullpath": float64(3),
	},
}

//Prompt is the struct with the prompt options/config
type Prompt struct {
	options map[string]interface{}
	cache   *Cache
	plugins map[string]Plugin
	format  func(string, ...termcolor.Mode) string

	debug   bool
	tmpMode []termcolor.Mode
}

//Plugin is the interface all the plugins MUST implement
type Plugin interface {
	Name() string
	Help() (description string, options map[string]string)
	Load(pr plugin.Prompter) error
	Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode)
}

//New returns a new promp
func New(options map[string]interface{}, color, debug bool) Prompt {
	c, err := newCache()
	if err != nil {
		log.Printf("unable to initializa cache: %v", err)
	}
	// map plugin by name
	mPlugins := make(map[string]Plugin)
	for _, p := range availablePlugins {
		mPlugins[p.Name()] = p
	}
	var format func(string, ...termcolor.Mode) string

	if color {
		shell := detectShell()
		switch shell {
		case "bash":
			format = termcolor.EscapedFormat
		case "fish":
			format = termcolor.Format
		case "zsh":
			format = termcolor.Format

		default:
			//Defaut failsafe
			format = func(s string, modes ...termcolor.Mode) string { return s }
		}
	} else {
		format = func(s string, modes ...termcolor.Mode) string { return s }
	}

	return Prompt{
		options: options,
		cache:   c,
		plugins: mPlugins,
		format:  format,
		debug:   debug,
		tmpMode: nil,
	}
}

//GetOption returns the option value for key
func (pr Prompt) GetOption(key string) (interface{}, bool) {
	value, ok := pr.options[key]
	return value, ok
}

//GetCache recovers a value from cache
func (pr Prompt) GetCache(key string) (interface{}, bool) {
	return pr.cache.get(key)
}

//Cache caches a key, value on cache
func (pr Prompt) Cache(key string, value interface{}) error {
	return pr.cache.set(key, value)
}

//Compile processes the template and returns a prompt string
func (pr *Prompt) Compile(tmpl string) string {

	t, err := template.New("prompt").Funcs(pr.getFuncMap()).Parse(tmpl)
	if err != nil {
		log.Fatalf("unable to parse tmpl %s: %v", tmpl, err)
	}

	b := &bytes.Buffer{}
	err = t.Execute(b, struct{}{})
	if err != nil {
		log.Fatalf("unable to execute tmpl %s: %v", tmpl, err)
	}
	err = pr.cache.save()
	if err != nil {
		log.Printf("Unable to save cache: %v", err)
	}
	return b.String()
}

func (pr *Prompt) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"load":    pr.Load,
		"wrap":    pr.Wrap,
		"suffix":  pr.Suffix,
		"prefix":  pr.Prefix,
		"replace": pr.Replace,
	}
}

func (pr *Prompt) Load(plugin string) (string, error) {
	var p Plugin
	var ok bool
	var output string

	if p, ok = pr.plugins[plugin]; !ok {
		return "", fmt.Errorf("unable to find plugin: %s", plugin)
	}
	err := p.Load(pr)
	if err != nil {
		return "", fmt.Errorf("unable to load plugin %s: %v", plugin, err)
	}
	output, pr.tmpMode = p.Get(pr.format)

	if pr.debug {
		fmt.Fprintf(os.Stderr, "plugin %s output: %s\n", plugin, output)
	}
	return output, nil
}

func (pr *Prompt) Wrap(prefix, suffix string, input string) string {
	if input == "" {
		return ""
	}
	if prefix != "" {
		prefix = pr.format(prefix, pr.tmpMode...)
	}
	if suffix != "" {
		suffix = pr.format(suffix, pr.tmpMode...)
	}
	return fmt.Sprintf("%s%s%s", prefix, input, suffix)
}
func (pr *Prompt) Prefix(prefix, input string) string {
	return pr.Wrap(prefix, "", input)
}
func (pr *Prompt) Suffix(suffix, input string) string {
	return pr.Wrap("", suffix, input)
}
func (pr *Prompt) Replace(expr, repl, src string) (string, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(src, repl), nil
}
