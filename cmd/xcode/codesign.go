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

var codesignCmd = &cobra.Command{
	Use:   "codesign",
	Short: "对 macOS 应用进行代码签名",
	Long:  color.Success.Render("\r\n对 macOS 应用程序进行代码签名，包括 Sparkle 框架的各个组件"),
	Run: func(cmd *cobra.Command, args []string) {
		scheme, _ := cmd.Flags().GetString("scheme")
		buildPath, _ := cmd.Flags().GetString("build-path")
		signingIdentity, _ := cmd.Flags().GetString("identity")
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		// 显示环境信息
		showCodesignEnvironmentInfo(scheme, buildPath, signingIdentity, verbose)
		
		// 检查必需的参数
		if scheme == "" {
			scheme = detectScheme()
			if scheme == "" {
				color.Error.Println("❌ 错误: 未设置 SCHEME 且无法自动检测")
				showAvailableSchemes()
				os.Exit(1)
			}
		}
		
		if signingIdentity == "" {
			color.Error.Println("❌ 错误: 未设置代码签名身份")
			showAvailableIdentities()
			os.Exit(1)
		}
		
		// 设置默认值
		if buildPath == "" {
			buildPath = "./temp"
		}
		
		// 构建应用路径
		appPath := buildAppPath(buildPath, scheme)
		
		// 检查应用是否存在
		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			color.Error.Printf("❌ 应用程序不存在: %s\n", appPath)
			searchAndSuggestAppPaths(scheme)
			os.Exit(1)
		}
		
		// 显示应用信息
		showAppInfo(appPath, scheme)
		
		// 执行代码签名
		err := performCodesign(appPath, signingIdentity, verbose)
		if err != nil {
			color.Error.Printf("❌ 代码签名失败: %s\n", err.Error())
			os.Exit(1)
		}
		
		color.Success.Println("✅ 代码签名成功完成！")
		
		// 显示开发路线图
		showDevelopmentRoadmap("codesign")
	},
}

func init() {
	codesignCmd.Flags().StringP("scheme", "s", "", "应用方案名称")
	codesignCmd.Flags().StringP("build-path", "b", "./temp", "构建输出路径")
	codesignCmd.Flags().StringP("identity", "i", "", "代码签名身份")
	codesignCmd.Flags().BoolP("verbose", "v", false, "显示详细签名日志")
}

// showCodesignEnvironmentInfo 显示代码签名环境信息
func showCodesignEnvironmentInfo(scheme, buildPath, signingIdentity string, verbose bool) {
	color.Blue.Println("===========================================")
	color.Blue.Println("         代码签名环境信息                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 系统信息
	color.Green.Println("📱 系统信息:")
	fmt.Printf("   操作系统: %s %s\n", runtime.GOOS, runtime.GOARCH)
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("   主机名称: %s\n", hostname)
	}
	if user := os.Getenv("USER"); user != "" {
		fmt.Printf("   当前用户: %s\n", user)
	}
	fmt.Printf("   当前时间: %s\n", getCommandOutput("date", "+%Y-%m-%d %H:%M:%S"))
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
			fmt.Printf("   Swift 版本: %s\n", strings.TrimSpace(lines[0]))
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
	
	// 签名环境变量
	color.Green.Println("🌍 签名环境变量:")
	fmt.Printf("   应用方案: %s\n", scheme)
	fmt.Printf("   构建路径: %s\n", buildPath)
	fmt.Printf("   签名身份: %s\n", signingIdentity)
	fmt.Printf("   详细日志: %t\n", verbose)
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("   工作目录: %s\n", cwd)
	}
	fmt.Println()
}

