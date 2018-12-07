package util

import (
	"encoding/json"
	"os"
)

type Config struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Groups []Group `json:"groups"`
}

type Group struct {
	Name  string   `json:"name"`
	Repos []string `json:"repo"`
}

func ReadConfigFile(confilgFile string) *Config {
	jsonFile, err := os.Open(confilgFile)
	if err != nil {
		Error.Fatalln("Failed to open config file", err)
	}
	defer jsonFile.Close()

	var config *Config
	err = json.NewDecoder(jsonFile).Decode(&config)
	if err != nil {
		Error.Fatalln("failed to decode config file", err)
	}

	return config
}
