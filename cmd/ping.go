package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "输出pang",
	Long:  `输出pang，用于测试。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pang")
	},
}
