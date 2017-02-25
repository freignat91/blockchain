package main

import (
	"fmt"

	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add an entry in the blockchain",
	Long:  `add an entry in the blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.addEntry(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}

func (m *bchainCLI) addEntry(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("addEntry needs at least one argument")
	}
	m.pInfo("Execute: add entry\n")
	tapi := api.New(m.server)
	if err := m.setAPI(tapi); err != nil {
		return err
	}
	if err := tapi.AddEntry([]byte(args[0]), args[1:]); err != nil {
		return err
	}
	m.pSuccess("entry sent\n")
	return nil
}
