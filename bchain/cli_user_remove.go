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
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	UserCmd.AddCommand(UserRemoveCmd)
	UserRemoveCmd.Flags().Bool("force", false, `WARNING: force to removce user with its associated files`)
}

func (m *bchainCLI) userRemove(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Error number of argument, needs [user] format userName:token")
	}
	user := args[0]
	force := false
	if cmd.Flag("force").Value.String() == "true" {
		force = true
	}
	m.pInfo("Execute: Remove user %s\n", user)
	tapi := api.New(m.server)
	m.setAPILogLevel(tapi)
	if err := tapi.UserRemove(user, force); err != nil {
		return err
	}
	m.pSuccess("User removed %s\n", user)
	return nil
}
