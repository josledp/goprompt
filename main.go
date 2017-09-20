package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	git2go "gopkg.in/libgit2/git2go.v26"

	"takiones.com/goprompt/colors"
)

func getPythonVirtualEnv() string {
	virtualEnv, ve := os.LookupEnv("VIRTUAL_ENV")
	if ve {
		ave := strings.Split(virtualEnv, "/")
		virtualEnv = fmt.Sprintf("(%s) ", ave[len(ave)-1])
	}
	return virtualEnv
}

func main() {
	//Get basicinfo
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current path", err)
	}
	user := os.Getenv("USER")
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname", err)
	}
	//Get Python VirtualEnv info
	virtualEnv := getPythonVirtualEnv()

	//Get git information
	_ = git2go.Repository{}

	//discover:= git2go.Discover(".",false, ["/","."] )

	//Format output
	userInfo := colors.Format(user+"@"+hostname, colors.Bold, colors.FgGreen)
	pwdInfo := colors.Format(pwd, colors.Bold, colors.FgBlue)
	virtualEnvInfo := colors.Format(virtualEnv, colors.FgBlue)

	fmt.Printf("%s[%s %s]$ ", virtualEnvInfo, userInfo, pwdInfo)
}
