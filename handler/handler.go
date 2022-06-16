package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/melbahja/goph"
	"gopkg.in/yaml.v2"
)

type commands struct {
	Commands []RemoteCommand `yaml:"commands"`
}

type RemoteCommand struct {
	Name        string `yaml:"name"`
	ExecCommand string `yaml:"command"`
	KeyPath     string `yaml:"key-path"`
	IP          string `yaml:"ip"`
	Username    string `yaml:"username"`
}

func ConfigOption(configPath string) ([]RemoteCommand, error) {
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

func Command(cmd RemoteCommand) error {
	auth, err := goph.Key(cmd.KeyPath, "")
	if err != nil {
		return err
	}
	client, err := goph.New(cmd.Username, cmd.IP, auth)
	if err != nil {
		return err
	}
	defer client.Close()
	cliComand := fmt.Sprintf("bash -l -c '%s'", cmd.ExecCommand)
	cmdExec, err := client.Command(cliComand)
	if err != nil {
		return err
	}

	stdout, err := cmdExec.StdoutPipe()
	stderr, err := cmdExec.StderrPipe()
	err = cmdExec.Start()
	if err != nil {
		log.Println(err)
	}
	defer cmdExec.Wait()
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)
	return nil
}
