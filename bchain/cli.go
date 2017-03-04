package main

//TODO:
//remove user code from nodeFunction

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use:   `bchain [OPTIONS] COMMAND [arg...]`,
		Short: "AntBlockchain storage cluster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}
)

func cli() {
	RootCmd.PersistentFlags().StringVar(&bCLI.server, "server", "127.0.0.1:30103", "Server addresses format addr1:port, addr2:port, ...")
	RootCmd.PersistentFlags().BoolVarP(&bCLI.verbose, "verbose", "v", false, `Verbose output`)
	RootCmd.PersistentFlags().BoolVarP(&bCLI.silence, "silence", "s", false, `Silence output`)
	RootCmd.PersistentFlags().BoolVar(&bCLI.debug, "debug", false, `Silence output`)
	RootCmd.PersistentFlags().StringVarP(&bCLI.userName, "user", "u", config.userName, `user name, default in config file`)
	RootCmd.PersistentFlags().StringVarP(&bCLI.keyPath, "key", "k", config.keyPath, `key path, default in config file`)
	cobra.OnInitialize(func() {
		if err := bCLI.init(); err != nil {
			fmt.Printf("Init error: %v\n", err)
			os.Exit(1)
		}
	})

	// versionCmd represents the antblockchain version
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version number of antblockchain",
		Long:  `Display the version number of antblockchain`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("antblockchain version: %s\n", Version)
		},
	}
	RootCmd.AddCommand(versionCmd)

	// infoCmd represents the antblockchain information
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display antblockchain version and server information",
		Long:  `Display antblockchain version and server information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("antblockchain version: %s\n", Version)
			fmt.Printf("Server: %s\n", config.serverAddress)
		},
	}
	RootCmd.AddCommand(infoCmd)

	//Execute commad
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error during: %s: %v\n", cmd.Name(), err)
		os.Exit(1)
	}

	os.Exit(0)
}
