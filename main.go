package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	git2go "gopkg.in/libgit2/git2go.v26"

	"github.com/josledp/termcolor"
)

const (
	s_DownArrow = "↓"
	s_UpArrow   = "↑"
	s_ThreeDots = "…"
	s_Dot       = "●"
	s_Check     = "✔"
	s_Flag      = "⚑"
)

var logger *log.Logger

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
				logger.Println("StatusCurrent")
				got = true
			}
			if entry.Status&git2go.StatusIndexNew > 0 {
				logger.Println("StatusIndexNew")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexModified > 0 {
				logger.Println("StatusIndexModified")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexDeleted > 0 {
				logger.Println("StatusIndexDeleted")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexRenamed > 0 {
				logger.Println("StatusIndexRenamed")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexTypeChange > 0 {
				logger.Println("StatusIndexTypeChange")
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusWtNew > 0 {
				logger.Println("StatusWtNew")
				gi.untracked++
				got = true
			}
			if entry.Status&git2go.StatusWtModified > 0 {
				logger.Println("StatusWtModified")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtDeleted > 0 {
				logger.Println("StatusWtDeleted")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtTypeChange > 0 {
				logger.Println("StatusWtTypeChange")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtRenamed > 0 {
				logger.Println("StatusWtRenamed")
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusIgnored > 0 {
				logger.Println("StatusIgnored")
				got = true
			}
			if entry.Status&git2go.StatusConflicted > 0 {
				logger.Println("StatusConflicted")
				gi.conflict = true
				got = true
			}
			if !got {
				logger.Println("Unknown: ", entry.Status)
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
			gi.upstream = true
			defer remoteRef.Free()

			if !remoteRef.Target().Equal(localRef.Target()) {
				logger.Println("Local & remore differ:", remoteRef.Target().String(), localRef.Target().String())
				//git rev-list --left-right localRef...remoteRef
				oids, err := repository.MergeBases(localRef.Target(), remoteRef.Target())
				if err != nil {
					log.Fatalln("Error getting merge bases")
				}

				gi.commitsAhead = gitCount(repository, localRef.Target(), oids)
				gi.commitsBehind = gitCount(repository, remoteRef.Target(), oids)
				logger.Println(gi.commitsAhead, gi.commitsBehind)
			}
		}
		// stash
		repository.Stashes.Foreach(func(i int, m string, o *git2go.Oid) error {
			gi.stashed = i + 1
			return nil
		})
		logger.Println("Stashes: ", gi.stashed)
	}
	return gi
}

func gitCount(r *git2go.Repository, oid *git2go.Oid, until []*git2go.Oid) int {
	c, err := r.LookupCommit(oid)
	defer c.Free()
	if err != nil {
		log.Fatalln("Error getting commit from oid ", oid, ": ", err)
	}
	mUntil := make(map[string]struct{})
	for _, u := range until {
		mUntil[u.String()] = struct{}{}
	}
	return _gitCount(r, c, mUntil)

}
func _gitCount(r *git2go.Repository, c *git2go.Commit, until map[string]struct{}) int {
	var s int
	for i := uint(0); i < c.ParentCount(); i++ {
		s++
		pc := c.ParentId(i)
		if _, ok := until[pc.String()]; !ok {
			s += _gitCount(r, c.Parent(i), until)
		}

	}
	return s
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
	lastrc     string
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
	var debug bool

	flag.BoolVar(&debug, "debug", false, "enable debug messages")
	flag.Parse()
	logger = log.New(os.Stderr, "", log.LstdFlags)

	if !debug {
		logger.SetOutput(ioutil.Discard)
	}
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
	ti.lastrc = os.Getenv("LAST_COMMAND_RC")

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
	var userInfo, lastCommandInfo, pwdInfo, virtualEnvInfo, awsInfo, gitInfo string

	promptEnd := "$"

	if ti.user == "root" {
		userInfo = termcolor.EscapedFormat(ti.hostname, termcolor.Bold, termcolor.FgRed)
		promptEnd = "#"
	} else {
		userInfo = termcolor.EscapedFormat(ti.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	if ti.lastrc != "" {
		lastCommandInfo = termcolor.EscapedFormat(ti.lastrc, termcolor.FgHiYellow) + " "
	}

	pwdInfo = termcolor.EscapedFormat(ti.pwd, termcolor.Bold, termcolor.FgBlue)
	if ti.virtualEnv != "" {
		virtualEnvInfo = termcolor.EscapedFormat(ti.virtualEnv, termcolor.FgBlue)
	}
	if ti.gi.branch != "" {
		gitInfo = " " + termcolor.EscapedFormat(ti.gi.branch, termcolor.FgMagenta)
		space := " "
		if ti.gi.commitsBehind > 0 {
			gitInfo += space + s_DownArrow + "·" + strconv.Itoa(ti.gi.commitsBehind)
			space = ""
		}
		if ti.gi.commitsAhead > 0 {
			gitInfo += space + s_UpArrow + "·" + strconv.Itoa(ti.gi.commitsAhead)
			space = ""
		}
		if !ti.gi.upstream {
			gitInfo += space + "*"
			space = ""
		}
		gitInfo += "|"
		synced := true
		if ti.gi.staged > 0 {
			gitInfo += termcolor.EscapedFormat(s_Dot+strconv.Itoa(ti.gi.staged), termcolor.FgCyan)
			synced = false
		}
		if ti.gi.changed > 0 {
			gitInfo += termcolor.EscapedFormat("+"+strconv.Itoa(ti.gi.changed), termcolor.FgCyan)
			synced = false
		}
		if ti.gi.untracked > 0 {
			gitInfo += termcolor.EscapedFormat(s_ThreeDots+strconv.Itoa(ti.gi.untracked), termcolor.FgCyan)
			synced = false
		}
		if ti.gi.stashed > 0 {
			gitInfo += termcolor.EscapedFormat(s_Flag+strconv.Itoa(ti.gi.stashed), termcolor.FgHiMagenta)
		}
		if synced {
			gitInfo += termcolor.EscapedFormat(s_Check, termcolor.FgHiGreen)
		}
	}
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

	return fmt.Sprintf("%s[%s%s %s%s%s]%s ", virtualEnvInfo, awsInfo, userInfo, lastCommandInfo, pwdInfo, gitInfo, promptEnd)
}
