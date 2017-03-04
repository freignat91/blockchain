package main

import (
	"fmt"

	"github.com/freignat91/blockchain/api"
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var AddBranchCmd = &cobra.Command{
	Use:   "branch [branchLabel...]",
	Short: "add a branch in the blockchain using labels, label format: name=value",
	Long:  `add a branch in the blockchain using labels, label format: name=value`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.addBranch(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	AddCmd.AddCommand(AddBranchCmd)
}

func (m *bchainCLI) addBranch(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("addEntry needs at least one argument")
	}
	m.pInfo("Execute: add entry\n")
	tapi := api.New(m.server)
	if err := m.setAPI(tapi); err != nil {
		return err
	}
	if err := tapi.AddBranch(args); err != nil {
		return err
	}
	m.pSuccess("branch sent\n")
	return nil
}
