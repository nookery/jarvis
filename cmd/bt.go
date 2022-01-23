package cmd

import (
	"github.com/spf13/cobra"
	"jarvis/cmd/bt"
)

var BtCmd = &cobra.Command{
	Use:   "bt",
	Short: "宝塔相关操作",
}

func init() {
	rootCmd.AddCommand(BtCmd)
	BtCmd.AddCommand(bt.SiteCmd)
}
