/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/PieroNarciso/comander/handler"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "comander",
	Short: "Application for sending commands to a remote machine",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		cmds, err := handler.ConfigOption(path)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		var wg sync.WaitGroup

		for ix := range cmds {
			wg.Add(1)
			go func(command handler.RemoteCommand, wg *sync.WaitGroup) {
				defer wg.Done()
				err = handler.Command(command)
				if err != nil {
					log.Println(err)
				}
			}(cmds[ix], &wg)
		}
		wg.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.comander.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("config", "c", "./config.yaml", "Path to config file in `.yml`")
}
