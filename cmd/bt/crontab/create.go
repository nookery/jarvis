package crontab

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"jarvis/cmd/bt/utils"
	"net/url"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "创建crontab",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		name, _ := cmd.Flags().GetString("name")
		shell, _ := cmd.Flags().GetString("shell")
		link := host + "/crontab?action=AddCrontab"

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"name":           {name},
			"type":           {"minute-n"},
			"sType":          {"toShell"},
			"sBody":          {shell},
			"where1":         {"1"},
			"hour":           {""},
			"minute":         {""},
			"week":           {""},
			"sName":          {""},
			"backupTo":       {""},
			"save":           {""},
			"urladdress":     {""},
			"save_local":     {"1"},
			"notice":         {""},
			"notice_channel": {""},
		}))
		color.Infoln(result)
	},
}

func init() {
	create.Flags().String("name", "", color.Blue.Render("名称"))
	create.Flags().String("shell", "", color.Blue.Render("脚本"))
	create.MarkFlagRequired("name")
	create.MarkFlagRequired("shell")
}
