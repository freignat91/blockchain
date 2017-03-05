package main

import (
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeWriteStatsCmd = &cobra.Command{
	Use:   "writeStats",
	Short: "write stats in log file",
	Long:  `write stats in log file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.NodeWriteStats(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodeWriteStatsCmd)
	NodeWriteStatsCmd.Flags().StringP("node", "n", "*", "Target a specific node")
}

func (m *bchainCLI) NodeWriteStats(cmd *cobra.Command, args []string) error {
	node := "*"
	if len(args) >= 1 {
		node = args[0]
	}
	m.pInfo("Execute: writeStats\n")
	tapi := api.New(m.server)
	m.setAPI(tapi)
	if err := tapi.InfoWriteStats(node); err != nil {
		return err
	}
	m.pSuccess("Stats written for node %s\n", node)
	return nil
}
