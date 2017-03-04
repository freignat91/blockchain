package main

import (
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeConnectionCmd = &cobra.Command{
	Use:   "connections",
	Short: "get the blockchain nodes list with their connections",
	Long:  `get the blockchain nodes list with their connections`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.getNodeConnections(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodeConnectionCmd)
}

func (m *bchainCLI) getNodeConnections(cmd *cobra.Command, args []string) error {
	m.pInfo("Execute: getNodeConnections\n")

	tapi := api.New(m.server)
	m.setAPI(tapi)
	list, err := tapi.NodeConnections()
	if err != nil {
		return err
	}
	for _, line := range list {
		m.pSuccess("%s\n", line)
	}
	m.pSuccess("number=%d\n", len(list))
	return nil
}
