package main

import (
	"fmt"
	"github.com/freignat91/agrid/agridapi"
	"github.com/spf13/cobra"
	"time"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var NodePingFromToCmd = &cobra.Command{
	Use:   "pingFromTo",
	Short: "pingFromTo",
	Long:  `pingFromTo`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.nodePingFromTo(cmd, args); err != nil {
			bCLI.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(NodePingFromToCmd)
}

func (m *bchainCLI) nodePingFromTo(cmd *cobra.Command, args []string) error {
	m.pInfo("Execute: pingFromTo\n")
	if len(args) < 2 {
		return fmt.Errorf("Needs two arguements: sender node name and targeted node name")
	}
	t0 := time.Now()
	api := agridapi.New(m.server)
	m.setAPILogLevel(api)
	path, err := api.NodePingFromTo(args[0], args[1], m.debug)
	if err != nil {
		return err
	}
	t1 := time.Now()
	m.pSuccess("Ping %s -> %s time=%d ms path: %v\n", args[0], args[1], t1.Sub(t0).Nanoseconds()/1000000, path)
	return nil
}
