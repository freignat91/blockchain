package main

import (
	"fmt"

	"github.com/freignat91/blockchain/api"
	"github.com/freignat91/blockchain/server/gnode"
	"github.com/spf13/cobra"
	"time"
)

// PlatformCmd is the main command for attaching topic subcommands.
var AddStatusCmd = &cobra.Command{
	Use:   "status [request id]",
	Short: "return the status of the request",
	Long:  `return the status of the request`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.addStatus(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	AddCmd.AddCommand(AddStatusCmd)
}

func (m *bchainCLI) addStatus(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("needs the request id")
	}
	m.pInfo("Execute: display add request status\n")
	tapi := api.New(m.server)
	if err := m.setAPI(tapi); err != nil {
		return err
	}
	status, erra := tapi.AddRequestStatus(args[0])
	if erra != nil {
		return erra
	}
	m.fullColor = true
	m.displayStatus(status)
	return nil
}

func (m *bchainCLI) displayStatus(status *gnode.RequestStatus) error {
	dt := time.Unix(status.Date, 0)
	sdate := dt.Format("2006-01-02 15:04:05")
	if status.Status == gnode.ReqStatusOnError {
		m.pError("%s %s %s %s\n", sdate, status.Id, status.Status, status.ReqError)
	} else if status.Status == gnode.ReqStatusOnGoing {
		m.pInfo("%s %s %s %s nbOk=%d nbSame=%d done=%t\n", sdate, status.UserName, status.Id, status.Status, status.NbOk, status.NbSame, status.Done)
	} else {
		m.pSuccess("%s %s %s %s %s\n", sdate, status.UserName, status.Id, status.Status, status.ReqError)
	}
	return nil
}
