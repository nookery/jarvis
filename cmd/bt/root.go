package bt

import (
	"jarvis/cmd/bt/site"

	"github.com/spf13/cobra"
)

var BtCmd = &cobra.Command{
	Use:   "bt",
	Short: "宝塔相关操作",
}

func init() {
	BtCmd.AddCommand(site.SiteCmd)
}
