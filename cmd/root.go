package cmd

import (
	"jarvis/cmd/bt"
	"jarvis/cmd/database"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "jarvis",
	Long: color.Success.Render("\r\n我是Jarvis，你的得力助理。"),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(database.DatabaseCmd)
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(jokeCmd)
	rootCmd.AddCommand(bt.BtCmd)

	// 关闭completion命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 自定义Usage提示，关键字染色
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.Replace(usageTemplate, "Usage:", color.Yellow.Render("Usage:"), 1)
	usageTemplate = strings.Replace(usageTemplate, "Available Commands:", color.Yellow.Render("Available Commands:"), 1)
	usageTemplate = strings.Replace(usageTemplate, "Flags:", color.Yellow.Render("Flags:"), 1)
	usageTemplate = strings.Replace(usageTemplate, "Global Flags:", color.Yellow.Render("Global Flags:"), 1)
	usageTemplate = strings.Replace(usageTemplate, "Use \"{{.CommandPath}} [command] --help\" for more information about a command.", "", 1)
	rootCmd.SetUsageTemplate(usageTemplate)

	// 自定义Help提示
	// helpTemplate := rootCmd.HelpTemplate()
	// helpTemplate = strings.Replace(usageTemplate, "Usage", color.Yellow.Render("Usage"), 1)
	// rootCmd.SetHelpTemplate(helpTemplate)
}
