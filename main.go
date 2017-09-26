package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/josledp/goprompt/prompt"
	"github.com/josledp/termcolor"
)

var logger *log.Logger

type formatFunc func(string, ...termcolor.Mode) string

func main() {
	var fullpath, noFullpath bool
	var color, noColor bool
	var pFullpath, pColor *bool
	var style string

	//pushing flag package to its limits :)
	flag.StringVar(&style, "style", "Evermeet", "Select style: Evermeet, Mac, Fedora")
	flag.BoolVar(&fullpath, "fullpath", true, "Show fullpath on prompt. The default value depends on the style")
	flag.BoolVar(&noFullpath, "no-fullpath", false, "Show fullpath on prompt. The default value depends on the style")
	flag.BoolVar(&color, "color", true, "Show color on prompt. The default value depends on the style")
	flag.BoolVar(&noColor, "no-color", false, "Show color on prompt. The default value depends on the style")
	flag.Parse()

	flagsSet := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { flagsSet[f.Name] = struct{}{} })

	_, colorSet := flagsSet["color"]
	_, noColorSet := flagsSet["no-color"]
	_, fullpathSet := flagsSet["fullpath"]
	_, noFullpathSet := flagsSet["no-fullpath"]

	if colorSet && noColorSet {
		fmt.Fprintf(os.Stderr, "Please use --color or --no-color, not both")
		os.Exit(1)
	}

	if fullpathSet && noFullpathSet {
		fmt.Fprintf(os.Stderr, "Please use --fullpath or --no-fullpath, not both")
		os.Exit(1)
	}

	if colorSet {
		pColor = &color
	} else if noColorSet {
		pColor = &noColor
		*pColor = !noColor
	}

	if fullpathSet {
		pFullpath = &fullpath
	} else if noFullpathSet {
		pFullpath = &noFullpath
		*pFullpath = !noFullpath
	}

	pr := prompt.New(style, pColor, pFullpath)
	fmt.Println(pr.GetPrompt())

}
