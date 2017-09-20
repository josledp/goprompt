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

	gitpath, err := git2go.Discover(".", false, []string{"/"})
	if err == nil {
		repository, err := git2go.OpenRepository(gitpath)
		if err != nil {
			log.Fatalf("Error opening repository at %s: %v", gitpath, err)
		}
		defer repository.Free()
		repostate, err := repository.StatusList(nil)
		if err != nil {
			log.Fatalf("Impsible to get repository status at %s: %v", gitpath, err)
		}
		defer repostate.Free()
		n, err := repostate.EntryCount()
		for i := 0; i < n; i++ {
			entry, _ := repostate.ByIndex(i)
			fmt.Println(entry.Status)
		}
	}
	userInfo := colors.Format(user+"@"+hostname, colors.Bold, colors.FgGreen)
	pwdInfo := colors.Format(pwd, colors.Bold, colors.FgBlue)
	virtualEnvInfo := colors.Format(virtualEnv, colors.FgBlue)

	fmt.Printf("%s[%s %s]$ ", virtualEnvInfo, userInfo, pwdInfo)
}
