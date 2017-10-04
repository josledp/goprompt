package prompt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/josledp/goprompt/plugins"
	"github.com/josledp/termcolor"
)

//Prompt is the struct with the prompt options/config
type Prompt struct {
	options map[string]interface{}
}

//Plugin is the interface all the plugins MUST implement
type Plugin interface {
	Name() string
	Load(pr plugins.Prompter) error
	Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode)
}

//New returns a new prompt
func New(options map[string]interface{}) Prompt {
	return Prompt{options}
}

//GetOption returns the option value for key
func (pr Prompt) GetOption(key string) (interface{}, bool) {
	value, ok := pr.options[key]
	return value, ok
}

//GetConfig return the config value for key
func (pr Prompt) GetConfig(key string) (interface{}, bool) {
	value, ok := 0, false
	return value, ok
}

//SetConfig sets a config value
func (pr Prompt) SetConfig(key string, value interface{}) error {
	return nil
}

//Compile processes the template and returns a prompt string
func (pr Prompt) Compile(template string, color bool) string {
	var format func(string, ...termcolor.Mode) string
	output := template

	if color {

		shell := pr.detectShell()
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

	// map plugin by name
	mPlugins := make(map[string]Plugin)
	for _, p := range availablePlugins {
		mPlugins[p.Name()] = p
	}

	//Regular expresions for matching on template
	reChunk, _ := regexp.Compile("<[^<>]*>")
	rePlugin, _ := regexp.Compile("%[a-z]*%")
	chunks := reChunk.FindAllString(template, -1)

	//Channel for plugins to write (parallel plugin processing)
	pluginsOutput := make(chan []string)
	pluginsWg := sync.WaitGroup{}
	pluginsWg.Add(len(chunks))
	go func() {
		pluginsWg.Wait()
		close(pluginsOutput)
	}()

	//For each chunk of <[^<>]*> we process it in parallel putting in the channel the string to replace on template
	for _, chunk := range chunks {
		go func(chunk string) {
			defer pluginsWg.Done()
			processedChunk := chunk[1 : len(chunk)-1]
			rawPlugin := rePlugin.FindString(chunk)
			plugin := rawPlugin[1 : len(rawPlugin)-1]
			if p, ok := mPlugins[plugin]; ok {
				//TODO +options
				err := p.Load(pr)
				if err != nil {
					log.Printf("Unable to load plugin %s", plugin)
					return
				}
				output, modes := p.Get(format)
				if output != "" {
					extra := strings.Split(processedChunk, rawPlugin)
					for _, e := range extra {
						useless, _ := regexp.MatchString("^[ ]*$", e)
						if !useless {
							processed := format(e, modes...)
							processedChunk = strings.Replace(processedChunk, e, processed, -1)
						}
					}
					processedChunk = strings.Replace(processedChunk, rawPlugin, output, -1)
					pluginsOutput <- []string{chunk, format(processedChunk, modes...)}
				} else {
					pluginsOutput <- []string{chunk, ""}
				}
			} else {
				log.Printf("Plugin %s not found", plugin)
			}

		}(chunk)
	}
	for rep := range pluginsOutput {
		output = strings.Replace(output, rep[0], rep[1], -1)
	}
	return output
}

var availablePlugins []Plugin

func init() {
	availablePlugins = []Plugin{
		&plugins.Aws{},
		&plugins.Git{},
		&plugins.LastCommand{},
		&plugins.Path{},
		&plugins.Python{},
		&plugins.User{},
		&plugins.Hostname{},
		&plugins.UserChar{},
		&plugins.Golang{},
	}

}

func (p Prompt) detectShell() string {
	pid := os.Getppid()
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdline, err := ioutil.ReadFile(cmdlineFile)
	if err != nil {
		return "unknown"
	}

	shells := []string{"bash", "zsh", "fish"}
	for _, shell := range shells {
		if matches, _ := regexp.Match(shell, cmdline); matches {
			return shell
		}
	}
	return "unknown"
}
