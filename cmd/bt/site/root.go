package site

import (
	"github.com/spf13/cobra"
)

var SiteCmd = &cobra.Command{
	Use:   "site",
	Short: "网站相关操作",
}

func init() {
	SiteCmd.AddCommand(show)
	SiteCmd.AddCommand(types)
	SiteCmd.AddCommand(php)
	SiteCmd.AddCommand(delete)
	SiteCmd.AddCommand(Create)
}
