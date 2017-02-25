package main

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var UserCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create an user",
	Long:  `create an user`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.userCreate(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	UserCmd.AddCommand(UserCreateCmd)
	UserCreateCmd.Flags().String("token", "", `force user token`)
}

func (m *bchainCLI) userCreate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Needs user name as first argument")
	}
	token := cmd.Flag("token").Value.String()
	user := args[0]
	m.pInfo("Execute: Create user %s\n", user)
	tapi := api.New(m.server)
	m.setAPILogLevel(tapi)
	token, err := tapi.UserCreate(user, token)
	if err != nil {
		return err
	}
	m.pSuccess("User create: %s:%s\n", user, token)
	return nil
}
