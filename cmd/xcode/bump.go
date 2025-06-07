package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "递增应用版本号",
	Long:  color.Success.Render("\r\n自动递增应用程序的修订版本号（最后一位数字）"),
	Run: func(cmd *cobra.Command, args []string) {
		projectFile, _ := cmd.Flags().GetString("project")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		// 显示环境信息
		showBumpEnvironmentInfo()
		
		// 如果没有指定项目文件，自动查找
		if projectFile == "" {
			var err error
			projectFile, err = findPbxprojFile()
			if err != nil {
				color.Error.Printf("❌ %s\n", err.Error())
				os.Exit(1)
			}
		}
		
		color.Info.Printf("📁 项目文件: %s\n", projectFile)
		
		// 获取当前版本号
		currentVersion, err := getVersionFromProject(projectFile)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(2)
		}
		
		color.Info.Printf("📱 当前版本: %s\n", currentVersion)
		
		// 计算新版本号
		newVersion, err := incrementVersion(currentVersion)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(3)
		}
		
		color.Success.Printf("🚀 新版本: %s\n", newVersion)
		
		if dryRun {
			color.Yellow.Println("🔍 预览模式，不会实际修改文件")
			return
		}
		
		// 更新项目文件
		err = updateVersionInProject(projectFile, currentVersion, newVersion)
		if err != nil {
			color.Error.Printf("❌ 更新版本失败: %s\n", err.Error())
			os.Exit(4)
		}
		
		color.Success.Println("✅ 版本号更新成功！")
		
		// 显示 Git 状态
		showGitStatus()
		
		// 显示开发路线图
		showDevelopmentRoadmap("version")
	},
}

func init() {
	bumpCmd.Flags().StringP("project", "p", "", "指定 .pbxproj 文件路径")
	bumpCmd.Flags().Bool("dry-run", false, "预览模式，不实际修改文件")
}

// showBumpEnvironmentInfo 显示版本管理环境信息
func showBumpEnvironmentInfo() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         版本管理环境信息                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 系统信息
	color.Green.Println("📱 系统信息:")
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("   主机名称: %s\n", hostname)
	}
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("   工作目录: %s\n", cwd)
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
}

// incrementVersion 递增版本号的最后一位
func incrementVersion(version string) (string, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("版本号格式不正确，期望格式: x.y.z")
	}
	
	// 解析最后一位数字
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("无法解析版本号的修订版本: %v", err)
	}
	
	// 递增
	patch++
	parts[2] = strconv.Itoa(patch)
	
	return strings.Join(parts, "."), nil
}

// updateVersionInProject 更新项目文件中的版本号
func updateVersionInProject(projectFile, oldVersion, newVersion string) error {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return fmt.Errorf("无法读取项目文件: %v", err)
	}
	
	// 替换版本号
	oldPattern := fmt.Sprintf("MARKETING_VERSION = %s", oldVersion)
	newPattern := fmt.Sprintf("MARKETING_VERSION = %s", newVersion)
	newContent := strings.ReplaceAll(string(content), oldPattern, newPattern)
	
	// 写回文件
	err = os.WriteFile(projectFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("无法写入项目文件: %v", err)
	}
	
	return nil
}

// showGitStatus 显示 Git 状态
func showGitStatus() {
	color.Green.Println("📝 Git 状态变更:")
	
	if status := getCommandOutput("git", "status", "--porcelain"); status != "" {
		lines := strings.Split(strings.TrimSpace(status), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	} else {
		fmt.Println("   无变更")
	}
	fmt.Println()
	
	color.Yellow.Println("💡 提示: 请手动提交 Git 变更")
	color.Cyan.Println("   git add .")
	color.Cyan.Printf("   git commit -m \"Bump version to %s\"\n", "新版本")
	fmt.Println()
}

// getCommandOutput 执行命令并返回输出
func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// showDevelopmentRoadmap 显示开发路线图
func showDevelopmentRoadmap(currentStep string) {
	fmt.Println()
	color.Blue.Println("===========================================")
	color.Blue.Println("         🗺️  开发分发路线图                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	steps := []string{
		"setup:⚙️ 环境设置:配置代码签名环境",
		"version:📝 版本管理:查看或更新应用版本号",
		"build:🔨 构建应用:编译源代码，生成可执行文件",
		"codesign:🔐 代码签名:为应用添加数字签名，确保安全性",
		"package:📦 打包分发:创建 DMG 安装包",
		"notarize:✅ 公证验证:Apple 官方验证（可选）",
		"distribute:🚀 发布分发:上传到分发平台或直接分发",
	}
	
	color.Cyan.Print("📍 当前位置: ")
	switch currentStep {
	case "setup":
		color.Green.Println("环境设置")
	case "version":
		color.Green.Println("版本管理")
	case "build":
		color.Green.Println("构建应用")
	case "codesign":
		color.Green.Println("代码签名")
	case "package":
		color.Green.Println("打包分发")
	case "notarize":
		color.Green.Println("公证验证")
	case "distribute":
		color.Green.Println("发布分发")
	default:
		color.Yellow.Println("未知步骤")
	}
	fmt.Println()
	
	// 显示路线图
	for _, step := range steps {
		parts := strings.Split(step, ":")
		stepId := parts[0]
		stepIcon := parts[1]
		stepDesc := parts[2]
		
		if stepId == currentStep {
			color.Green.Printf("▶ %s %s\n", stepIcon, stepDesc)
		} else {
			fmt.Printf("  %s %s\n", stepIcon, stepDesc)
		}
	}
	
	fmt.Println()
	color.Yellow.Println("💡 下一步建议:")
	switch currentStep {
	case "setup":
		color.Cyan.Println("   查看版本信息: go run main.go xcode version")
		color.Cyan.Println("   或直接构建应用: go run main.go xcode build")
	case "version":
		color.Cyan.Println("   构建应用: go run main.go xcode build")
	case "build":
		color.Cyan.Println("   运行代码签名: go run main.go xcode codesign")
	case "codesign":
		color.Cyan.Println("   创建安装包: go run main.go xcode package")
	case "package":
		fmt.Println("   进行公证验证或直接分发应用")
	case "notarize":
		fmt.Println("   发布到分发平台或提供下载链接")
	case "distribute":
		fmt.Println("   🎉 开发分发流程已完成！")
	}
	
	fmt.Println()
	color.Blue.Println("===========================================")
}