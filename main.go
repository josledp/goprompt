package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/josledp/goprompt/config"
	"github.com/josledp/goprompt/prompt"

	"github.com/josledp/termcolor"
)

var logger *log.Logger

type formatFunc func(string, ...termcolor.Mode) string

func main() {
	var noColor bool
	var template string
	var customTemplate string
	var t string

	config, err := config.New(os.Getenv("HOME") + "/.config/goprompt/goprompt.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	currentTemplates := strings.Join(config.GetTemplates(), ",")
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

	var options map[string]interface{}
	if customTemplateSet {
		t = customTemplate
	} else {
		var ok bool
		t, ok = config.GetTemplate(template)
		if !ok {
			fmt.Fprintf(os.Stderr, "Template %s not found", template)
		}
		options = config.GetTemplateOptions(template)

	}
	pr := prompt.New(options)
	output := pr.Compile(t, !noColor)
	fmt.Println(output)

}
