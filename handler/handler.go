package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/melbahja/goph"
	"gopkg.in/yaml.v2"
)

type commands struct {
	Commands []remoteCommand `yaml:"commands"`
}

type remoteCommand struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	ExecCommand string `yaml:"command"`
	KeyPath     string `yaml:"key-path"`
	IP          string `yaml:"ip"`
	Username    string `yaml:"username"`
}

func ConfigOption(configPath string) ([]remoteCommand, error) {
	absPath, _ := filepath.Abs(configPath)
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	cmds := &commands{}
	err = yaml.Unmarshal(file, cmds)
	if err != nil {
		return nil, err
	}
	return cmds.Commands, nil
}

func Command(cmd remoteCommand) error {
	auth, err := goph.Key(cmd.KeyPath, "")
	if err != nil {
		return err
	}
	client, err := goph.New(cmd.Username, cmd.IP, auth)
	if err != nil {
		return err
	}
	defer client.Close()
	// Move to path
	strCommand := fmt.Sprintf("cd %s && %s", cmd.Path, cmd.ExecCommand)
	output, err := client.Run(strCommand)
	if err != nil {
		return err
	}
	log.Println(string(output))
	log.Println(cmd.Name, "command run correctly")
	return nil
}
