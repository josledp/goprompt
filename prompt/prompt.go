package prompt

import (
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/josledp/goprompt/plugins"
	"github.com/josledp/termcolor"
)

//Plugin is the interface all the plugins MUST implement
type Plugin interface {
	Name() string
	Load(options map[string]interface{}) error
	Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode)
}

//Compile processes the template and returns a prompt string
func Compile(template string, color bool, options map[string]interface{}) string {
	var format func(string, ...termcolor.Mode) string
	prompt := template

	if color {
		format = termcolor.EscapedFormat
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
				err := p.Load(options)
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
		prompt = strings.Replace(prompt, rep[0], rep[1], -1)
	}
	return prompt
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
