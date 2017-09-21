package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	git2go "gopkg.in/libgit2/git2go.v24"

	"github.com/josledp/termcolor"
)

func getPythonVirtualEnv() string {
	virtualEnv, ve := os.LookupEnv("VIRTUAL_ENV")
	if ve {
		ave := strings.Split(virtualEnv, "/")
		virtualEnv = fmt.Sprintf("(%s) ", ave[len(ave)-1])
	}
	return virtualEnv
}

func getAwsInfo() string {
	role := os.Getenv("AWS_ROLE")
	if role != "" {
		tmp := strings.Split(role, ":")
		role = tmp[0]
		tmp = strings.Split(tmp[1], "-")
		role += ":" + tmp[2]
	}
	return role
}

func main() {
	//Get basicinfo
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current path", err)
	}
	home := os.Getenv("HOME")
	if home != "" {
		pwd = strings.Replace(pwd, home, "~", -1)
	}
	user := os.Getenv("USER")
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname", err)
	}

	//Get Python VirtualEnv info
	virtualEnv := getPythonVirtualEnv()

	//AWS
	awsRole := getAwsInfo()
	iExpire, _ := strconv.ParseInt(os.Getenv("AWS_SESSION_EXPIRE"), 10, 0)
	awsExpire := time.Unix(iExpire, int64(0))

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

	//Formatting
	var userInfo, pwdInfo, virtualEnvInfo, awsInfo string

	if user == "root" {
		userInfo = termcolor.Format(hostname, termcolor.Bold, termcolor.FgRed)
	} else {
		userInfo = termcolor.Format(hostname, termcolor.Bold, termcolor.FgGreen)
	}
	pwdInfo = termcolor.Format(pwd, termcolor.Bold, termcolor.FgBlue)
	virtualEnvInfo = termcolor.Format(virtualEnv, termcolor.FgBlue)

	if awsRole != "" {
		t := termcolor.FgGreen
		d := time.Until(awsExpire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		awsInfo = termcolor.Format(awsRole, t) + "|"
	}

	fmt.Printf("%s[%s%s %s]$ ", virtualEnvInfo, awsInfo, userInfo, pwdInfo)
}
