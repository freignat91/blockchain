package main

import (
	"fmt"

	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var AddEntryCmd = &cobra.Command{
	Use:   "entry [entry] [branchLabel...]",
	Short: "add an entry in the blockchain at the end of the branch carring all the branch labels, label format: name=value",
	Long:  `add an entry in the blockchain at the end of the branch carring all the branch labels, label format: name=value`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.addEntry(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	AddCmd.AddCommand(AddEntryCmd)
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
