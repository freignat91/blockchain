package main

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var UserCreateCmd = &cobra.Command{
	Use:   "signup",
	Short: "sigup in  antblockchain",
	Long:  `sigup in antblockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.userSignup(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	UserCmd.AddCommand(UserCreateCmd)
	UserCreateCmd.Flags().String("key-path", "./private.key", "path to copy the private key")
}

func (m *bchainCLI) userSignup(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Needs user name as first argument")
	}
	keyPath := cmd.Flag("key-path").Value.String()
	user := args[0]
	m.pInfo("Execute: Signup user %s\n", user)
	tapi := api.New(m.server)
	m.setAPI(tapi)
	if err := tapi.UserSignup(user, keyPath); err != nil {
		return err
	}
	m.pSuccess("User %s created\n", user)
	m.pSuccess("Private key file create: %s\n", keyPath)
	m.pWarn("Keep this private on a secure place, it'll be mandatory for every request on the cluster\n")
	return nil
}
