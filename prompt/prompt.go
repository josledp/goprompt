package prompt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/josledp/goprompt/prompt/plugins"
	"github.com/josledp/termcolor"
)

//Prompt is the struct with the prompt options/config
type Prompt struct {
	options map[string]interface{}
	cache   *Cache
}

//Plugin is the interface all the plugins MUST implement
type Plugin interface {
	Name() string
	Load(pr plugins.Prompter) error
	Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode)
}

//New returns a new promp
func New(options map[string]interface{}) Prompt {
	c, err := newCache()
	if err != nil {
		log.Printf("unable to initializa cache: %v", err)
	}
	return Prompt{options, c}
}

//GetOption returns the option value for key
func (pr Prompt) GetOption(key string) (interface{}, bool) {
	value, ok := pr.options[key]
	return value, ok
}

//GetCache recovers a value from cache
func (pr Prompt) GetCache(key string) (interface{}, bool) {
	//Encapsulate more cache?
	if pr.cache == nil {
		return nil, false
	}
	value, ok := pr.cache.data[key]
	return value, ok
}

//Cache caches a key, value on cache
func (pr Prompt) Cache(key string, value interface{}) error {
	if pr.cache == nil {
		return fmt.Errorf("Cache not initialized")
	}
	if pr.cache.data == nil {
		pr.cache.data = make(map[string]interface{})
	}
	pr.cache.data[key] = value
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
					log.Printf("Unable to load plugin %s: %v", plugin, err)
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
	err := pr.cache.save()
	if err != nil {
		log.Printf("Unable to save cache: %v", err)
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

func (pr Prompt) detectShell() string {
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
