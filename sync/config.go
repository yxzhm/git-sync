package main

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

type Config struct {
	SourceURL        string `json:"sourceUrl"`
	SourceKey        string `json:"sourceKey"`
	SourcePrivateKey ssh.Signer
	TargetURL        string `json:"targetUrl"`
	TargetKey        string `json:"targetKey"`
	TargetPrivateKey ssh.Signer
	Groups           []Group `json:"groups"`
}

type Group struct {
	Name  string   `json:"name"`
	Repos []string `json:"repos"`
}

func ReadConfigFile(configFile string) *Config {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		Error.Fatalln("Failed to open config file", err)
	}
	defer jsonFile.Close()

	var config *Config
	err = json.NewDecoder(jsonFile).Decode(&config)
	if err != nil {
		Error.Fatalln("failed to decode config file", err)
	}

	if config.SourceURL == config.TargetURL {
		Error.Fatalln("The Source URL is same with the Target URL")
	}

	sourceSSHKey, err := ioutil.ReadFile(config.SourceKey)
	if err != nil {
		Error.Fatalln("Can't read Source SSH Key")
	}
	sourceSinger, err := ssh.ParsePrivateKey(sourceSSHKey)
	if err != nil {
		Error.Fatalln("Invalid Source SSH Key")
	}
	config.SourcePrivateKey = sourceSinger

	targetSSHKey, err := ioutil.ReadFile(config.TargetKey)
	if err != nil {
		Error.Fatalln("Can't read Target SSH Key")
	}
	targetSinger, err := ssh.ParsePrivateKey(targetSSHKey)
	if err != nil {
		Error.Fatalln("Invalid Target SSH Key")
	}
	config.TargetPrivateKey = targetSinger

	return config
}
