package main

import (
	"github.com/yxzhm/git-sync/util"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	goconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"strings"
)

// Example of how to:
// - Clone a repository into memory
// - Get the HEAD reference
// - Using the HEAD reference, obtain the commit this reference is pointing to
// - Using the commit, obtain its history and print it
func main() {

	config := util.ReadConfigFile("D:/Sirius/go/workspace/src/github.com/yxzhm/git-sync/config/config.json")
	sourceSSHkey, _ := ioutil.ReadFile("D:/Sirius/go/workspace/src/github.com/yxzhm/git-sync/config/git_key")
	singer, _ := ssh.ParsePrivateKey(sourceSSHkey)
	for i := 0; i < len(config.Groups); i++ {
		group := &config.Groups[i]
		for j := 0; j < len(group.Repos); j++ {
			repo := group.Repos[j]
			println(repo)
			sourceURL := "git@" + config.Source + ":" + group.Name + "/" + repo
			targetURL := "git@" + config.Target + ":" + group.Name + "/" + repo
			println(sourceURL)

			r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
				Progress: os.Stdout,
				URL:      sourceURL,
				Auth:     &ssh2.PublicKeys{User: "git", Signer: singer},
			})
			if err != nil {
				println(err)
			}
			println(r)

			_, err = r.CreateRemote(&goconfig.RemoteConfig{
				Name: "Chengdu",
				URLs: []string{targetURL},
			})
			cIter, err := r.References()
			//w, _ := r.Worktree()
			err = cIter.ForEach(func(ref *plumbing.Reference) error {
				if ref.Name() != "HEAD" {
					branchName := strings.Replace(ref.Name().String(), "refs/remotes/origin/", "", 1)

					//w.Checkout(&git.CheckoutOptions{
					//	Force:  true,
					//	Hash:   ref.Hash(),
					//	Branch: plumbing.NewBranchReferenceName(branchName),
					//})

					head := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), ref.Hash())
					//h:=plumbing.NewSymbolicReference(plumbing.HEAD,)

					err = r.Storer.SetReference(head)

					if err != nil {
						return err
					}
				}
				println(ref.Name())
				err = r.Push(&git.PushOptions{
					RemoteName: "Chengdu",
					Auth:       &ssh2.PublicKeys{User: "git", Signer: singer},
				})
				return nil
			})

		}
	}
}
