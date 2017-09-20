package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	reset        = "0"
	bold         = "1"
	underline    = "4"
	boldOff      = "21"
	underlineOff = "24"
	black        = "30"
	red          = "31"
	green        = "32"
	yellow       = "33"
	blue         = "34"
	magenta      = "35"
	cyan         = "36"
	white        = "37"
	start        = "\033["
	end          = "m"
)

func resetCode() {

}
func getPythonVirtualEnv() string {
	virtualEnv, ve := os.LookupEnv("VIRTUAL_ENV")
	if ve {
		ave := strings.Split(virtualEnv, "/")
		virtualEnv = fmt.Sprintf("(%s) ", ave[len(ave)-1])
	}
	return virtualEnv
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current path", err)
	}
	user := os.Getenv("USER")
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname", err)
	}
	virtualEnv := getPythonVirtualEnv()

	//Does not work exitCode, _ := os.LookupEnv("PIPESTATUS")
	fmt.Printf("%s[%s@%s  %s]$ ", virtualEnv, user, hostname, pwd)
}
