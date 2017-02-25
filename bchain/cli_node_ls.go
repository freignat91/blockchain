package main

import (
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodeLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "get the blockchain nodes list",
	Long:  `get the blockchain nodes list`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.getNodeList(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodeLsCmd)
}

func (m *bchainCLI) getNodeList(cmd *cobra.Command, args []string) error {
	m.pInfo("Execute: getNodeList\n")

	tapi := api.New(m.server)
	m.setAPI(tapi)
	list, err := tapi.NodeLs()
	if err != nil {
		return err
	}
	for _, line := range list {
		m.pSuccess("%s\n", line)
	}
	m.pSuccess("number=%d\n", len(list))
	return nil
}
