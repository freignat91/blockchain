package main

import (
	"fmt"
	"time"

	"github.com/freignat91/blockchain/api"
	"github.com/freignat91/blockchain/server/gnode"
	"github.com/spf13/cobra"
)

type treeDisplayer struct {
	isBlocks   bool
	isEntries  bool
	debug      bool
	hash       bool
	fullHashId bool
	cli        *bchainCLI
}

// PlatformMonitor is the main command for attaching platform subcommands.
var DisplayCmd = &cobra.Command{
	Use:   "display",
	Short: "display the blockchain tree starting from the branch matchign the labels",
	Long:  `display the blockchain tree starting from the branch matchign the labels`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bCLI.displayTree(cmd, args); err != nil {
			bCLI.fatal("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(DisplayCmd)
	DisplayCmd.Flags().Bool("blocks", false, `display branches and block`)
	DisplayCmd.Flags().Bool("entries", false, `display branches and block and entries`)
	DisplayCmd.Flags().Bool("debug", false, `display branches and block and entries and debug information`)
	DisplayCmd.Flags().Bool("hash", false, `display hash branches/blocks instead of id`)
	DisplayCmd.Flags().Bool("full-hash-id", false, `display the entire hash and id`)
}

//For debug
func (m *bchainCLI) displayTree(cmd *cobra.Command, args []string) error {
	m.pInfo("Display\n")
	d := treeDisplayer{cli: m}
	if cmd.Flag("blocks").Value.String() == "true" {
		d.isBlocks = true
	}
	if cmd.Flag("entries").Value.String() == "true" {
		d.isEntries = true
		d.isBlocks = true
	}
	if cmd.Flag("debug").Value.String() == "true" {
		d.debug = true
	}
	if cmd.Flag("full-hash-id").Value.String() == "true" {
		d.fullHashId = true
	}
	if cmd.Flag("hash").Value.String() == "true" {
		d.hash = true
	}
	tapi := api.New(m.server)
	m.fullColor = true
	m.setAPI(tapi)
	err := tapi.GetTree(args, d.isBlocks, d.isEntries, d.displayBlock)
	if err != nil {
		return err
	}
	return nil
}

func (d *treeDisplayer) displayBlock(id string, blockType string, block *gnode.TreeBlock) error {
	if blockType == "branch" || d.isBlocks {
		tab := ""
		for i := 0; i < int(block.Depth); i++ {
			tab += "  "
		}
		id := fmt.Sprintf("id=%s", d.shorters(block.Id))
		if d.hash {
			id = fmt.Sprintf("hash=%x", d.shorterb(block.FullHash))
		}
		if blockType == "branch" {
			d.cli.pSuccess("%s%s: %s nb=%d branchLabel=%s:%s\n", tab, blockType, id, block.Size, block.LabelName, block.LabelValue)
		} else {
			tab += "  "
			d.cli.pRegular("%s%s: %s nb=%d\n", tab, blockType, id, block.Size)
		}
		if d.debug {
			d.cli.pWarn("%sdebug:[parentId=%s childId=%s loaded=%t updated=%t]\n", tab, d.shorters(block.ParentId), d.shorters(block.ChildId), block.Loaded, block.Updated)
		}
		if d.isEntries {
			for _, entry := range block.Entries {
				d.cli.pInfo("%s->entry: %s date=%s user=%s labels=[%s]\n", tab, string(entry.Payload), d.getDate(entry), entry.UserName, d.getLabels(entry))
			}
		}
	}
	return nil
}

func (d *treeDisplayer) getDate(entry *gnode.BCEntry) string {
	dt := time.Unix(entry.Date, 0)
	return dt.Format("2006-01-02T15:04:05")
}

func (d *treeDisplayer) getLabels(entry *gnode.BCEntry) string {
	labels := ""
	for _, label := range entry.Labels {
		labels += fmt.Sprintf("%s:%s ", label.Name, label.Value)
	}
	return labels
}

func (d *treeDisplayer) shorters(val string) string {
	if d.fullHashId {
		return val
	}
	if len(val) > 15 {
		return val[0:15]
	}
	return val
}

func (d *treeDisplayer) shorterb(valb []byte) string {
	val := fmt.Sprintf("%x", valb)
	return d.shorters(val)
}
