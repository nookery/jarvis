package site

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var SiteCmd = &cobra.Command{
	Use:   "site",
	Short: color.Blue.Render("网站相关操作"),
	Long:  color.Success.Render("\r\n网站相关操作"),
}

func init() {
	SiteCmd.AddCommand(show)
	SiteCmd.AddCommand(types)
	SiteCmd.AddCommand(php)
	SiteCmd.AddCommand(delete)
	SiteCmd.AddCommand(Create)
	SiteCmd.AddCommand(Conf)
}
