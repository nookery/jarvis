package xcode

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "显示 Xcode 版本信息",
	Long:  color.Success.Render("\r\n显示当前系统中安装的 Xcode 版本信息，包括版本号、构建号和安装路径"),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		// 显示标题
		showXcodeInfoHeader()
		
		// 显示 Xcode 信息
		showXcodeInfo(verbose)
	},
}

func init() {
	infoCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

// showXcodeInfoHeader 显示信息标题
func showXcodeInfoHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         📱 Xcode 版本信息              ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showXcodeInfo 显示 Xcode 详细信息
func showXcodeInfo(verbose bool) {
	color.Blue.Println("🔍 检查 Xcode 安装")
	
	// 检查 Xcode 路径
	xcodePath := getCommandOutput("xcode-select", "-p")
	if xcodePath == "" {
		color.Error.Println("❌ 未找到 Xcode 安装")
		color.Info.Println("💡 请从 App Store 安装 Xcode")
		return
	}
	
	color.Success.Printf("✅ Xcode 路径: %s\n", xcodePath)
	fmt.Println()
	
	// 获取 Xcode 版本信息
	color.Blue.Println("📋 版本信息")
	xcodeVersion := getCommandOutput("xcodebuild", "-version")
	if xcodeVersion != "" {
		lines := strings.Split(xcodeVersion, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				if strings.HasPrefix(line, "Xcode") {
					color.Info.Printf("🚀 %s\n", line)
				} else if strings.HasPrefix(line, "Build version") {
					color.Info.Printf("🔨 %s\n", line)
				} else {
					color.Info.Printf("📝 %s\n", line)
				}
			}
		}
	} else {
		color.Error.Println("❌ 无法获取 Xcode 版本信息")
	}
	fmt.Println()
	
	// 显示 SDK 信息
	if verbose {
		showSDKInfo()
	}
	
	// 显示命令行工具信息
	showCommandLineToolsInfo(verbose)
	
	// 显示 Swift 版本
	showSwiftInfo()
	
	// 显示模拟器信息
	if verbose {
		showSimulatorInfo()
	}
}

// showSDKInfo 显示 SDK 信息
func showSDKInfo() {
	color.Blue.Println("📦 可用 SDK")
	
	// 获取 macOS SDK
	macosSDK := getCommandOutput("xcrun", "--show-sdk-path", "--sdk", "macosx")
	if macosSDK != "" {
		sdkVersion := getCommandOutput("xcrun", "--show-sdk-version", "--sdk", "macosx")
		color.Info.Printf("🖥️  macOS SDK: %s\n", sdkVersion)
		if strings.Contains(macosSDK, "/") {
			color.Gray.Printf("   路径: %s\n", macosSDK)
		}
	}
	
	// 获取 iOS SDK
	iosSDK := getCommandOutput("xcrun", "--show-sdk-path", "--sdk", "iphoneos")
	if iosSDK != "" {
		sdkVersion := getCommandOutput("xcrun", "--show-sdk-version", "--sdk", "iphoneos")
		color.Info.Printf("📱 iOS SDK: %s\n", sdkVersion)
		if strings.Contains(iosSDK, "/") {
			color.Gray.Printf("   路径: %s\n", iosSDK)
		}
	}
	
	// 获取 iOS 模拟器 SDK
	iosSimSDK := getCommandOutput("xcrun", "--show-sdk-path", "--sdk", "iphonesimulator")
	if iosSimSDK != "" {
		sdkVersion := getCommandOutput("xcrun", "--show-sdk-version", "--sdk", "iphonesimulator")
		color.Info.Printf("📲 iOS 模拟器 SDK: %s\n", sdkVersion)
		if strings.Contains(iosSimSDK, "/") {
			color.Gray.Printf("   路径: %s\n", iosSimSDK)
		}
	}
	
	fmt.Println()
}

// showCommandLineToolsInfo 显示命令行工具信息
func showCommandLineToolsInfo(verbose bool) {
	color.Blue.Println("🛠️  命令行工具")
	
	// 检查命令行工具版本
	clangVersion := getCommandOutput("clang", "--version")
	if clangVersion != "" {
		lines := strings.Split(clangVersion, "\n")
		if len(lines) > 0 {
			firstLine := strings.TrimSpace(lines[0])
			if strings.Contains(firstLine, "clang") {
				color.Info.Printf("🔧 %s\n", firstLine)
			}
		}
	}
	
	// 检查关键工具
	tools := []struct {
		name string
		desc string
	}{
		{"xcodebuild", "Xcode 构建工具"},
		{"xcrun", "Xcode 运行工具"},
		{"codesign", "代码签名工具"},
		{"security", "安全工具"},
		{"hdiutil", "磁盘映像工具"},
		{"plutil", "属性列表工具"},
		{"lipo", "架构工具"},
	}
	
	if verbose {
		for _, tool := range tools {
			if _, err := exec.LookPath(tool.name); err == nil {
				color.Success.Printf("✅ %s (%s)\n", tool.name, tool.desc)
			} else {
				color.Error.Printf("❌ %s (%s)\n", tool.name, tool.desc)
			}
		}
	} else {
		availableCount := 0
		for _, tool := range tools {
			if _, err := exec.LookPath(tool.name); err == nil {
				availableCount++
			}
		}
		color.Info.Printf("可用工具: %d/%d\n", availableCount, len(tools))
	}
	
	fmt.Println()
}

// showSwiftInfo 显示 Swift 信息
func showSwiftInfo() {
	color.Blue.Println("🚀 Swift 编译器")
	
	swiftVersion := getCommandOutput("swift", "--version")
	if swiftVersion != "" {
		lines := strings.Split(swiftVersion, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				if strings.Contains(line, "Swift version") {
					color.Info.Printf("⚡ %s\n", line)
				} else if strings.Contains(line, "Target:") {
					color.Info.Printf("🎯 %s\n", line)
				}
			}
		}
	} else {
		color.Error.Println("❌ 无法获取 Swift 版本信息")
	}
	
	fmt.Println()
}

// showSimulatorInfo 显示模拟器信息
func showSimulatorInfo() {
	color.Blue.Println("📲 iOS 模拟器")
	
	// 获取可用的模拟器
	simulators := getCommandOutput("xcrun", "simctl", "list", "devices", "available")
	if simulators != "" {
		lines := strings.Split(simulators, "\n")
		iosCount := 0
		for _, line := range lines {
			if strings.Contains(line, "iOS") && strings.Contains(line, "--") {
				iosCount++
				if iosCount <= 3 { // 只显示前3个版本
					line = strings.TrimSpace(line)
					color.Info.Printf("📱 %s\n", line)
				}
			}
		}
		if iosCount > 3 {
			color.Gray.Printf("   ... 还有 %d 个版本\n", iosCount-3)
		}
		if iosCount == 0 {
			color.Yellow.Println("⚠️  未找到可用的 iOS 模拟器")
		}
	} else {
		color.Error.Println("❌ 无法获取模拟器信息")
	}
	
	fmt.Println()
}