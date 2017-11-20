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

//type formatFunc func(string, ...termcolor.Mode) string

func enableCpuProf() func() {

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
	defer enableCpuProf()()

	config, err := prompt.NewConfigFromFile(os.Getenv("HOME") + "/.config/goprompt/goprompt.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	currentTemplates := strings.Join(prompt.GetDefaultTemplates(), ",")
	flag.StringVar(&template, "template", "Evermeet", "template to use for the prompt ("+currentTemplates+")")
	flag.StringVar(&customTemplate, "custom-template", "<(%python%) ><%aws%|><%user% ><%lastcommand% ><%path%>< %git%>$ ", "template to use for the prompt")
	flag.BoolVar(&noColor, "no-color", false, "Disable color on prompt")

	flag.Parse()

	flagsSet := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { flagsSet[f.Name] = struct{}{} })

	_, templateSet := flagsSet["template"]
	_, customTemplateSet := flagsSet["custom-template"]

	if templateSet && customTemplateSet {
		fmt.Fprintf(os.Stderr, "Please provice --template or --custom-template, but not both!")
		os.Exit(1)
	}

	var t string
	var options map[string]interface{}
	options, _ = config.GetOptions()

	if customTemplateSet {
		t = customTemplate
	} else if !templateSet {
		t, _ = config.GetCustomTemplate()
	}

	if t == "" {
		var ok bool
		t, ok = prompt.GetTemplate(template)
		if !ok {
			fmt.Fprintf(os.Stderr, "Template %s not found", template)
		}
		if options == nil {
			options, _ = prompt.GetTemplateOptions(template)
		}
	}
	pr := prompt.New(options)
	output := pr.Compile(t, !noColor)
	fmt.Println(output)

}
