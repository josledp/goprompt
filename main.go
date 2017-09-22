package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	git2go "gopkg.in/libgit2/git2go.v26"

	"github.com/josledp/termcolor"
)

const (
	downArrow   = "↓"
	upArrow     = "↑"
	threePoints = "…"
	dot         = "●"
	check       = "✔"
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

func getGitInfo() gitInfo {
	gi := gitInfo{}

	gitpath, err := git2go.Discover(".", false, []string{"/"})
	if err == nil {
		repository, err := git2go.OpenRepository(gitpath)
		if err != nil {
			log.Fatalf("Error opening repository at %s: %v", gitpath, err)
		}
		defer repository.Free()

		//Get current tracked & untracked files status
		statusOpts := git2go.StatusOptions{
			Flags: git2go.StatusOptIncludeUntracked | git2go.StatusOptRenamesHeadToIndex,
		}
		repostate, err := repository.StatusList(&statusOpts)
		if err != nil {
			log.Fatalf("Error getting repository status at %s: %v", gitpath, err)
		}
		defer repostate.Free()
		n, err := repostate.EntryCount()
		for i := 0; i < n; i++ {
			entry, _ := repostate.ByIndex(i)
			got := false
			if entry.Status&git2go.StatusCurrent > 0 {
				log.Println("StatusCurrent")
				got = true
			}
			if entry.Status&git2go.StatusIndexNew > 0 {
				log.Println("StatusIndexNew")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexModified > 0 {
				log.Println("StatusIndexModified")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexDeleted > 0 {
				log.Println("StatusIndexDeleted")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexRenamed > 0 {
				log.Println("StatusIndexRenamed")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexTypeChange > 0 {
				log.Println("StatusIndexTypeChange")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusWtNew > 0 {
				log.Println("StatusWtNew")
				gi.untracked++
				got = true
			}
			if entry.Status&git2go.StatusWtModified > 0 {
				log.Println("StatusWtModified")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtDeleted > 0 {
				log.Println("StatusWtDeleted")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtTypeChange > 0 {
				log.Println("StatusWtTypeChange")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtRenamed > 0 {
				log.Println("StatusWtRenamed")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusIgnored > 0 {
				log.Println("StatusIgnored")
				got = true
			}
			if entry.Status&git2go.StatusConflicted > 0 {
				log.Println("StatusConflicted")
				gi.conflict = true
				got = true
			}
			if !got {
				log.Println("Unknown: ", entry.Status)
			}
		}
		//Get current branch name
		localRef, err := repository.Head()
		if err != nil {
			log.Fatalln("error getting head: ", err)
		}
		defer localRef.Free()

		ref := strings.Split(localRef.Name(), "/")
		gi.branch = ref[len(ref)-1]
		//Get commits Ahead/Behind

		localBranch := localRef.Branch()
		if err != nil {
			log.Fatalln("Error getting local branch: ", err)
		}

		remoteRef, err := localBranch.Upstream()
		if err == nil {
			defer remoteRef.Free()

			if !remoteRef.Target().Equal(localRef.Target()) {
				log.Println("Local & remore differ:", remoteRef.Target().String(), localRef.Target().String())
				//git rev-list --left-right localRef...remoteRef
				oids, err := repository.MergeBases(localRef.Target(), remoteRef.Target())
				if err != nil {
					log.Fatalln("Error getting merge bases")
				}
				for _, oid := range oids {
					log.Println(oid.String())
				}
			}
		}

	}
	return gi
}

type gitInfo struct {
	conflict      bool
	changed       int
	staged        int
	untracked     int
	commitsAhead  int
	commitsBehind int
	stashed       int
	branch        string
	upstream      bool
}

type termInfo struct {
	pwd        string
	user       string
	hostname   string
	virtualEnv string
	awsRole    string
	awsExpire  time.Time
	gi         gitInfo
}

func main() {
	var err error

	ti := termInfo{}
	//Get basicinfo
	ti.pwd, err = os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current path", err)
	}
	home := os.Getenv("HOME")
	if home != "" {
		ti.pwd = strings.Replace(ti.pwd, home, "~", -1)
	}
	ti.user = os.Getenv("USER")
	ti.hostname, err = os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname", err)
	}

	//Get Python VirtualEnv info
	ti.virtualEnv = getPythonVirtualEnv()

	//AWS
	ti.awsRole = getAwsInfo()
	iExpire, _ := strconv.ParseInt(os.Getenv("AWS_SESSION_EXPIRE"), 10, 0)
	ti.awsExpire = time.Unix(iExpire, int64(0))

	//Get git information
	_ = git2go.Repository{}

	ti.gi = getGitInfo()

	fmt.Println(makePrompt(ti))
}

func makePrompt(ti termInfo) string {
	//Formatting
	var userInfo, pwdInfo, virtualEnvInfo, awsInfo string
	promptEnd := "$"

	if ti.user == "root" {
		userInfo = termcolor.EscapedFormat(ti.hostname, termcolor.Bold, termcolor.FgRed)
		promptEnd = "#"
	} else {
		userInfo = termcolor.EscapedFormat(ti.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	pwdInfo = termcolor.EscapedFormat(ti.pwd, termcolor.Bold, termcolor.FgBlue)
	virtualEnvInfo = termcolor.EscapedFormat(ti.virtualEnv, termcolor.FgBlue)

	if ti.awsRole != "" {
		t := termcolor.FgGreen
		d := time.Until(ti.awsExpire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		awsInfo = termcolor.EscapedFormat(ti.awsRole, t) + "|"
	}

	return fmt.Sprintf("%s[%s%s %s]%s ", virtualEnvInfo, awsInfo, userInfo, pwdInfo, promptEnd)
}
