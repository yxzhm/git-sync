package main

import (
	"flag"
	"gopkg.in/src-d/go-git.v4"
	goconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"strings"
	"sync"
)

var (
	TargetRemoteName = "Chengdu"
	GitUser          = "git"
	wg               sync.WaitGroup
	ch               chan int
)

func syncRepo(config Config, groupName string, targetRroupName string, repoName string) error {
	Info.Printf("Handling the Group: %s, the Repo: %s", groupName, repoName)
	sourceGit := "git@" + config.SourceURL + ":" + groupName + "/" + repoName
	targetGit := "git@" + config.TargetURL + ":" + groupName + "/" + repoName + ".git"
	if targetRroupName != "" {
		targetGit = "git@" + config.TargetURL + ":" + targetRroupName + "/" + repoName + ".git"
	}

	//Clone Source
	source, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  sourceGit,
		Auth: &ssh2.PublicKeys{User: GitUser, Signer: config.SourcePrivateKey},
	})
	if err != nil {
		return err
	}
	_, err = source.CreateRemote(&goconfig.RemoteConfig{
		Name: TargetRemoteName,
		URLs: []string{targetGit},
	})

	if err != nil {
		return err
	}

	//Clone Target
	target, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  targetGit,
		Auth: &ssh2.PublicKeys{User: GitUser, Signer: config.TargetPrivateKey},
	})

	//if err != nil {
	//	return err
	//}

	cIter, err := source.References()

	if err != nil {
		return err
	}

	// Hading branches
	err = cIter.ForEach(func(ref *plumbing.Reference) error {
		if ref != nil && strings.HasPrefix(ref.Name().String(), "refs/remotes/origin/") {
			branchName := strings.Replace(ref.Name().String(), "refs/remotes/origin/", "", 1)
			Info.Printf("Handling the Group: %s, the Repo: %s, Branch: %s", groupName, repoName, branchName)
			head := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), ref.Hash())
			err = source.Storer.SetReference(head)

			if err != nil {
				Warning.Printf("Set Reference fail.")
				return err
			}
			err = source.Push(&git.PushOptions{
				RemoteName: TargetRemoteName,
				Auth:       &ssh2.PublicKeys{User: GitUser, Signer: config.TargetPrivateKey},
			})
			//err.Error() == "non-fast-forward update: refs/heads/"+branchName
			if err != nil && strings.HasPrefix(err.Error(), "non-fast-forward update") && config.ForcePush {
				Warning.Printf("Need to remove the Target branch " + branchName + " firstly.")
				err = target.Push(&git.PushOptions{
					RefSpecs: []goconfig.RefSpec{goconfig.RefSpec(":refs/heads/" + branchName)},
					Auth:     &ssh2.PublicKeys{User: GitUser, Signer: config.TargetPrivateKey},
				})

				if err != nil {
					return err
				}
				Warning.Printf("Removed " + branchName)
				err = source.Push(&git.PushOptions{
					RemoteName: TargetRemoteName,
					Auth:       &ssh2.PublicKeys{User: GitUser, Signer: config.TargetPrivateKey},
				})
			}

			if err != git.NoErrAlreadyUpToDate {
				return err
			}

			Info.Printf("Pushed Group: %s, the Repo: %s, Branch: %s", groupName, repoName, branchName)
		}
		return nil
	})

	if err != nil {
		return err
	}

	//Handing tags
	err = source.Push(&git.PushOptions{
		RefSpecs:   []goconfig.RefSpec{goconfig.RefSpec("+refs/tags/*:refs/tags/*")},
		RemoteName: TargetRemoteName,
		Auth:       &ssh2.PublicKeys{User: GitUser, Signer: config.TargetPrivateKey},
	})

	if err != git.NoErrAlreadyUpToDate {
		return err
	}
	Info.Printf("Pushed Tags for Group: %s, the Repo: %s", groupName, repoName)
	return nil

}

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "/config.json", "config file")
	flag.Parse()

	config := ReadConfigFile(configFile)
	ch = make(chan int, config.Concurrence)
	for i := 0; i < len(config.Groups); i++ {
		group := &config.Groups[i]
		for j := 0; j < len(group.Repos); j++ {
			repo := group.Repos[j]
			wg.Add(1)
			ch <- 1
			go func() {
				defer wg.Done()
				err := syncRepo(*config, group.Name, group.TargetName, repo)
				if err != nil {
					Error.Fatalln("Sync fail", err)
				}
				<-ch
			}()

		}
	}

	wg.Wait()
	Info.Println("Sync complete")

}
