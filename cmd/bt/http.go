package bt

import (
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var Http = &cobra.Command{
	Use:   "http",
	Short: color.Blue.Render("发送HTTP请求"),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")

		query, _ := cmd.Flags().GetString("query")
		data, _ := cmd.Flags().GetString("data")
		link := host + query

		parsedData, _ := url.ParseQuery(data)
		result := utils.Post(link, utils.PatchSign(key, parsedData))
		color.Infoln(result)
	},
}

func init() {
	Http.Flags().String("query", "", color.Blue.Render("URL查询参数，如：/plugin?action=a&name=supervisor&s=AddProcess"))
	Http.Flags().String("data", "", color.Blue.Render("发送的数据，如：pjname=abcd&user=www&path=/&command=tail -f /dev/null&numprocs=1"))
	Http.MarkFlagRequired("query")
}
