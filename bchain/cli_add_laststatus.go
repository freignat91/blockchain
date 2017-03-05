package main

import (
	"fmt"
	"strconv"

	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var LastAddStatusCmd = &cobra.Command{
	Use:   "last [number]",
	Short: "return the status of last requests",
	Long:  `return the status of last requests`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.lastAddStatus(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	AddCmd.AddCommand(LastAddStatusCmd)
	LastAddStatusCmd.Flags().String("user", "", `get the status belonging to user only`)
	LastAddStatusCmd.Flags().Bool("error-only", false, `get the status on error only`)
}

func (m *bchainCLI) lastAddStatus(cmd *cobra.Command, args []string) error {
	nb := 10
	if len(args) > 0 {
		nn, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("argument shoud be a number")
		}
		nb = nn
	}
	m.fullColor = true
	m.pInfo("Execute: display last add request status\n")
	tapi := api.New(m.server)
	if err := m.setAPI(tapi); err != nil {
		return err
	}
	userName := cmd.Flag("user").Value.String()
	errorOnly := false
	if cmd.Flag("error-only").Value.String() == "true" {
		errorOnly = true
	}
	if err := tapi.LastAddRequestStatus(nb, userName, errorOnly, m.displayStatus); err != nil {
		return err
	}
	return nil
}
