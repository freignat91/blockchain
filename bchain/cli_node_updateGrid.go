package main

import (
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeUpdateGridCmd = &cobra.Command{
	Use:   "updateGrid",
	Short: "update grid connections",
	Long:  `update grid connections`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.nodeUpdateGrid(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeUpdateGridCmd.Flags().StringP("node", "n", "*", "Target a specific node")
	NodeUpdateGridCmd.Flags().BoolP("force", "f", false, "Force to recompute")
	NodeCmd.AddCommand(NodeUpdateGridCmd)
}

func (m *bchainCLI) nodeUpdateGrid(cmd *cobra.Command, args []string) error {
	node := "*"
	if len(args) >= 1 {
		node = args[0]
	}
	force := false
	if cmd.Flag("force").Value.String() == "true" {
		force = true
	}

	if len(args) >= 1 {
		node = args[0]
	}
	m.pInfo("Execute: update grid\n")
	tapi := api.New(m.server)
	m.setAPILogLevel(tapi)
	if err := tapi.NodeUpdateGrid(node, force); err != nil {
		return err
	}
	m.pSuccess("Grid updated for node %s\n", node)
	return nil
}
