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

type awsInfo struct {
	role   string
	expire time.Time
}

type termInfo struct {
	lastrc     string
	pwd        string
	user       string
	hostname   string
	virtualEnv string
}

type promptInfo struct {
	term termInfo
	aws  awsInfo
	git  gitInfo
}

func getPythonVirtualEnv() string {
	virtualEnv, ve := os.LookupEnv("VIRTUAL_ENV")
	if ve {
		ave := strings.Split(virtualEnv, "/")
		virtualEnv = ave[len(ave)-1]
	}
	return virtualEnv
}

func getAwsInfo() awsInfo {
	ai := awsInfo{}
	role := os.Getenv("AWS_ROLE")
	if role != "" {
		tmp := strings.Split(role, ":")
		role = tmp[0]
		tmp = strings.Split(tmp[1], "-")
		role += ":" + tmp[2]
	}
	ai.role = role
	iExpire, _ := strconv.ParseInt(os.Getenv("AWS_SESSION_EXPIRE"), 10, 0)
	ai.expire = time.Unix(iExpire, int64(0))

	return ai
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
			//Probably there are no commits yet. How to know the current branch??
			gi.branch = "No_Commits"
			return gi
			//log.Fatalln("error getting head: ", err)
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
				if err != nil {
					log.Fatalln("Error getting merge bases")
				}
				gi.commitsAhead, gi.commitsBehind, err = repository.AheadBehind(localRef.Target(), remoteRef.Target())
				if err != nil {
					log.Fatalln("Error getting commits ahead/behind")
				}
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

func getTermInfo() termInfo {
	ti := termInfo{}
	//Get basicinfo
	ti.pwd = os.Getenv("PWD")
	home := os.Getenv("HOME")
	if home != "" {
		ti.pwd = strings.Replace(ti.pwd, home, "~", -1)
	}
	ti.user = os.Getenv("USER")
	ti.hostname = os.Getenv("HOSTNAME")

	ti.lastrc = os.Getenv("LAST_COMMAND_RC")

	//Get Python VirtualEnv info
	ti.virtualEnv = getPythonVirtualEnv()

	return ti

}

type promptOptions struct {
	style      string
	fullpath   bool
	noFullpath bool
	color      bool
	noColor    bool
}

func main() {
	var debug bool
	var po promptOptions
	var style string

	flag.BoolVar(&debug, "debug", false, "enable debug messages")
	flag.StringVar(&po.style, "style", "Evermeet", "Select style: Evermeet, Mac, Fedora")
	flag.BoolVar(&po.fullpath, "fullpath", true, "Show fullpath on prompt. Depends on the style")
	flag.BoolVar(&po.noFullpath, "no-fullpath", false, "Show fullpath on prompt. The default value depends on the style")
	flag.BoolVar(&po.color, "color", true, "Show color on prompt. The default value depends on the style")
	flag.BoolVar(&po.noColor, "no-color", false, "Show color on prompt. The default value depends on the style")
	flag.Parse()

	/*colorSet := flag.CommandLine.LookupActual("color")
	noColorSet := flag.CommandLine.LookupActual("color")
	fullpathSet := flag.CommandLine.LookupActual("color")
	noFullpathSet := flag.CommandLine.LookupActual("color")*/

	switch style {
	case "Evermeet":
	case "Mac":
	case "Fedora":
	default:
		fmt.Fprintln(os.Stderr, "Invalid style. Valid styles: Evermeet, Mac, Fedora")
	}
	logger = log.New(os.Stderr, "", log.LstdFlags)

	if !debug {
		logger.SetOutput(ioutil.Discard)
	}
	ti := getTermInfo()

	//AWS
	ai := getAwsInfo()

	//Get git information
	gi := getGitInfo()

	pi := promptInfo{term: ti, git: gi, aws: ai}

	fmt.Println(makePrompt(style, pi))

}

func makePrompt(style string, pi promptInfo) string {
	switch style {
	case "Evermeet":
		return makePromptEvermeet(pi)
	case "Mac":
		return makePromptMac(pi)
	case "Fedora":
		return makePromptFedora(pi)
	}
	return "Not suppported"
}
func makePromptMac(pi promptInfo) string {
	return "Not implemented"
}

func makePromptFedora(pi promptInfo) string {
	return "Not implemented"
}

func makePromptEvermeet(pi promptInfo) string {
	//Formatting
	var userInfo, lastCommandInfo, pwdInfo, virtualEnvInfo, awsInfo, gitInfo string

	promptEnd := "$"

	if pi.term.user == "root" {
		userInfo = termcolor.EscapedFormat(pi.term.hostname, termcolor.Bold, termcolor.FgRed)
		promptEnd = "#"
	} else {
		userInfo = termcolor.EscapedFormat(pi.term.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	if pi.term.lastrc != "" {
		lastCommandInfo = termcolor.EscapedFormat(pi.term.lastrc, termcolor.FgHiYellow) + " "
	}

	pwdInfo = termcolor.EscapedFormat(pi.term.pwd, termcolor.Bold, termcolor.FgBlue)
	if pi.term.virtualEnv != "" {
		virtualEnvInfo = termcolor.EscapedFormat(fmt.Sprintf("(%s) ", pi.term.virtualEnv), termcolor.FgBlue)
	}
	if pi.git.branch != "" {
		gitInfo = " " + termcolor.EscapedFormat(pi.git.branch, termcolor.FgMagenta)
		space := " "
		if pi.git.commitsBehind > 0 {
			gitInfo += space + s_DownArrow + "·" + strconv.Itoa(pi.git.commitsBehind)
			space = ""
		}
		if pi.git.commitsAhead > 0 {
			gitInfo += space + s_UpArrow + "·" + strconv.Itoa(pi.git.commitsAhead)
			space = ""
		}
		if !pi.git.upstream {
			gitInfo += space + "*"
			space = ""
		}
		gitInfo += "|"
		synced := true
		if pi.git.staged > 0 {
			gitInfo += termcolor.EscapedFormat(s_Dot+strconv.Itoa(pi.git.staged), termcolor.FgCyan)
			synced = false
		}
		if pi.git.changed > 0 {
			gitInfo += termcolor.EscapedFormat("+"+strconv.Itoa(pi.git.changed), termcolor.FgCyan)
			synced = false
		}
		if pi.git.untracked > 0 {
			gitInfo += termcolor.EscapedFormat(s_ThreeDots+strconv.Itoa(pi.git.untracked), termcolor.FgCyan)
			synced = false
		}
		if pi.git.stashed > 0 {
			gitInfo += termcolor.EscapedFormat(s_Flag+strconv.Itoa(pi.git.stashed), termcolor.FgHiMagenta)
		}
		if synced {
			gitInfo += termcolor.EscapedFormat(s_Check, termcolor.FgHiGreen)
		}
	}
	if pi.aws.role != "" {
		t := termcolor.FgGreen
		d := time.Until(pi.aws.expire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		awsInfo = termcolor.EscapedFormat(pi.aws.role, t) + "|"
	}

	return fmt.Sprintf("%s[%s%s %s%s%s]%s ", virtualEnvInfo, awsInfo, userInfo, lastCommandInfo, pwdInfo, gitInfo, promptEnd)
}