// showAvailableIdentities 显示可用的代码签名证书
func showAvailableIdentities() {
	color.Yellow.Println("正在检查可用的代码签名证书...")
	
	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	output, err := cmd.Output()
	if err != nil {
		color.Error.Printf("❌ 无法获取代码签名证书: %v\n", err)
		return
	}
	
	lines := strings.Split(string(output), "\n")
	color.Green.Println("📋 检测到的可用代码签名证书:")
	
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Developer ID Application") || 
		   strings.Contains(line, "Apple Development") || 
		   strings.Contains(line, "Mac Developer") {
			
			// 提取证书名称
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start < end {
				certName := line[start+1 : end]
				
				// 根据证书类型添加说明
				if strings.Contains(certName, "Developer ID Application") {
					fmt.Printf("  - %s [分发证书 - 可公开分发]\n", certName)
				} else if strings.Contains(certName, "Apple Development") {
					fmt.Printf("  - %s [开发证书 - 仅限开发测试]\n", certName)
				} else if strings.Contains(certName, "Mac Developer") {
					fmt.Printf("  - %s [开发证书 - 仅限开发测试]\n", certName)
				} else {
					fmt.Printf("  - %s\n", certName)
				}
				count++
			}
		}
	}
	
	if count == 0 {
		color.Error.Println("   未检测到可用的代码签名证书")
	}
	
	fmt.Println()
	color.Yellow.Println("💡 使用示例:")
	color.Cyan.Println(`   go run main.go xcode codesign --identity "Developer ID Application: Your Name (XXXXXXXXXX)"`)
	fmt.Println()
	color.Yellow.Println("📋 证书类型说明:")
	fmt.Println("   🟢 Developer ID Application: 用于 Mac App Store 外分发，可被所有用户安装")
	fmt.Println("   🟡 Apple Development: 用于开发测试，仅限开发团队内部使用")
	fmt.Println("   🔴 Mac App Store: 用于 App Store 上架（需单独申请）")
}

// buildAppPath 构建应用路径
func buildAppPath(buildPath, scheme string) string {
	// 检查 BuildPath 是否已经包含 Build/Products 路径
	if strings.Contains(buildPath, "/Build/Products/") {
		// 如果已经包含，直接使用
		return filepath.Join(buildPath, scheme+".app")
	} else {
		// 如果不包含，添加标准路径
		return filepath.Join(buildPath, "Build/Products/Release", scheme+".app")
	}
}

// searchAndSuggestAppPaths 搜索并建议可能的应用路径
func searchAndSuggestAppPaths(scheme string) {
	color.Green.Println("🔍 搜索可能的应用程序位置...")
	
	possiblePaths := []string{
		fmt.Sprintf("./temp/Build/Products/Release/%s.app", scheme),
		fmt.Sprintf("./temp/Build/Products/Debug/%s.app", scheme),
		fmt.Sprintf("./Build/Products/Release/%s.app", scheme),
		fmt.Sprintf("./Build/Products/Debug/%s.app", scheme),
		fmt.Sprintf("./build/Release/%s.app", scheme),
		fmt.Sprintf("./build/Debug/%s.app", scheme),
	}
	
	foundApps := []string{}
	
	// 检查预定义路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			foundApps = append(foundApps, path)
		}
	}
	
	// 使用 find 命令搜索更多可能的位置
	cmd := exec.Command("find", ".", "-name", scheme+".app", "-type", "d", "-not", "-path", "*/.*")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				// 避免重复添加
				alreadyFound := false
				for _, existing := range foundApps {
					if existing == line {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound {
					foundApps = append(foundApps, line)
				}
			}
		}
	}
	
	if len(foundApps) > 0 {
		fmt.Println()
		color.Info.Printf("📍 发现 %d 个可能的应用程序:\n", len(foundApps))
		for i, appPath := range foundApps {
			appSize := "未知"
			if info, err := os.Stat(appPath); err == nil && info.IsDir() {
				if sizeOutput := getCommandOutput("du", "-sh", appPath); sizeOutput != "" {
					parts := strings.Fields(sizeOutput)
					if len(parts) > 0 {
						appSize = parts[0]
					}
				}
			}
			fmt.Printf("   %d. %s (%s)\n", i+1, appPath, appSize)
		}
		fmt.Println()
		color.Info.Println("💡 建议: 请设置正确的构建路径，例如:")
		fmt.Println()
		for _, appPath := range foundApps {
			buildPath := filepath.Dir(appPath)
			fmt.Printf(" go run main.go xcode codesign --build-path '%s'\n", buildPath)
		}
		fmt.Println()
	} else {
		color.Info.Println("💡 建议: 请先运行构建命令: go run main.go xcode build")
	}
}

