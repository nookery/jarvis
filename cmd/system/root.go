package system

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: color.Blue.Render("系统信息查看工具"),
	Long:  color.Success.Render("\r\n系统信息查看工具，包含操作系统基础信息、资源占用情况、进程信息等"),
}

func init() {
	SystemCmd.AddCommand(infoCmd)
	SystemCmd.AddCommand(resourceCmd)
	SystemCmd.AddCommand(processCmd)
	SystemCmd.AddCommand(networkCmd)
	SystemCmd.AddCommand(diskCmd)
}