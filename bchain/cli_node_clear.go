package main

import (
	"github.com/freignat91/agrid/agridapi"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear an agrid node",
	Long:  `clear an agrid node`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.clear(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodeClearCmd)
}

func (m *bchainCLI) clear(cmd *cobra.Command, args []string) error {
	node := "*"
	if len(args) >= 1 {
		node = args[0]
	}
	m.pInfo("Execute: clear node %s\n", node)
	api := agridapi.New(m.server)
	m.setAPILogLevel(api)
	if err := api.NodeClear(node); err != nil {
		return err
	}
	m.pSuccess("done\n")
	return nil
}
