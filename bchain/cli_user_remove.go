package main

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var UserRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove an user",
	Long:  `remove an user`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.userRemove(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	UserCmd.AddCommand(UserRemoveCmd)
}

func (m *bchainCLI) userRemove(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Needs user name as first argument")
	}
	user := args[0]
	m.pInfo("Execute: Remove user %s\n", user)
	tapi := api.New(m.server)
	m.setAPI(tapi)
	if err := tapi.UserRemove(user); err != nil {
		return err
	}
	m.pSuccess("User removed %s\n", user)
	return nil
}
