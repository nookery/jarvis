package bt

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var SiteCmd = &cobra.Command{
	Use:   "site",
	Short: "网站相关操作",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println("pang")
	},
}
