package main

import (
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add entry or branch on blockchain",
	Long:  `add entry or branch on blockchain`,
	//Aliases: []string{"pf"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
