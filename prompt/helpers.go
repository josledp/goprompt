package prompt

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	git2go "gopkg.in/libgit2/git2go.v26"
)

type gitInfo struct {
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
				got = true
			}
			if entry.Status&git2go.StatusIndexNew > 0 {
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexModified > 0 {
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexDeleted > 0 {
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexRenamed > 0 {
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusIndexTypeChange > 0 {
				gi.staged++
				got = true
			}
			if entry.Status&git2go.StatusWtNew > 0 {
				gi.untracked++
				got = true
			}
			if entry.Status&git2go.StatusWtModified > 0 {
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtDeleted > 0 {
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtTypeChange > 0 {
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusWtRenamed > 0 {
				gi.changed++
				got = true
			}
			if entry.Status&git2go.StatusIgnored > 0 {
				got = true
			}
			if entry.Status&git2go.StatusConflicted > 0 {
				gi.conflicted++
				got = true
			}
			if !got {
				log.Println("Unknown: ", entry.Status)
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

		//Get commits Ahead/Behind

		localBranch := localRef.Branch()
		if err != nil {
			log.Fatalln("Error getting local branch: ", err)
		}

		if isHead, _ := localBranch.IsHead(); isHead {
			gi.branch = localRef.Shorthand()
		} else {
			gi.branch = localRef.Target().String()[:7]
			gi.detached = true
		}

		remoteRef, err := localBranch.Upstream()

		if err == nil {
			gi.hasUpstream = true
			remoteBranchName, err := remoteRef.Branch().Name()
			if err != nil {
				log.Println("Error getting branch name", err)
			} else {
				upstream := strings.Split(remoteBranchName, "/")[0]
				remote, err := repository.Remotes.Lookup(upstream)
				defer remote.Free()
				if err != nil {
					log.Println("Error looking for remote Upstream: ", err)
				} else {
					//It does not work with authentication
					//It does not work using host alias on .ssh/config
					cb := git2go.RemoteCallbacks{}
					cb.CertificateCheckCallback = func(*git2go.Certificate, bool, string) git2go.ErrorCode { return git2go.ErrOk }
					cb.CredentialsCallback = func(url string, username_from_url string, allowed_types git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
						return git2go.ErrOk, nil
					}
					err = remote.Fetch([]string{}, &git2go.FetchOptions{RemoteCallbacks: cb}, "")
					if err != nil {
						log.Println("error connecting fetching remote: ", err)
					}

				}
			}

			defer remoteRef.Free()

			if !remoteRef.Target().Equal(localRef.Target()) {
				if err != nil {
					log.Fatalln("Error getting merge bases")
				}
				gi.commitsAhead, gi.commitsBehind, err = repository.AheadBehind(localRef.Target(), remoteRef.Target())
				if err != nil {
					log.Fatalln("Error getting commits ahead/behind")
				}
			}
		}
		// only works if libgit >= 0.25
		repository.Stashes.Foreach(func(i int, m string, o *git2go.Oid) error {
			gi.stashed = i + 1
			return nil
		})
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
	ti.hostname, _ = os.Hostname()

	ti.lastrc = os.Getenv("LAST_COMMAND_RC")

	//Get Python VirtualEnv info
	ti.virtualEnv = getPythonVirtualEnv()

	return ti

}
