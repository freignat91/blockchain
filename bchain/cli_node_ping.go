package main

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
	"time"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodePingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping  an blockchain node",
	Long:  `ping an blockchain node`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.ping(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodePingCmd)
}

func (m *bchainCLI) ping(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Needs node name as first argument")
	}
	node := args[0]
	m.pInfo("Execute: ping %s\n", node)
	t0 := time.Now()
	tapi := api.New(m.server)
	m.setAPILogLevel(tapi)
	path, err := tapi.NodePing(node, false)
	if err != nil {
		return err
	}
	t1 := time.Now()
	m.pSuccess("Ping time=%dms path: %s\n", t1.Sub(t0).Nanoseconds()/1000000, path)
	return nil
}
