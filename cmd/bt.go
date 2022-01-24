package cmd

import (
	"jarvis/cmd/bt"

	"github.com/spf13/cobra"
)

var BtCmd = &cobra.Command{
	Use:   "bt",
	Short: "宝塔相关操作",
}

func init() {
	rootCmd.AddCommand(BtCmd)
	BtCmd.AddCommand(bt.SiteCmd)
}