// showAppInfo 显示应用信息
func showAppInfo(appPath, scheme string) {
	color.Green.Println("🎯 应用程序信息:")
	fmt.Printf("   应用路径: %s\n", appPath)
	
	// 读取 Info.plist
	infoPath := filepath.Join(appPath, "Contents/Info.plist")
	if _, err := os.Stat(infoPath); err == nil {
		if version := getCommandOutput("plutil", "-p", infoPath); version != "" {
			lines := strings.Split(version, "\n")
			for _, line := range lines {
				if strings.Contains(line, "CFBundleShortVersionString") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   应用版本: %s\n", parts[3])
					}
				} else if strings.Contains(line, "CFBundleVersion") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   构建版本: %s\n", parts[3])
					}
				} else if strings.Contains(line, "CFBundleIdentifier") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   Bundle ID: %s\n", parts[3])
					}
				}
			}
		}
		fmt.Printf("   应用名称: %s\n", scheme)
	} else {
		color.Yellow.Println("   ⚠️  无法读取应用信息")
	}
	fmt.Println()
}

// performCodesign 执行代码签名
func performCodesign(appPath, signingIdentity string, verbose bool) error {
	color.Blue.Println("===========================================")
	color.Yellow.Println("🔐 开始代码签名过程...")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 签名 Sparkle 框架组件（如果存在）
	sparkleFramework := filepath.Join(appPath, "Contents/Frameworks/Sparkle.framework")
	if _, err := os.Stat(sparkleFramework); err == nil {
		color.Info.Println("🔧 签名 Sparkle 框架组件...")
		
		// 签名 Sparkle 框架内的各个组件
		sparkleComponents := []string{
			"Versions/B/Resources/Autoupdate.app/Contents/MacOS/Autoupdate",
			"Versions/B/Resources/Autoupdate.app",
			"Versions/B/Sparkle",
			"Sparkle",
		}
		
		for _, component := range sparkleComponents {
			componentPath := filepath.Join(sparkleFramework, component)
			if _, err := os.Stat(componentPath); err == nil {
				err := executeCodesign(componentPath, signingIdentity, verbose)
				if err != nil {
					return fmt.Errorf("签名 Sparkle 组件失败 (%s): %v", component, err)
				}
			}
		}
	}
	
	// 签名主应用程序
	color.Info.Println("🔧 签名主应用程序...")
	err := executeCodesign(appPath, signingIdentity, verbose)
	if err != nil {
		return fmt.Errorf("签名主应用程序失败: %v", err)
	}
	
	// 验证代码签名
	color.Info.Println("🔍 验证代码签名...")
	err = verifyCodesign(appPath, verbose)
	if err != nil {
		return fmt.Errorf("验证代码签名失败: %v", err)
	}
	
	return nil
}

// executeCodesign 执行代码签名命令
func executeCodesign(path, identity string, verbose bool) error {
	args := []string{
		"--sign", identity,
		"--force",
		"--options", "runtime",
		"--deep",
		"--timestamp",
		path,
	}
	
	if verbose {
		args = append([]string{"--verbose"}, args...)
		fmt.Printf("执行命令: codesign %s\n", strings.Join(args, " "))
	}
	
	cmd := exec.Command("codesign", args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	err := cmd.Run()
	if err != nil {
		return err
	}
	
	if verbose {
		color.Success.Printf("✅ 签名完成: %s\n", path)
	}
	
	return nil
}

// verifyCodesign 验证代码签名
func verifyCodesign(appPath string, verbose bool) error {
	args := []string{"--verify", "--deep", "--strict", appPath}
	
	if verbose {
		args = append([]string{"--verbose"}, args...)
		fmt.Printf("执行命令: codesign %s\n", strings.Join(args, " "))
	}
	
	cmd := exec.Command("codesign", args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	err := cmd.Run()
	if err != nil {
		return err
	}
	
	color.Success.Println("✅ 代码签名验证通过")
	return nil
}