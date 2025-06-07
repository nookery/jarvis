package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "构建 iOS/macOS 应用",
	Long:  color.Success.Render("\r\n构建 iOS/macOS 应用程序，支持多架构和详细日志输出"),
	Run: func(cmd *cobra.Command, args []string) {
		scheme, _ := cmd.Flags().GetString("scheme")
		buildPath, _ := cmd.Flags().GetString("build-path")
		arch, _ := cmd.Flags().GetString("arch")
		verbose, _ := cmd.Flags().GetBool("verbose")
		clean, _ := cmd.Flags().GetBool("clean")
		
		// 显示环境信息
		showBuildEnvironmentInfo(scheme, buildPath, arch, verbose)
		
		// 检查必需的环境变量
		if scheme == "" {
			scheme = detectScheme()
			if scheme == "" {
				color.Error.Println("❌ 错误: 未设置 SCHEME 且无法自动检测")
				showAvailableSchemes()
				os.Exit(1)
			}
		}
		
		// 设置默认值
		if buildPath == "" {
			buildPath = "./temp"
		}
		if arch == "" {
			arch = "universal"
		}
		
		// 检测项目文件
		projectFile, projectType, err := detectProjectFile()
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(1)
		}
		
		// 显示构建目标信息
		showBuildTargetInfo(projectFile, projectType, scheme, arch)
		
		// 执行构建
		err = performBuild(projectFile, projectType, scheme, buildPath, arch, verbose, clean)
		if err != nil {
			color.Error.Printf("❌ 构建失败: %s\n", err.Error())
			os.Exit(1)
		}
		
		color.Success.Println("✅ 构建成功完成！")
		color.Green.Printf("📦 构建产物位置: %s/Build/Products/Release/\n", buildPath)
		
		// 显示开发路线图
		showDevelopmentRoadmap("build")
	},
}

func init() {
	buildCmd.Flags().StringP("scheme", "s", "", "构建方案名称")
	buildCmd.Flags().StringP("build-path", "b", "./temp", "构建输出路径")
	buildCmd.Flags().StringP("arch", "a", "universal", "目标架构 (universal, x86_64, arm64)")
	buildCmd.Flags().BoolP("verbose", "v", false, "显示详细构建日志")
	buildCmd.Flags().Bool("clean", true, "构建前清理")
}

// showBuildEnvironmentInfo 显示构建环境信息
func showBuildEnvironmentInfo(scheme, buildPath, arch string, verbose bool) {
	color.Blue.Println("===========================================")
	color.Blue.Println("         应用构建环境信息                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 系统信息
	color.Green.Println("📱 系统信息:")
	fmt.Printf("   操作系统: %s %s\n", runtime.GOOS, runtime.GOARCH)
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("   主机名称: %s\n", hostname)
	}
	fmt.Println()
	
	// Xcode 信息
	color.Green.Println("🔨 Xcode 开发环境:")
	if xcodeVersion := getCommandOutput("xcodebuild", "-version"); xcodeVersion != "" {
		lines := strings.Split(xcodeVersion, "\n")
		if len(lines) >= 1 {
			fmt.Printf("   Xcode 版本: %s\n", lines[0])
		}
		if len(lines) >= 2 {
			fmt.Printf("   构建版本: %s\n", lines[1])
		}
	}
	if sdkPath := getCommandOutput("xcrun", "--show-sdk-path"); sdkPath != "" {
		fmt.Printf("   SDK 路径: %s\n", sdkPath)
	}
	if devDir := getCommandOutput("xcode-select", "-p"); devDir != "" {
		fmt.Printf("   开发者目录: %s\n", devDir)
	}
	fmt.Println()
	
	// Swift 信息
	color.Green.Println("🚀 Swift 编译器:")
	if swiftVersion := getCommandOutput("swift", "--version"); swiftVersion != "" {
		lines := strings.Split(swiftVersion, "\n")
		if len(lines) >= 1 {
			fmt.Printf("   Swift 版本: %s\n", lines[0])
		}
	}
	fmt.Println()
	
	// Git 信息
	color.Green.Println("📝 Git 版本控制:")
	if gitVersion := getCommandOutput("git", "--version"); gitVersion != "" {
		fmt.Printf("   Git 版本: %s\n", gitVersion)
	}
	if branch := getCommandOutput("git", "branch", "--show-current"); branch != "" {
		fmt.Printf("   当前分支: %s\n", branch)
	}
	if commit := getCommandOutput("git", "log", "-1", "--pretty=format:%h - %s (%an, %ar)"); commit != "" {
		fmt.Printf("   最新提交: %s\n", commit)
	}
	fmt.Println()
	
	// 构建环境变量
	color.Green.Println("🌍 构建环境变量:")
	fmt.Printf("   构建方案: %s\n", scheme)
	fmt.Printf("   构建路径: %s\n", buildPath)
	fmt.Printf("   目标架构: %s\n", arch)
	fmt.Printf("   构建配置: Release\n")
	fmt.Printf("   详细日志: %t\n", verbose)
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("   工作目录: %s\n", cwd)
	}
	fmt.Println()
}

