package plugin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/josledp/termcolor"
	"gopkg.in/src-d/go-billy.v3/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const (
	sDownArrow = "↓"
	sUpArrow   = "↑"
	sThreeDots = "…"
	sDot       = "●"
	sCheck     = "✔"
	sFlag      = "⚑"
	sAsterisk  = "*"
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

func findGitRepository() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error determining current pwd: %v", err)
	}
	p := strings.Split(pwd, "/")
	for i := range p {
		path := strings.Join(p[:len(p)-i], "/")
		if info, err := os.Stat(path + "/.git"); os.IsNotExist(err) {
			continue
		} else if err == nil && info.IsDir() {
			return path, nil
		}
	}
	return "", fmt.Errorf("unable to find .git directory")
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

//Load is the load function of the plugin
func (g *Git) Load(pr Prompter) error {
	gitPwd, err := findGitRepository()
	if err != nil {
		return nil //Unable to find valid git repo, so good so far
	}
	fs := osfs.New(gitPwd + "/.git")
	storage, err := filesystem.NewStorage(fs)
	if err != nil {
		return fmt.Errorf("unable to get storer: %v", err)
	}
	repository, err := git.Open(storage, osfs.New(gitPwd))
	if err != nil {
		return fmt.Errorf("unable to open repository: %v", err)
	}
	wt, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}
	st, err := wt.Status()
	if err != nil {
		return fmt.Errorf("error getting worktree status: %v", err)
	}
	for _, status := range st {
		switch status.Staging {
		case git.Unmodified:
		case git.Untracked:
		default:
			g.staged++
		}

		switch status.Worktree {
		case git.Unmodified:
		case git.Untracked:
			g.untracked++
		default:
			g.changed++
		}
	}
	//TODO: missing conflict files!
	if fstash, err := os.Open(gitPwd + "/.git/logs/refs/stash"); err == nil {
		defer fstash.Close()
		g.stashed, err = lineCounter(fstash)
		if err != nil {
			return fmt.Errorf("unable to count stashes:%v", err)
		}
	}
	localRef, err := repository.Head()
	if err != nil {
		g.branch = "No_Commits"
		return nil
	}
	localName := strings.Split(localRef.Name().String(), "/")
	if len(localName) == 1 {
		g.branch = ":" + localRef.Hash().String()[:7]
		g.detached = true
	} else {
		g.branch = localName[len(localName)-1]
	}
	c, err := repository.Config()
	remote := c.Raw.Section("branch").Subsection(g.branch).Option("remote")
	if remote != "" {
		g.hasUpstream = true
		g.fetchIfNeeded(pr)
		remoteRef, err := repository.Reference(plumbing.ReferenceName("refs/remotes/"+remote+"/"+g.branch), false)
		if localRef != remoteRef && err == nil {
			localCo, err := repository.CommitObject(localRef.Hash())
			if err != nil {
				return fmt.Errorf("unable to get local commit from local reference: %v", err)
			}
			remoteCo, err := repository.CommitObject(localRef.Hash())
			if err != nil {
				return fmt.Errorf("unable to get local commit from remote reference: %v", err)
			}
			g.commitsAhead, g.commitsBehind = aheadBehind(repository, localCo, remoteCo)
		}

	}
	return nil
}

func fillMap(r *git.Repository, co *object.Commit, m map[string]struct{}) {
	log.Print(co.Hash)
	m[co.Hash] = struct{}{}
	for _, _p := range co.ParentHashes {
		p, _ := r.CommitObject(_p)
		fillMap(r, p, m)
	}
}
func count(r *git.Repository, co *object.Commit, m map[string]struct{}) int {
	if _, ok := m[co.Hash]; ok {
		return 0
	}
	c := 1
	for _, _p := range co.ParentHashes {
		p, _ := r.CommitObject(_p)
		c += count(r, p, m)
	}
	return c
}
func aheadBehind(repository *git.Repository, local *object.Commit, remote *object.Commit) (ahead, behind int) {
	localMap := make(map[string]struct{})
	remoteMap := make(map[string]struct{})
	log.Print("local")
	fillMap(repository, local, localMap)
	log.Print("remote")
	fillMap(repository, remote, remoteMap)
	ahead = count(repository, local, remoteMap)
	behind = count(repository, remote, localMap)
	return ahead, behind
}

/*

			if entry.Status&git2go.StatusConflicted > 0 {
				g.conflicted++
				got = true
			}

		}


		remoteRef, err := localBranch.Upstream()

		if err == nil {
			defer remoteRef.Free()


			if !remoteRef.Target().Equal(localRef.Target()) {
				g.commitsAhead, g.commitsBehind, err = repository.AheadBehind(localRef.Target(), remoteRef.Target())
				if err != nil {
					return fmt.Errorf("Error getting commitsAhead/Behing: %v", err)
				}
			}
		}
		})
	}
	return nil*/

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

func (g *Git) fetchIfNeeded(pr Prompter) {
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
}
