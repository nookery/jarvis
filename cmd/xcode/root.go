package xcode

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var XcodeCmd = &cobra.Command{
	Use:   "xcode",
	Short: color.Blue.Render("Xcode 开发工具集"),
	Long:  color.Success.Render("\r\nXcode 开发工具集，包含构建、签名、打包等功能"),
}

func init() {
	XcodeCmd.AddCommand(versionCmd)
	XcodeCmd.AddCommand(bumpCmd)
	XcodeCmd.AddCommand(buildCmd)
	XcodeCmd.AddCommand(codesignCmd)
	XcodeCmd.AddCommand(packageCmd)
	XcodeCmd.AddCommand(setupCmd)
	XcodeCmd.AddCommand(infoCmd)
}