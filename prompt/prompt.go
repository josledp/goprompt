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
	Get(format func(string, ...termcolor.Mode) string) string
}

//Compile processes the template and returns a prompt string
func Compile(predefinedTemplate, template string, color bool) string {
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
			processedChunk := chunk[1 : len(chunk)-1]
			use := false
			plugins := rePlugin.FindAllString(chunk, -1)
			for _, rawPlugin := range plugins {
				plugin := rawPlugin[1 : len(rawPlugin)-1]
				if p, ok := mPlugins[plugin]; ok {
					//TODO +options
					err := p.Load(defaultOptions[predefinedTemplate])
					if err != nil {
						log.Printf("Unable to load plugin %s", plugin)
						continue
					}
					output := p.Get(format)
					if output != "" {
						use = true
					}
					processedChunk = strings.Replace(processedChunk, rawPlugin, output, -1)
				} else {
					log.Printf("Plugin %s not found", plugin)
				}
			}
			if use {
				pluginsOutput <- []string{chunk, processedChunk}
			} else {
				pluginsOutput <- []string{chunk, ""}
			}
			pluginsWg.Done()

		}(chunk)
	}
	for rep := range pluginsOutput {
		prompt = strings.Replace(prompt, rep[0], rep[1], -1)
	}
	return prompt
}

//Templates is a map with predefined templates
var Templates map[string]string
var defaultOptions map[string]map[string]interface{}
var availablePlugins []Plugin

func init() {
	availablePlugins = []Plugin{
		&plugins.Aws{},
		&plugins.Git{},
		&plugins.LastCommand{},
		&plugins.Path{},
		&plugins.Python{},
		&plugins.User{},
		&plugins.UserChar{},
		&plugins.Golang{},
	}

	Templates = map[string]string{
		"Evermeet": "<(%python%) ><%aws%|><%user% ><%lastcommand% ><%path%>< %git%><%userchar%> ",
		"Fedora":   "[ <(%python%) ><%aws%|><%user% ><%lastcommand% ><%path%>< %git%> ]<%userchar%> ",
	}
	defaultOptions = map[string]map[string]interface{}{
		"Evermeet": map[string]interface{}{
			"path.fullpath": true,
		},
		"Fedora": map[string]interface{}{
			"path.fullpath": false,
		},
	}
}
