package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/josledp/goprompt/prompt"
)

var logger *log.Logger

//enableCPUProf call me this way if you want CPU Profiling:
// defer enableCPUProf()()
func enableCPUProf() func() {

	f, err := os.Create("/tmp/cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

func main() {
	var noColor bool
	var template string
	var customTemplate string
	var helpPlugin, helpTemplate bool
	var debug bool

	config, err := prompt.NewConfigFromFile(os.Getenv("HOME") + "/.config/goprompt/goprompt.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	defaultTemplate, ok := config.GetTemplate()
	if !ok {
		defaultTemplate = "Evermeet"
	}
	currentTemplates := strings.Join(prompt.GetDefaultTemplates(), ",")
	flag.StringVar(&template, "template", defaultTemplate, "template to use for the prompt ("+currentTemplates+")")
	flag.StringVar(&customTemplate, "custom-template", "<(%python%) ><%aws%|><%user% ><%lastcommand% ><%path%>< %git%>$ ", "template to use for the prompt")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.BoolVar(&noColor, "no-color", false, "Disable color on prompt")
	flag.BoolVar(&helpPlugin, "help-plugin", false, "Shows plugins help")
	flag.BoolVar(&helpTemplate, "help-template", false, "Shows templating help")

	flag.Parse()

	if helpPlugin {
		prompt.ShowHelpPlugin(os.Stdout)
		os.Exit(0)
	}
	if helpTemplate {
		prompt.ShowHelpTemplate(os.Stdout)
		os.Exit(0)
	}

	flagsSet := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { flagsSet[f.Name] = struct{}{} })

	_, templateSet := flagsSet["template"]
	_, customTemplateSet := flagsSet["custom-template"]

	if templateSet && customTemplateSet {
		fmt.Fprintf(os.Stderr, "please provice -template or -custom-template, but not both!")
		os.Exit(1)
	}

	var t string
	var options map[string]interface{}
	options, _ = config.GetOptions()

	//If we provide a customTemplate in the command line use it. Otherwise, if template parameter is not set try to load the template from the config
	if customTemplateSet {
		t = customTemplate
	} else if !templateSet {
		t, _ = config.GetCustomTemplate()
	}

	//If we have not a template yet, get it (from the template parameter, or from the template option in the config)
	if t == "" {
		var ok bool
		t, ok = prompt.GetTemplate(template)
		if !ok {
			fmt.Fprintf(os.Stderr, "template %s not found", template)
		}
		if options == nil {
			options, _ = prompt.GetTemplateOptions(template)
		}
	}
	pr := prompt.New(options, !noColor, debug)
	output := pr.Compile(t)
	fmt.Println(output)

}
