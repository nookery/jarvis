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

	SiteCmd.PersistentFlags().String("host", "http://127.0.0.1", "宝塔地址")
	SiteCmd.PersistentFlags().String("key", "", "宝塔密钥")
}
