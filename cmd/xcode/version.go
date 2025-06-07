package xcode

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "获取应用版本号",
	Long:  color.Success.Render("\r\n从 Xcode 项目配置文件中获取应用程序的营销版本号（MARKETING_VERSION）"),
	Run: func(cmd *cobra.Command, args []string) {
		projectFile, _ := cmd.Flags().GetString("project")
		
		// 如果没有指定项目文件，自动查找
		if projectFile == "" {
			var err error
			projectFile, err = findPbxprojFile()
			if err != nil {
				color.Error.Printf("❌ %s\n", err.Error())
				os.Exit(1)
			}
		}
		
		// 获取版本号
		version, err := getVersionFromProject(projectFile)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(2)
		}
		
		color.Success.Printf("📱 当前版本: %s\n", version)
	},
}

func init() {
	versionCmd.Flags().StringP("project", "p", "", "指定 .pbxproj 文件路径")
}

// findPbxprojFile 自动查找 .pbxproj 文件
func findPbxprojFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("无法获取当前目录: %v", err)
	}
	
	// 在当前目录及其子目录中查找 .pbxproj 文件（排除 Resources 和 temp 目录）
	var projectFile string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续查找
		}
		
		// 跳过深度超过2层的目录
		relPath, _ := filepath.Rel(cwd, path)
		if strings.Count(relPath, string(filepath.Separator)) > 2 {
			return filepath.SkipDir
		}
		
		// 跳过 Resources 和 temp 目录
		if info.IsDir() && (strings.Contains(path, "Resources") || strings.Contains(path, "temp")) {
			return filepath.SkipDir
		}
		
		// 查找 .pbxproj 文件
		if strings.HasSuffix(path, ".pbxproj") {
			projectFile = path
			return fmt.Errorf("found") // 用错误来停止遍历
		}
		
		return nil
	})
	
	if projectFile == "" {
		return "", fmt.Errorf("未找到 .pbxproj 配置文件")
	}
	
	return projectFile, nil
}

// getVersionFromProject 从项目文件中提取版本号
func getVersionFromProject(projectFile string) (string, error) {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return "", fmt.Errorf("无法读取项目文件: %v", err)
	}
	
	// 使用正则表达式查找 MARKETING_VERSION
	re := regexp.MustCompile(`MARKETING_VERSION\s*=\s*([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(string(content))
	
	if len(matches) < 2 {
		return "", fmt.Errorf("未找到 MARKETING_VERSION")
	}
	
	return matches[1], nil
}