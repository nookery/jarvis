package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: color.Blue.Render("输出pang"),
	Long:  color.Success.Render(`输出pang，用于测试。`),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pang")
	},
}