// detectScheme 自动检测可用的 scheme
func detectScheme() string {
	projectFile, projectType, err := detectProjectFile()
	if err != nil {
		return ""
	}
	
	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-list")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-list")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	// 解析输出，查找 schemes
	lines := strings.Split(string(output), "\n")
	inSchemes := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Schemes:" {
			inSchemes = true
			continue
		}
		if inSchemes && line != "" && !strings.Contains(line, ":") {
			return line // 返回第一个找到的 scheme
		}
		if inSchemes && line == "" {
			break
		}
	}
	
	return ""
}

// showAvailableSchemes 显示可用的 schemes
func showAvailableSchemes() {
	color.Yellow.Println("正在检查项目中可用的 scheme...")
	
	projectFile, projectType, err := detectProjectFile()
	if err != nil {
		color.Error.Printf("❌ %s\n", err.Error())
		return
	}
	
	color.Green.Printf("在项目 %s 中找到以下可用的 scheme:\n", projectFile)
	
	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-list")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-list")
	}
	
	output, err := cmd.Output()
	if err != nil {
		color.Error.Printf("❌ 无法获取 scheme 列表: %v\n", err)
		return
	}
	
	// 解析并显示 schemes
	lines := strings.Split(string(output), "\n")
	inSchemes := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Schemes:" {
			inSchemes = true
			continue
		}
		if inSchemes && line != "" && !strings.Contains(line, ":") {
			fmt.Printf("  - %s\n", line)
		}
		if inSchemes && line == "" {
			break
		}
	}
	
	fmt.Println()
	color.Yellow.Println("💡 使用示例:")
	color.Cyan.Println("   go run main.go xcode build --scheme YourSchemeName")
}

// detectProjectFile 检测项目文件
func detectProjectFile() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("无法获取当前目录: %v", err)
	}
	
	// 查找 .xcworkspace 文件
	workspaces, err := filepath.Glob(filepath.Join(cwd, "*.xcworkspace"))
	if err == nil && len(workspaces) > 0 {
		return workspaces[0], "workspace", nil
	}
	
	// 查找 .xcodeproj 文件
	projects, err := filepath.Glob(filepath.Join(cwd, "*.xcodeproj"))
	if err == nil && len(projects) > 0 {
		return projects[0], "project", nil
	}
	
	return "", "", fmt.Errorf("未找到 .xcodeproj 或 .xcworkspace 文件")
}

// showBuildTargetInfo 显示构建目标信息
func showBuildTargetInfo(projectFile, projectType, scheme, arch string) {
	color.Green.Println("🎯 构建目标信息:")
	fmt.Printf("   项目文件: %s\n", projectFile)
	if projectType == "workspace" {
		fmt.Println("   项目类型: Xcode Workspace")
	} else {
		fmt.Println("   项目类型: Xcode Project")
	}
	fmt.Printf("   构建方案: %s\n", scheme)
	
	// 显示支持的架构
	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-scheme", scheme, "-showBuildSettings", "-configuration", "Release")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-scheme", scheme, "-showBuildSettings", "-configuration", "Release")
	}
	
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "ARCHS =") {
				parts := strings.Split(line, "=")
				if len(parts) >= 2 {
					projectArchs := strings.TrimSpace(parts[1])
					fmt.Printf("   项目支持架构: %s\n", projectArchs)
					break
				}
			}
		}
	}
	
	fmt.Printf("   构建目标架构: %s\n", arch)
	fmt.Println()
}

// performBuild 执行构建
func performBuild(projectFile, projectType, scheme, buildPath, arch string, verbose, clean bool) error {
	color.Blue.Println("===========================================")
	color.Yellow.Println("🚀 开始构建过程...")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 构建基础参数
	args := []string{}
	if projectType == "workspace" {
		args = append(args, "-workspace", projectFile)
	} else {
		args = append(args, "-project", projectFile)
	}
	
	args = append(args, "-scheme", scheme, "-configuration", "Release", "-derivedDataPath", buildPath)
	
	// 设置目标和架构
	args = append(args, "-destination", "generic/platform=macOS")
	if arch != "universal" {
		args = append(args, "ARCHS="+arch, "ONLY_ACTIVE_ARCH=NO")
	} else {
		args = append(args, "ARCHS=x86_64 arm64", "ONLY_ACTIVE_ARCH=NO")
	}
	
	// 添加静默参数
	if !verbose {
		args = append(args, "-quiet")
	}
	
	// 清理构建
	if clean {
		color.Yellow.Println("正在清理之前的构建...")
		cleanArgs := append(args, "clean")
		cleanCmd := exec.Command("xcodebuild", cleanArgs...)
		if verbose {
			cleanCmd.Stdout = os.Stdout
			cleanCmd.Stderr = os.Stderr
		}
		err := cleanCmd.Run()
		if err != nil {
			return fmt.Errorf("清理失败: %v", err)
		}
	}
	
	// 开始构建
	if arch == "universal" {
		color.Yellow.Println("开始构建应用 (通用二进制: x86_64 arm64)...")
	} else {
		color.Yellow.Printf("开始构建应用 (架构: %s)...\n", arch)
	}
	
	buildArgs := append(args, "build")
	buildCmd := exec.Command("xcodebuild", buildArgs...)
	
	if verbose {
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		fmt.Printf("执行命令: xcodebuild %s\n", strings.Join(buildArgs, " "))
		fmt.Println()
	}
	
	err := buildCmd.Run()
	if err != nil {
		return fmt.Errorf("构建失败: %v", err)
	}
	
	return nil
}