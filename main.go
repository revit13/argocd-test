package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
)

func TemporalDir() (path string, clean func()) {
	fs := osfs.New(os.TempDir())
	path, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	return fs.Join(fs.Root(), path), func() {
		util.RemoveAll(fs, path)
	}
}

func publicKey(filePath string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshKey, _ := ioutil.ReadFile(filePath)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

// Example of how to open a repository in a specific path, and push to
// its default remote (origin).
func main() {
	//auth, _ := publicKey("/home/eres/.ssh/id_rsa.pub")
	r, err := git.PlainClone("/data/stam12/", false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "revit13",
			Password: "e4dea7c7048d04c1e2bcb2e144764e27a398e7d8",
		},
		URL:      "https://github.com/revit13/argocd-test.git",
		Progress: os.Stdout,
	})
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// ... we need a file to commit so let's create a new file inside of the
	// worktree of the project using the go standard library.
	Info("echo \"hello world!\" > example-git-file")
	filename := filepath.Join("/data/stam12/", "example-git-file")
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	CheckIfError(err)

	// Adds the new file to the staging area.
	Info("git add example-git-file")
	_, err = w.Add("example-git-file")
	CheckIfError(err)

	// We can verify the current status of the worktree using the method Status.
	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	Info("git commit -m \"example go-git commit\"")
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Revital Sur",
			Email: "eres@il.ibm.com",
			When:  time.Now(),
		},
	})

	remote, err := r.Remote("origin")

	fmt.Println(commit)
	po := &git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: "revit13",
			Password: "e4dea7c7048d04c1e2bcb2e144764e27a398e7d8",
		},
		RemoteName:      "origin",
		RefSpecs:        []config.RefSpec{config.RefSpec("refs/heads/*:refs/heads/*")},
		Progress:        os.Stdout,
		Force:           false,
		InsecureSkipTLS: true,
	}

	err = remote.Push(po)

	fmt.Println(err)
	CheckIfError(err)

}
