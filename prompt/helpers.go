package prompt

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	git2go "gopkg.in/libgit2/git2go.v26"
)

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
