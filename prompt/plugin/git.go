package plugin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/josledp/termcolor"
	git2go "github.com/jeffwelling/git2go/v37"
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

//Git is the plugin struct
type Git struct {
	conflicted    int
	detached      bool
	changed       int
	staged        int
	untracked     int
	commitsAhead  int
	commitsBehind int
	stashed       int
	branch        string
	hasUpstream   bool
}

//Name returns the plugin name
func (Git) Name() string {
	return "git"
}

//Help returns help information about this plugin
func (Git) Help() (description string, options map[string]string) {
	description = "This plugins show git information in the current git repo"
	return
}

//Load is the load function of the plugin
func (g *Git) Load(pr Prompter) error {
	gitpath, err := git2go.Discover(".", false, []string{"/"})
	if err == nil {
		repository, err := git2go.OpenRepository(gitpath)
		if err != nil {
			return fmt.Errorf("Error opening repository at %s: %v", gitpath, err)
		}
		defer repository.Free()

		//Get current tracked & untracked files status
		statusOpts := git2go.StatusOptions{
			Flags: git2go.StatusOptIncludeUntracked | git2go.StatusOptRenamesHeadToIndex,
		}
		repostate, err := repository.StatusList(&statusOpts)
		if err != nil {
			return fmt.Errorf("Error getting repository status at %s: %v", gitpath, err)
		}
		defer repostate.Free()
		n, err := repostate.EntryCount()
		for i := 0; i < n; i++ {
			entry, _ := repostate.ByIndex(i)
			got := false
			if entry.Status&git2go.StatusCurrent > 0 {
				got = true
			}
			if entry.Status&git2go.StatusIndexNew > 0 {
				g.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexModified > 0 {
				g.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexDeleted > 0 {
				g.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexRenamed > 0 {
				g.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexTypeChange > 0 {
				g.staged++
				got = true
			}
			if entry.Status&git2go.StatusWtNew > 0 {
				g.untracked++
				got = true
			}
			if entry.Status&git2go.StatusWtModified > 0 {
				g.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtDeleted > 0 {
				g.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtTypeChange > 0 {
				g.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtRenamed > 0 {
				g.changed++
				got = true
			}
			if entry.Status&git2go.StatusIgnored > 0 {
				got = true
			}
			if entry.Status&git2go.StatusConflicted > 0 {
				g.conflicted++
				got = true
			}
			if !got {
				log.Println("Git plugin. Unknown: ", entry.Status)
			}
		}
		//Get current branch name
		localRef, err := repository.Head()
		if err != nil {
			//Probably there are no commits yet. How to know the current branch??
			g.branch = "No_Commits"
			return nil
		}
		defer localRef.Free()

		localBranch := localRef.Branch()
		if err != nil {
			return fmt.Errorf("Error getting local branch: %v", err)
		}

		if isHead, _ := localBranch.IsHead(); isHead {
			g.branch = localRef.Shorthand()
		} else {
			g.branch = ":" + localRef.Target().String()[:7]
			g.detached = true
		}

		remoteRef, err := localBranch.Upstream()

		if err == nil {
			defer remoteRef.Free()

			g.hasUpstream = true
			pwd, err := os.Getwd()
			if err == nil {
				key := fmt.Sprintf("git-%s-fetch", pwd)
				last, ok := pr.GetCache(key)
				var lastTime time.Time
				if last != nil {
					lastTime, err = time.Parse(time.RFC3339, last.(string))
					if err != nil {
						log.Printf("Error loading git last fetch time: %v", err)
					}
				}
				if !ok || time.Since(lastTime) > 300*time.Second {
					pa := syscall.ProcAttr{}
					pa.Env = os.Environ()
					pa.Dir = pwd
					gitcommand, err := exec.LookPath("git")
					if err != nil {
						log.Printf("git command not found: %v", err)
					} else {
						_, err = syscall.ForkExec(gitcommand, []string{gitcommand, "fetch"}, &pa)
						if err != nil {
							//Silently fail?
							log.Printf("Error fetching: %v", err)
						} else {
							pr.Cache(key, time.Now())
						}
					}
				}
			}

			if !remoteRef.Target().Equal(localRef.Target()) {
				g.commitsAhead, g.commitsBehind, err = repository.AheadBehind(localRef.Target(), remoteRef.Target())
				if err != nil {
					return fmt.Errorf("Error getting commitsAhead/Behing: %v", err)
				}
			}
		}
		if fstash, err := os.Open(gitpath + "/logs/refs/stash"); err == nil {
			defer fstash.Close()
			g.stashed, err = lineCounter(fstash)
			if err != nil {
				return fmt.Errorf("unable to count stashes:%v", err)
			}
		}
	}
	return nil
}

//Get returns the string to use in the prompt
func (g Git) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	var gitPromptInfo string
	if g.branch != "" {
		gitPromptInfo = format(g.branch, termcolor.FgMagenta)
		space := " "
		if g.commitsBehind > 0 {
			gitPromptInfo += space + sDownArrow + "·" + strconv.Itoa(g.commitsBehind)
			space = ""
		}
		if g.commitsAhead > 0 {
			gitPromptInfo += space + sUpArrow + "·" + strconv.Itoa(g.commitsAhead)
			space = ""
		}
		if !g.hasUpstream {
			gitPromptInfo += space + sAsterisk
			space = ""
		}
		gitPromptInfo += "|"
		synced := true
		if g.conflicted > 0 {
			gitPromptInfo += format(sCross+strconv.Itoa(g.conflicted), termcolor.FgRed)
			synced = false
		}
		if g.staged > 0 {
			gitPromptInfo += format(sDot+strconv.Itoa(g.staged), termcolor.FgCyan)
			synced = false
		}
		if g.changed > 0 {
			gitPromptInfo += format("+"+strconv.Itoa(g.changed), termcolor.FgCyan)
			synced = false
		}
		if g.untracked > 0 {
			gitPromptInfo += format(sThreeDots+strconv.Itoa(g.untracked), termcolor.FgCyan)
			synced = false
		}
		if synced {
			gitPromptInfo += format(sCheck, termcolor.FgHiGreen)
		}
		if g.stashed > 0 {
			gitPromptInfo += format(sFlag+strconv.Itoa(g.stashed), termcolor.FgHiMagenta)
		}
	}
	return gitPromptInfo, []termcolor.Mode{termcolor.FgMagenta}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
