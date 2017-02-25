package main

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill an blockchain node",
	Long:  `kill an blockchain node`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.kill(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodeKillCmd)
}

func (m *bchainCLI) kill(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Needs node name as first argument")
	}
	node := args[0]
	m.pInfo("Execute: kill node %s\n", node)
	tapi := api.New(m.server)
	m.setAPILogLevel(tapi)
	if err := tapi.NodeKill(node); err != nil {
		return err
	}
	m.pSuccess("Container killed node %s\n", node)
	return nil
}
