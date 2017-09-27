package prompt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/josledp/termcolor"
)

const (
	sDownArrow = "↓"
	sUpArrow   = "↑"
	sThreeDots = "…"
	sDot       = "●"
	sCheck     = "✔"
	sFlag      = "⚑"
	sAsterisk  = "⭑"
	sCross     = "✖"
)

//The Prompt object
type Prompt struct {
	format   func(string, ...termcolor.Mode) string
	style    string
	fullpath bool
	term     termInfo
	aws      awsInfo
	git      gitInfo
}

//New returns a new Prompt object
func New(style string, pColor *bool, pFullpath *bool, noFetch bool) Prompt {
	return newWithInfo(style, pColor, pFullpath, getTermInfo(), getAwsInfo(), getGitInfo(noFetch))
}

func newWithInfo(style string, pColor *bool, pFullpath *bool, ti termInfo, ai awsInfo, gi gitInfo) Prompt {

	pr := Prompt{}
	pr.style = style
	pr.term = ti
	pr.aws = ai
	pr.git = gi

	var color bool

	switch style {
	case "Evermeet":
		color = true
		pr.fullpath = true
	case "Mac":
		color = false
		pr.fullpath = true
	case "Fedora":
		color = false
		pr.fullpath = false
	default:
		fmt.Fprintln(os.Stderr, "Invalid style. Valid styles: Evermmet, Mac, Fedora")
	}
	if pColor != nil {
		color = *pColor
	}
	if pFullpath != nil {
		pr.fullpath = *pFullpath
	}

	if color {
		pr.format = termcolor.EscapedFormat
	} else {

		pr.format = func(s string, modes ...termcolor.Mode) string { return s }
	}
	return pr
}

//GetPrompt returns prompt
func (pr Prompt) GetPrompt() string {

	switch pr.style {
	case "Evermeet":
		return pr.makePromptEvermeet()
	case "Mac":
		return pr.makePromptMac()
	case "Fedora":
		return pr.makePromptFedora()
	}
	return "Not suppported"
}

func (pr Prompt) makePromptMac() string {
	return "Not implemented"
}

func (pr Prompt) makePromptFedora() string {
	//Formatting
	var userPromptInfo, lastCommandPromptInfo, pwdPromptInfo, virtualEnvPromptInfo, awsPromptInfo, gitPromptInfo string

	promptEnd := "$"

	if pr.term.user == "root" {
		userPromptInfo = pr.format(pr.term.user+"@"+pr.term.hostname, termcolor.Bold, termcolor.FgRed)
		promptEnd = "#"
	} else {
		userPromptInfo = pr.format(pr.term.user+"@"+pr.term.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	if pr.term.lastrc != "" {
		lastCommandPromptInfo = pr.format(pr.term.lastrc, termcolor.FgHiYellow) + " "
	}
	if !pr.fullpath {
		pwd := strings.Split(pr.term.pwd, "/")
		pr.term.pwd = pwd[len(pwd)-1]
	}
	pwdPromptInfo = pr.format(pr.term.pwd, termcolor.Bold, termcolor.FgBlue)
	if pr.term.virtualEnv != "" {
		virtualEnvPromptInfo = pr.format(fmt.Sprintf("(%s) ", pr.term.virtualEnv), termcolor.FgBlue)
	}
	gitPromptInfo = pr.makeGitPrompt()

	if pr.aws.role != "" {
		t := termcolor.FgGreen
		d := time.Until(pr.aws.expire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		awsPromptInfo = pr.format(pr.aws.role, t) + "|"
	}

	return fmt.Sprintf("[%s%s%s %s%s%s]%s ", virtualEnvPromptInfo, awsPromptInfo, userPromptInfo, lastCommandPromptInfo, pwdPromptInfo, gitPromptInfo, promptEnd)
}

func (pr Prompt) makePromptEvermeet() string {
	//Formatting
	var userPromptInfo, lastCommandPromptInfo, pwdPromptInfo, virtualEnvPromptInfo, awsPromptInfo, gitPromptInfo string

	promptEnd := "$"

	if pr.term.user == "root" {
		userPromptInfo = pr.format(pr.term.hostname, termcolor.Bold, termcolor.FgRed)
		promptEnd = "#"
	} else {
		userPromptInfo = pr.format(pr.term.user+"@"+pr.term.hostname, termcolor.Bold, termcolor.FgGreen)
	}
	if pr.term.lastrc != "" {
		lastCommandPromptInfo = pr.format(pr.term.lastrc, termcolor.FgHiYellow) + " "
	}
	if !pr.fullpath {
		pwd := strings.Split(pr.term.pwd, "/")
		pr.term.pwd = pwd[len(pwd)-1]
	}
	pwdPromptInfo = pr.format(pr.term.pwd, termcolor.Bold, termcolor.FgBlue)
	if pr.term.virtualEnv != "" {
		virtualEnvPromptInfo = pr.format(fmt.Sprintf("(%s) ", pr.term.virtualEnv), termcolor.FgBlue)
	}
	gitPromptInfo = pr.makeGitPrompt()

	if pr.aws.role != "" {
		t := termcolor.FgGreen
		d := time.Until(pr.aws.expire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		awsPromptInfo = pr.format(pr.aws.role, t) + "|"
	}

	return fmt.Sprintf("%s%s%s %s%s%s%s ", virtualEnvPromptInfo, awsPromptInfo, userPromptInfo, lastCommandPromptInfo, pwdPromptInfo, gitPromptInfo, promptEnd)
}

func (pr Prompt) makeGitPrompt() string {
	var gitPromptInfo string
	if pr.git.branch != "" {
		gitPromptInfo = " " + pr.format(pr.git.branch, termcolor.FgMagenta)
		space := " "
		if pr.git.commitsBehind > 0 {
			gitPromptInfo += space + sDownArrow + "·" + strconv.Itoa(pr.git.commitsBehind)
			space = ""
		}
		if pr.git.commitsAhead > 0 {
			gitPromptInfo += space + sUpArrow + "·" + strconv.Itoa(pr.git.commitsAhead)
			space = ""
		}
		if !pr.git.hasUpstream {
			gitPromptInfo += space + sAsterisk
			space = ""
		}
		gitPromptInfo += "|"
		synced := true
		if pr.git.conflicted > 0 {
			gitPromptInfo += pr.format(sCross+strconv.Itoa(pr.git.conflicted), termcolor.FgRed)
			synced = false
		}
		if pr.git.staged > 0 {
			gitPromptInfo += pr.format(sDot+strconv.Itoa(pr.git.staged), termcolor.FgCyan)
			synced = false
		}
		if pr.git.changed > 0 {
			gitPromptInfo += pr.format("+"+strconv.Itoa(pr.git.changed), termcolor.FgCyan)
			synced = false
		}
		if pr.git.untracked > 0 {
			gitPromptInfo += pr.format(sThreeDots+strconv.Itoa(pr.git.untracked), termcolor.FgCyan)
			synced = false
		}
		if synced {
			gitPromptInfo += pr.format(sCheck, termcolor.FgHiGreen)
		}
		if pr.git.stashed > 0 {
			gitPromptInfo += pr.format(sFlag+strconv.Itoa(pr.git.stashed), termcolor.FgHiMagenta)
		}
	}
	return gitPromptInfo
}
